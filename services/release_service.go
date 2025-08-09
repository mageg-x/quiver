package services

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm/clause"
	"quiver/cache"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
	"quiver/utils"
	"strings"
	"time"
)

// ReleaseService 命名空间服务
type ReleaseService struct{}

// NewReleaseService 创建命名空间服务实例
func NewReleaseService() *ReleaseService {
	return &ReleaseService{}
}
func OnKeyUpdate4Release(env, key string) {
	if len(env) == 0 || len(key) == 0 {
		logger.GetLogger("quiver").Errorf("env %s or key %s is empty", env, key)
		return
	}

	parts := strings.Split(key, ":")
	if len(parts) != 6 {
		logger.GetLogger("quiver").Errorf("key %s is invalid format", key)
		return
	}
	appName, clusterName, namespaceName := parts[3], parts[4], parts[5]

	if !utils.ValidateAppName(appName) || !utils.ValidateClusterName(clusterName) || !utils.ValidateNamespaceName(namespaceName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s or cluster_name %s or namespace_name %s for key %s",
			appName, clusterName, namespaceName, key)
		return
	}

	s := NewReleaseService()
	_, err := s.GetLatestReleaseAll(env, appName, clusterName, namespaceName)
	logger.GetLogger("quiver").Infof("refresh release cache: %s %s %s %s for key %s, %v",
		env, appName, clusterName, namespaceName, key, err)
	return
}

// PublishRelease 发布指定命名空间的配置
func (s *ReleaseService) PublishRelease(env, appName, clusterName, namespaceName,
	releaseName, operator, comment string) (*models.NamespaceRelease, error) {
	// 1. 校验 App/Cluster/Namespace 是否存在，并获取 IDs
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, err
	}

	db := database.GetDB(env)

	// 2. 生成 release_id
	releaseID, err := utils.GenerateReleaseID()
	if err != nil {
		logger.GetLogger("quiver").Errorf("generate release id error: %v", err)
		return nil, fmt.Errorf("failed to generate release ID: %w", err)
	}

	// 3. 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		logger.GetLogger("quiver").Errorf("begin transaction failed: %v", tx.Error)
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// 4. 分页查询未发布的 items，并收集 kv_id
	const batchSize = 1000
	var allKvIDs []uint64
	offset := 0
	totalProcessed := 0

	for {
		var items []models.Item
		if err := tx.
			Where("namespace_id = ? AND is_deleted = ?", ids.NamespaceID, 0).
			Order("id ASC").Offset(offset).Limit(batchSize).Find(&items).Error; err != nil {
			logger.GetLogger("quiver").Errorf("query items failed: %v", err)
			return nil, fmt.Errorf("failed to query items: %w", err)
		}

		if len(items) == 0 {
			break // 查询完成
		}

		// 提取 kv_id
		for _, item := range items {
			allKvIDs = append(allKvIDs, item.KVId)
		}

		// 转换为 item_release 数据
		var itemReleases []models.ItemRelease
		for _, item := range items {
			itemReleases = append(itemReleases, models.ItemRelease{
				AppID:       item.AppID,
				ClusterID:   item.ClusterID,
				NamespaceID: item.NamespaceID,
				K:           item.K,
				V:           item.V,
				KvID:        item.KVId,
			})
		}

		// 批量插入到 item_release 表，跳过重复（基于 uk_namespace_kv_id）
		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).
			Create(&itemReleases).Error; err != nil {
			logger.GetLogger("quiver").Errorf("batch insert item_release failed: %v", err)
			return nil, fmt.Errorf("failed to insert into item_release: %w", err)
		}

		totalProcessed += len(items)
		offset += len(items)
	}

	// 5. 如果没有要发布的配置项
	if len(allKvIDs) == 0 {
		return nil, fmt.Errorf("no unreleased items found for namespace %s", namespaceName)
	}

	logger.GetLogger("quiver").Infof("Publishing %d items for namespace %s", totalProcessed, namespaceName)

	// 6. 序列化 kv_id 列表为 MessagePack
	configData, err := msgpack.Marshal(allKvIDs)
	if err != nil {
		logger.GetLogger("quiver").Errorf("msgpack marshal kv_ids failed: %v", err)
		return nil, fmt.Errorf("failed to serialize config: %w", err)
	}

	// 7. 写入 namespace_release 表（核心发布记录）
	namespaceRelease := &models.NamespaceRelease{
		AppID:         ids.AppID,
		AppName:       appName,
		ClusterID:     ids.ClusterID,
		ClusterName:   clusterName,
		NamespaceID:   ids.NamespaceID,
		NamespaceName: namespaceName,
		ReleaseName:   releaseName,
		ReleaseID:     releaseID,
		Config:        configData,
		Operator:      operator,
		Comment:       comment,
	}

	if err := tx.Create(namespaceRelease).Error; err != nil {
		logger.GetLogger("quiver").Errorf("create namespace_release failed: %v", err)
		return nil, fmt.Errorf("failed to create namespace release: %w", err)
	}

	// 8. 更新 item 表：标记为已发布
	if err := tx.Model(&models.Item{}).
		Where("namespace_id = ? AND is_released = ?", ids.NamespaceID, 0).
		Update("is_released", 1).Error; err != nil {
		logger.GetLogger("quiver").Errorf("update is_released flag failed: %v", err)
		return nil, fmt.Errorf("failed to mark items as released: %w", err)
	}

	// 9. 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.GetLogger("quiver").Errorf("transaction commit failed: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	tx = nil
	logger.GetLogger("quiver").Infof("Successfully published release %s with %d items", releaseID, len(allKvIDs))

	// 10. 返回发布结果
	return namespaceRelease, nil
}

// ListRelease 获取集群下的所有命名空间
func (s *ReleaseService) ListRelease(env, appName, clusterName, namespaceName string, page, size int) ([]models.NamespaceRelease, int64, error) {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, 0, err
	}

	db := database.GetDB(env)

	var releases []models.NamespaceRelease
	var total int64

	// 构建查询条件
	query := db.Where("namespace_id = ?", ids.NamespaceID).Order("release_time DESC,id DESC")

	// 获取总数（用于分页）
	if err := query.Model(&models.NamespaceRelease{}).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count release num : %s for cluster: %s", err.Error(), clusterName)
		return nil, 0, fmt.Errorf("failed to count release: %w", err)
	}

	offset := (page - 1) * size

	// 执行分页查询
	if err := query.Offset(offset).Limit(size).Find(&releases).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to query release: %s", err.Error())
		return nil, 0, fmt.Errorf("failed to query release: %w", err)
	}

	return releases, total, nil
}

func (s *ReleaseService) GetFixedReleaseAll(env, releaseID string) (*models.NamespaceRelease, error) {
	if releaseID == "" {
		return nil, fmt.Errorf("releaseID is empty")
	}

	// 1、先检查缓存里面有没有
	releaseKey := fmt.Sprintf("release:%s:%s", env, releaseID)
	data, ok, _ := cache.Get(releaseKey)
	if ok {
		if len(data) > 0 {
			nr := &models.NamespaceRelease{}
			if err := msgpack.Unmarshal(data, nr); err != nil {
				// 数据错误，删除缓存
				_ = cache.Delete(releaseKey)
				logger.GetLogger("quiver").Warnf("error unmarshaling release: %s, %v", releaseKey, err)
			}

			logger.GetLogger("quiver").Infof("get release %s from cache success", releaseKey)
			return nr, nil
		} else {
			// 数据已经破坏
			_ = cache.Delete(releaseKey)
		}
	}

	// 2、从数据库中获取
	db := database.GetDB(env)
	var baseRelease models.NamespaceRelease
	if err := db.Where("release_id = ?", releaseID).First(&baseRelease).Error; err != nil {
		logger.GetLogger("quiver").Warnf("Release %s not found", releaseID)
		return nil, errors.New("no releases found")
	}

	var allKvIDs []uint64
	if len(baseRelease.Config) > 0 {
		if err := msgpack.Unmarshal(baseRelease.Config, &allKvIDs); err != nil {
			logger.GetLogger("quiver").Errorf("failed to unmarshal %s config", releaseID)
			return nil, errors.New("failed to get release data")
		}
	}
	if len(allKvIDs) == 0 {
		return nil, errors.New("no kv found")
	}

	// 3、 根据kv_id 取 item_release 表获取数据
	var allItems []models.ItemRelease
	for i := 0; i < len(allKvIDs); i += 1000 {
		end := i + 1000
		if end > len(allKvIDs) {
			end = len(allKvIDs)
		}
		chunk := allKvIDs[i:end]
		var batch []models.ItemRelease
		if err := db.Table("item_release").
			Select("kv_id, `k`, `v`"). // 可优化为只查 key，但通常一起查
			Where("namespace_id =? AND kv_id IN (?)", baseRelease.NamespaceID, chunk).
			Scan(&batch).Error; err != nil {
			logger.GetLogger("quiver").Errorf("failed to query items: %s", err.Error())
			return nil, err
		}
		// 追加到结果
		allItems = append(allItems, batch...)
	}

	baseRelease.Items = allItems
	baseRelease.KvIDs = allKvIDs

	// 4、写入cache 缓存
	if data, err := msgpack.Marshal(&baseRelease); err == nil {
		if len(data) > 0 && len(baseRelease.ReleaseID) > 0 {
			releaseKey := fmt.Sprintf("release:%s:%s", env, baseRelease.ReleaseID)
			_ = cache.Set(releaseKey, data, 3600*24*7*time.Second)
		}
	}
	return &baseRelease, nil
}

func (s *ReleaseService) GetLatestReleaseAll(env, appName, clusterName, namespaceName string) (*models.NamespaceRelease, error) {
	// 1、先检查缓存里面有没有
	r := models.NamespaceRelease{AppName: appName, ClusterName: clusterName, NamespaceName: namespaceName}
	data, ok, err := cache.Get(r.CacheKey(env))
	logger.GetLogger("quiver").Infof("cache.Get %s %v %v %v", r.CacheKey(env), data, ok, err)
	if ok {
		if len(data) > 0 {
			// 把data 转成 string
			releaseID := string(data)
			nr, err := s.GetFixedReleaseAll(env, releaseID)
			logger.GetLogger("quiver").Warnf("get fixed release %s, err: %v", releaseID, err)
			if err == nil && nr != nil {
				return nr, nil
			}
		} else {
			// 数据已经破坏
			_ = cache.Delete(r.CacheKey(env))
		}
	}

	// 2. 校验并获取 namespace ID
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, err
	}

	db := database.GetDB(env)

	// 3. 获取最近一次的 release kv_id
	var latestRelease models.NamespaceRelease
	if err := db.Where("namespace_id = ?", ids.NamespaceID).
		Order("id DESC").
		First(&latestRelease).Error; err != nil {
		logger.GetLogger("quiver").Errorf("no releases found for %s/%s/%s/%s", env, appName, clusterName, namespaceName)
		return nil, errors.New("no releases found")
	}
	if latestRelease.ReleaseID == "" || len(latestRelease.Config) == 0 {
		return nil, errors.New("release not found")
	}

	// 尝试从缓存中获取
	nr, err := s.GetFixedReleaseAll(env, latestRelease.ReleaseID)
	if err != nil || nr == nil {
		logger.GetLogger("quiver").Errorf("failed to get latest release %s %s data", env, latestRelease.ReleaseID)
		return nil, fmt.Errorf("failed to get latest release data: %w", err)
	}

	latestRelease.Items = nr.Items
	latestRelease.KvIDs = nr.KvIDs

	// 写入cache 缓存
	if data, err := msgpack.Marshal(&latestRelease); err == nil {
		if len(data) > 0 && len(latestRelease.ReleaseID) > 0 {
			_ = cache.Set(latestRelease.CacheKey(env), []byte(latestRelease.ReleaseID), 300*time.Second)
			releaseKey := fmt.Sprintf("release:%s:%s", env, latestRelease.ReleaseID)
			_ = cache.Set(releaseKey, data, 3600*24*7*time.Second)
			logger.GetLogger("quiver").Infof("write %s %s to cache", latestRelease.CacheKey(env), releaseKey)
		}
	}
	return &latestRelease, nil
}

// GetRelease 获取特定命名空间
func (s *ReleaseService) GetRelease(env, appName, clusterName, namespaceName, releaseId string) (map[string]interface{}, error) {
	// 1、从缓存中获取最近一次发布的release内容
	latestRelease, err := s.GetLatestReleaseAll(env, appName, clusterName, namespaceName)
	if err != nil {
		logger.GetLogger("quiver").Errorf("get latest release all error: %v", err)
		return nil, err
	}
	//logger.GetLogger("quiver").Infof("get latest release all: %v", latestRelease)

	latestKvIDs := latestRelease.KvIDs
	// 2、建立映射表，方便后续查找
	latestKH := make(map[uint32]uint64, len(latestKvIDs))
	latestK := make([]uint32, 0, len(latestKvIDs)) // 预分配容量，提升性能
	for _, h := range latestKvIDs {
		k := uint32(h & 0xFFFFFFFF) // 取低 32 位
		latestK = append(latestK, k)
		latestKH[k] = h
	}

	//logger.GetLogger("quiver").Infof("latestKvIDs: %v", latestKvIDs)

	// 3. 获取 base release（客户端传来的release版本）
	var baseKvIDs []uint64
	baseRelease, err := s.GetFixedReleaseAll(env, releaseId)
	if err == nil && baseRelease != nil {
		baseKvIDs = baseRelease.KvIDs
	}

	// 4、同样建立映射表，方便后续查找
	baseKH := make(map[uint32]uint64, len(baseKvIDs))
	baseK := make([]uint32, 0, len(baseKvIDs)) // 预分配容量，提升性能
	for _, h := range baseKvIDs {
		k := uint32(h & 0xFFFFFFFF) // 取低 32 位
		baseK = append(baseK, k)
		baseKH[k] = h
	}

	// 5、得到增删改三部分
	addK, bothK, deleteK := utils.Diff32(latestK, baseK)
	// 计算新增，update, 和 delete 的 kvid中有变化的部分
	var allKvIDs, addKvIDs, updateKvIDs, deleteKvIDs []uint64
	for _, k := range bothK {
		if latestKH[k] != baseKH[k] {
			updateKvIDs = append(updateKvIDs, latestKH[k])
		}
	}
	allKvIDs = append(allKvIDs, updateKvIDs...)
	for _, k := range addK {
		addKvIDs = append(addKvIDs, latestKH[k])
	}
	allKvIDs = append(allKvIDs, addKvIDs...)
	for _, k := range deleteK {
		deleteKvIDs = append(deleteKvIDs, baseKvIDs[k])
	}
	allKvIDs = append(allKvIDs, deleteKvIDs...)

	logger.GetLogger("quiver").Infof("item del %d , updated %d, add %d",
		len(deleteKvIDs), len(updateKvIDs), len(addKvIDs))

	// 6. 构建 kvid -> {k,v} 映射
	type KvPair struct {
		Key   string
		Value string
	}

	idToKv := make(map[uint64]KvPair, len(allKvIDs))

	for _, item := range latestRelease.Items {
		idToKv[item.KvID] = KvPair{
			Key:   item.K,
			Value: item.V,
		}
	}
	if baseRelease != nil {
		for _, item := range baseRelease.Items {
			idToKv[item.KvID] = KvPair{
				Key:   item.K,
				Value: item.V,
			}
		}
	}

	// 6. 构造 changed 差异（只返回 key 列表）
	addKeys, updateKeys, deleteKeys := []string{}, []string{}, []string{}
	kvList := make([]KvPair, 0, len(addKvIDs)+len(updateKvIDs))

	for _, id := range addKvIDs {
		if item, exists := idToKv[id]; exists {
			addKeys = append(addKeys, item.Key)
			kvList = append(kvList, KvPair{Key: item.Key, Value: item.Value})
		}
	}

	for _, id := range updateKvIDs {
		if item, exists := idToKv[id]; exists {
			updateKeys = append(updateKeys, item.Key)
			kvList = append(kvList, KvPair{Key: item.Key, Value: item.Value})
		}
	}

	for _, id := range deleteKvIDs {
		if item, exists := idToKv[id]; exists {
			deleteKeys = append(deleteKeys, item.Key)
		}
	}

	// 7、回idToKv 和 addKeys, updateKeys, deleteKeys
	ret := map[string]interface{}{
		"env":            env,
		"app_name":       appName,
		"cluster_name":   clusterName,
		"namespace_name": namespaceName,
		"release_id":     latestRelease.ReleaseID,
		"release_name":   latestRelease.ReleaseName,
		"release_time":   latestRelease.ReleaseTime,
		"operator":       latestRelease.Operator,
		"comment":        latestRelease.Comment,
		"items":          kvList,
	}

	if releaseId != "" {
		ret["changed"] = map[string]interface{}{
			"added":   addKeys,
			"updated": updateKeys,
			"deleted": deleteKeys,
		}
	}

	//logger.GetLogger("quiver").Infof("GetRelease: %+v", ret)
	return ret, nil
}

func (s *ReleaseService) RollbackRelease(env, appName, clusterName, namespaceName, releaseId, operator, comment string) (*models.NamespaceRelease, error) {
	// 1. 校验并获取 namespace ID
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, err
	}

	db := database.GetDB(env)

	// 2、判断待回滚的 releaseId 是否存在
	if releaseId == "" {
		logger.GetLogger("quiver").Errorf("rollback release_id is empty")
		return nil, errors.New("release_id not found")
	}

	var release models.NamespaceRelease
	err = db.Where("release_id = ?", releaseId).First(&release).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("rollback release_id %s not found", releaseId)
		return nil, errors.New("release_id not found")
	}

	// 3. 获取 base release（客户端传来的release版本）
	var baseKvIDs []uint64
	if len(release.Config) > 0 {
		if err := msgpack.Unmarshal(release.Config, &baseKvIDs); err != nil {
			logger.GetLogger("quiver").Errorf("failed to unmarshal %s config", releaseId)
			return nil, errors.New("failed to get release data")
		}
	}

	// 4. 生成 release_id
	newReleaseID, err := utils.GenerateReleaseID()
	if err != nil {
		logger.GetLogger("quiver").Errorf("generate release id error: %v", err)
		return nil, fmt.Errorf("failed to generate release ID: %w", err)
	}

	// 5. 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		logger.GetLogger("quiver").Errorf("begin transaction failed: %v", tx.Error)
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	//6、设置item 表(草稿区) 的 namespace_id 为 ids 的 namespaceId 的 is_released 字段为 0
	err = tx.Model(&models.Item{}).Where("namespace_id = ? AND is_deleted = 0", ids.NamespaceID).
		Update("is_released", 0).Error

	if err != nil {
		logger.GetLogger("quiver").Errorf("update items failed: %v", err)
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	//7、重新发布该版本
	release.ID = 0
	release.ReleaseID = newReleaseID
	release.ReleaseName = "rollback-" + release.ReleaseName
	release.Operator = operator
	release.Comment = comment
	logger.GetLogger("quiver").Errorf("step 7 release :%v", release)
	if err := tx.Create(&release).Error; err != nil {
		logger.GetLogger("quiver").Errorf("create namespace_release failed: %v", err)
		return nil, fmt.Errorf("failed to create namespace release: %w", err)
	}

	// 8、 恢复item表中没有变化的 的is_release marked
	const batchSize = 1000
	for i := 0; i < len(baseKvIDs); i += batchSize {
		chunk := baseKvIDs[i:min(i+batchSize, len(baseKvIDs))]

		err = tx.Model(&models.Item{}).
			Where("namespace_id = ? AND kv_id IN (?)", ids.NamespaceID, chunk).
			Update("is_released", 1).Error

		if err != nil {
			logger.GetLogger("quiver").Errorf("update batch failed: %v", err)
			return nil, fmt.Errorf("failed to update batch: %w", err)
		}
	}

	// 9. 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.GetLogger("quiver").Errorf("transaction commit failed: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	tx = nil
	return &release, nil
}
