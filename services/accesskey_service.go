package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
)

type AccessKeyService struct{}

func NewAccessKeyService() *AccessKeyService {
	return &AccessKeyService{}
}

func (s *AccessKeyService) generateKeys() (string, string, error) {
	ak := make([]byte, 16)
	sk := make([]byte, 32)
	if _, err := rand.Read(ak); err != nil {
		return "", "", err
	}
	if _, err := rand.Read(sk); err != nil {
		return "", "", err
	}
	return hex.EncodeToString(ak), hex.EncodeToString(sk), nil
}

func (s *AccessKeyService) CreateAccessKey(env string, userID uint64) (*models.AccessKey, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}
	if userID == 0 {
		logger.GetLogger("quiver").Errorf("userID is 0")
		return nil, errors.New("userID is 0")
	}

	// 检查用户是否存在
	var userCount int64
	err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return nil, err
	}
	if userCount == 0 {
		logger.GetLogger("quiver").Errorf("user does not exist for id %d", userID)
		return nil, errors.New("user does not exist")
	}

	ak, sk, err := s.generateKeys()
	if err != nil {
		logger.GetLogger("quiver").Errorf("error generating keys: %v", err)
		return nil, err
	}

	accessKey := &models.AccessKey{
		UserID:    userID,
		AccessKey: ak,
		SecretKey: sk,
	}

	if err := db.Create(accessKey).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error creating access key: %v, userID: %d", err, userID)
		return nil, err
	}

	// 注意：返回 secretKey，仅本次可见
	return accessKey, nil
}

func (s *AccessKeyService) GetAccessKey(env string, userID uint64, accessKey string) (*models.AccessKey, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}

	if userID == 0 {
		logger.GetLogger("quiver").Errorf("userID is 0")
		return nil, errors.New("userID is 0")
	}

	// 检查用户是否存在
	var userCount int64
	err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return nil, err
	}

	if userCount == 0 {
		logger.GetLogger("quiver").Errorf("user does not exist for id %d", userID)
		return nil, errors.New("user does not exist")
	}

	if accessKey == "" {
		logger.GetLogger("quiver").Errorf("access key cannot be empty")
		return nil, errors.New("access key cannot be empty")
	}

	var ak models.AccessKey
	if err := db.Where("access_key = ?", accessKey).First(&ak).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error querying access key: %v", err)
		return nil, err
	}
	return &ak, nil
}

func (s *AccessKeyService) ListAccessKey(env string, userID uint64) ([]models.AccessKey, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}

	if userID == 0 {
		logger.GetLogger("quiver").Errorf("userID is 0")
		return nil, errors.New("userID is 0")
	}

	// 检查用户是否存在
	var userCount int64
	err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return nil, err
	}

	if userCount == 0 {
		logger.GetLogger("quiver").Errorf("user does not exist for id %d", userID)
		return nil, errors.New("user does not exist")
	}

	var aks []models.AccessKey
	if err := db.Where("user_id = ?", userID).Find(&aks).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error listing access keys for user %d: %v", userID, err)
		return nil, err
	}
	// 不返回 secret_key
	for i := range aks {
		aks[i].SecretKey = ""
	}
	return aks, nil
}

func (s *AccessKeyService) DeleteAccessKey(env string, userID uint64, accessKey string) error {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return errors.New("db not initialized")
	}

	if userID == 0 {
		logger.GetLogger("quiver").Errorf("userID is 0")
		return errors.New("userID is 0")
	}

	// 检查用户是否存在
	var userCount int64
	err := db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("error querying user existence: %v", err)
		return err
	}

	if userCount == 0 {
		logger.GetLogger("quiver").Errorf("user does not exist for id %d", userID)
		return errors.New("user does not exist")
	}

	if accessKey == "" {
		logger.GetLogger("quiver").Errorf("access key cannot be empty")
		return errors.New("access key cannot be empty")
	}

	result := db.Where("access_key = ?", accessKey).Delete(&models.AccessKey{})
	if result.Error != nil {
		logger.GetLogger("quiver").Errorf("error deleting accesskey: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.GetLogger("quiver").Warnf("accesskey delete failed %s", "accesskey not found")
		return errors.New("accesskey not found")
	}
	return nil
}
