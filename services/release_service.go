package services

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm/clause"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
	"quiver/utils"
)

// ReleaseService 命名空间服务
type ReleaseService struct{}

// NewReleaseService 创建命名空间服务实例
func NewReleaseService() *ReleaseService {
	return &ReleaseService{}
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

// GetRelease 获取特定命名空间
func (s *ReleaseService) GetRelease(env, appName, clusterName, namespaceName, releaseId string) (map[string]interface{}, error) {
	// 1. 校验并获取 namespace ID
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, err
	}

	db := database.GetDB(env)

	// 2. 获取最近一次的 release
	var latestRelease models.NamespaceRelease
	if err := db.Where("namespace_id = ?", ids.NamespaceID).
		Order("id DESC").
		First(&latestRelease).Error; err != nil {
		logger.GetLogger("quiver").Errorf("no releases found for %s/%s/%s/%s", env, appName, clusterName, namespaceName)
		return nil, errors.New("no releases found")
	}

	var latestKvIDs []uint64

	if releaseId != latestRelease.ReleaseID {
		if err := msgpack.Unmarshal(latestRelease.Config, &latestKvIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal latest config: %w", err)
		}
	}

	latestKK := make(map[uint32]uint64, len(latestKvIDs))
	latestK := make([]uint32, 0, len(latestKvIDs)) // 预分配容量，提升性能
	for _, hash := range latestKvIDs {
		k := uint32(hash & 0xFFFFFFFF) // 取低 32 位
		latestK = append(latestK, k)
		latestKK[k] = hash
	}

	logger.GetLogger("quiver").Errorf("latestKvIDs: %v", latestKvIDs)

	// 3. 获取 base release（客户端传来的release版本）
	var baseKvIDs []uint64

	if releaseId != "" && releaseId != latestRelease.ReleaseID {
		var baseRelease models.NamespaceRelease
		if err := db.Where("release_id = ?", ids.NamespaceID, releaseId).First(&baseRelease).Error; err != nil {
			logger.GetLogger("quiver").Errorf("Release %s not found", releaseId)
		}

		if len(baseRelease.Config) > 0 {
			if err := msgpack.Unmarshal(baseRelease.Config, &baseKvIDs); err != nil {
				logger.GetLogger("quiver").Errorf("failed to unmarshal %s config", releaseId)
			}
		}
	}

	baseKK := make(map[uint32]uint64, len(baseKvIDs))
	baseK := make([]uint32, 0, len(baseKvIDs)) // 预分配容量，提升性能
	for _, hash := range baseKvIDs {
		k := uint32(hash & 0xFFFFFFFF) // 取低 32 位
		baseK = append(baseK, k)
		baseKK[k] = hash
	}

	// 4、得到增删改三部分
	addK, bothK, deleteK := utils.Diff32(latestK, baseK)
	// 计算新增，update, 和 delete 的 kvid中有变化的部分
	var allKvIDs, addKvIDs, updateKvIDs, deleteKvIDs []uint64
	for _, k := range bothK {
		if latestKK[k] != baseKK[k] {
			updateKvIDs = append(updateKvIDs, latestKK[k])
		}
	}
	allKvIDs = append(allKvIDs, updateKvIDs...)
	for _, k := range addK {
		addKvIDs = append(addKvIDs, latestKK[k])
	}
	allKvIDs = append(allKvIDs, addKvIDs...)
	for _, k := range deleteK {
		deleteKvIDs = append(deleteKvIDs, baseKvIDs[k])
	}
	allKvIDs = append(allKvIDs, deleteKvIDs...)

	logger.GetLogger("quiver").Errorf("item num %d vs %d", len(latestKvIDs)+len(baseKvIDs), len(allKvIDs))

	// 4. 批量查询所有需要的 items（只查 key，value 仅用于 update 判断，但此处省略比较）
	var allItems []models.ItemRelease
	// 分批读出
	for i := 0; i < len(allKvIDs); i += 1000 {
		end := i + 1000
		if end > len(allKvIDs) {
			end = len(allKvIDs)
		}
		chunk := allKvIDs[i:end]
		var batch []models.ItemRelease
		if err := db.Table("item_release").
			Select("kv_id, `k`, `v`"). // 可优化为只查 key，但通常一起查
			Where("namespace_id =? AND kv_id IN (?)", ids.NamespaceID, chunk).
			Scan(&batch).Error; err != nil {
			logger.GetLogger("quiver").Errorf("failed to query items: %s", err.Error())
			return nil, err
		}
		// 追加到结果
		allItems = append(allItems, batch...)
	}

	// 5. 构建 kvid -> {k,v} 映射
	type KvPair struct {
		Key   string
		Value string
	}

	idToKv := make(map[uint64]KvPair, len(allItems))
	for _, item := range allItems {
		idToKv[item.KvID] = KvPair{
			Key:   item.K,
			Value: item.V,
		}
	}

	// 6. 构造 changed 差异（只返回 key 列表）
	addKeys, updateKeys, deleteKeys := []string{}, []string{}, []string{}
	if len(baseKvIDs) > 0 {
		for _, id := range addKvIDs {
			if item, exists := idToKv[id]; exists {
				addKeys = append(addKeys, item.Key)
			}
		}

		for _, id := range updateKvIDs {
			if item, exists := idToKv[id]; exists {
				updateKeys = append(updateKeys, item.Key)
			}
		}

		for _, id := range deleteKvIDs {
			if item, exists := idToKv[id]; exists {
				deleteKeys = append(deleteKeys, item.Key)
				// 从 idToKv 中删除
				delete(idToKv, id)
			}
		}
	}

	kvList := make([]KvPair, 0, len(idToKv))
	for _, kv := range idToKv {
		kvList = append(kvList, kv)
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
	logger.GetLogger("quiver").Infof("GetRelease: %+v", ret)
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
