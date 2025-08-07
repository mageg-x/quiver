package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
)

type PermissionService struct{}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

func (s *PermissionService) CreatePermission(env string, perm *models.Permission) error {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return errors.New("db not initialized")
	}

	if perm == nil {
		logger.GetLogger("quiver").Errorf("permission object cannot be nil")
		return errors.New("permission object cannot be nil")
	}
	if perm.UserID == 0 {
		logger.GetLogger("quiver").Errorf("user_id is required")
		return errors.New("user_id is required")
	}
	validResourceTypes := map[string]bool{"APP": true, "CLUSTER": true, "NAMESPACE": true}
	if !validResourceTypes[perm.ResourceType] {
		logger.GetLogger("quiver").Errorf("invalid resource type %s", perm.ResourceType)
		return errors.New("invalid resource type")
	}

	if perm.ResourceID == 0 {
		logger.GetLogger("quiver").Errorf("required resource_id")
		return errors.New("resource_id is required")
	}

	if perm.Action == "" {
		logger.GetLogger("quiver").Errorf("required action")
		return errors.New("action is required")
	}

	// 2. 校验 UserID 是否存在
	var userCount int64
	if err := db.Model(&models.User{}).Where("id = ?", perm.UserID).Count(&userCount).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return err
	}
	if userCount == 0 {
		logger.GetLogger("quiver").Errorf("user does not exist")
		return errors.New("user does not exist")
	}

	// 3. 校验 ResourceID 是否存在 (根据 ResourceType 查询不同表)
	var modelMap = map[string]interface{}{
		"APP":       &models.App{},
		"CLUSTER":   &models.Cluster{},
		"NAMESPACE": &models.Namespace{},
	}
	model, ok := modelMap[perm.ResourceType]
	if !ok {
		logger.GetLogger("quiver").Errorf("unexpected resource type")
		return errors.New("unexpected resource type")
	}
	var count int64
	err := db.Model(model).Where("id = ?", perm.ResourceID).Count(&count).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("error querying %s existence: %v", perm.ResourceType, err)
		return err
	}
	if count == 0 {
		logger.GetLogger("quiver").Errorf("%s with id %d does not exist", perm.ResourceType, perm.ResourceID)
		return fmt.Errorf("%s with id %d does not exist", perm.ResourceType, perm.ResourceID)
	}

	return db.Create(perm).Error
}

func (s *PermissionService) GetPermission(env string, userID uint64, resourceId uint64) (*models.Permission, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	if resourceId == 0 {
		return nil, errors.New("resource_id is required")
	}

	// 2. 校验 UserID 是否存在
	var userCount int64
	if err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return nil, err
	}
	if userCount == 0 {
		return nil, errors.New("user does not exist")
	}

	var perm models.Permission
	if err := db.Where("id = ?", resourceId).First(&perm).Error; err != nil {
		logger.GetLogger("quiver").Errorf("permission not found for id %d", resourceId)
		return nil, errors.New("permission not found")
	}
	return &perm, nil
}

func (s *PermissionService) ListPermission(env string, userID uint64) ([]models.Permission, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("user_id is required")
		return nil, errors.New("user_id is required")
	}
	// 2. 校验 UserID 是否存在
	var userCount int64
	if err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return nil, err
	}
	if userCount == 0 {
		logger.GetLogger("quiver").Errorf("user does not exist")
		return nil, errors.New("user does not exist")
	}

	var perms []models.Permission
	if err := db.Where("user_id = ?", userID).Find(&perms).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error listing permissions for user %d: %v", userID, err)
		return nil, err
	}
	return perms, nil
}

func (s *PermissionService) UpdatePermission(env string, update *models.Permission) (*models.Permission, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}
	if update == nil {
		logger.GetLogger("quiver").Errorf("permission is nil")
		return nil, errors.New("permission is required")
	}

	if update.UserID == 0 {
		logger.GetLogger("quiver").Errorf("user_id is required")
		return nil, errors.New("user_id is required")
	}

	validResourceTypes := map[string]bool{"APP": true, "CLUSTER": true, "NAMESPACE": true}
	if !validResourceTypes[update.ResourceType] {
		logger.GetLogger("quiver").Errorf("invalid resource type %s", update.ResourceType)
		return nil, errors.New("invalid resource type")
	}
	if update.ResourceID == 0 {
		logger.GetLogger("quiver").Errorf("resource_id is nil")
		return nil, errors.New("resource_id is required")
	}

	if update.Action == "" {
		logger.GetLogger("quiver").Errorf("action is required")
		return nil, errors.New("action is required")
	}
	// 2. 校验 UserID 是否存在
	var userCount int64
	if err := db.Model(&models.User{}).Where("id = ?", update.UserID).Count(&userCount).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return nil, err
	}
	if userCount == 0 {
		return nil, errors.New("user does not exist")
	}

	var perm models.Permission
	if err := db.Where("id = ?", update.ID).First(&perm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger("quiver").Errorf("permission not found for id %d", update.ID)
			return nil, errors.New("permission not found")
		}
		logger.GetLogger("quiver").Errorf("error querying permission for id %d: %v", update.ID, err)
		return nil, errors.New("query failed")
	}

	// 3. 校验 ResourceID 是否存在 (根据 ResourceType 查询不同表)
	var modelMap = map[string]interface{}{
		"APP":       &models.App{},
		"CLUSTER":   &models.Cluster{},
		"NAMESPACE": &models.Namespace{},
	}
	model, ok := modelMap[perm.ResourceType]
	if !ok {
		logger.GetLogger("quiver").Errorf("unexpected resource type")
		return nil, errors.New("unexpected resource type")
	}
	var count int64
	err := db.Model(model).Where("id = ?", perm.ResourceID).Count(&count).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("error querying %s existence: %v", perm.ResourceType, err)
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("%s with id %d does not exist", perm.ResourceType, perm.ResourceID)
	}

	if err := db.Model(&perm).Updates(update).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error updating permission for id %d: %v", update.ID, err)
		return nil, errors.New("update failed")
	}

	return s.GetPermission(env, update.UserID, update.ResourceID)
}

func (s *PermissionService) DeletePermission(env string, userID uint64, resourceId uint64) error {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return errors.New("db not initialized")
	}

	if userID == 0 {
		return errors.New("user_id is required")
	}
	// 2. 校验 UserID 是否存在
	var userCount int64
	if err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return err
	}
	if userCount == 0 {
		return errors.New("user does not exist")
	}

	result := db.Where("id = ?", resourceId).Delete(&models.Permission{})
	if result.Error != nil {
		logger.GetLogger("quiver").Errorf("error deleting permission for id %d: %v", resourceId, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.GetLogger("quiver").Warnf("permission delete failed %s", "permission not found")
		return errors.New("permission not found")
	}

	return nil
}
