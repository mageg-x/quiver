package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
	"quiver/utils"
	"strconv"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}
func GetUserID(c *fiber.Ctx) (uint64, error) {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		logger.GetLogger("quiver").Errorf("invalid user_id")
		return 0, utils.BadRequest(c, "invalid user_id")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", err)
		return 0, utils.BadRequest(c, "invalid user_id format")
	}

	if userID <= 0 {
		logger.GetLogger("quiver").Errorf("invalid user_id: %v", userID)
		return 0, utils.BadRequest(c, "invalid user_id")
	}

	return uint64(userID), nil
}

func (s *UserService) CreateUser(env string, user *models.User) error {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return errors.New("db not initialized")
	}

	if len(user.UserName) == 0 || len(user.Password) == 0 {
		logger.GetLogger("quiver").Errorf("username or password cannot be empty")
		return errors.New("username or password cannot be empty")
	}

	var existing models.User
	if err := db.Where("username = ? ", user.UserName).First(&existing).Error; err == nil {
		if existing.UserName == user.UserName {
			logger.GetLogger("quiver").Errorf("username %s already exists", user.UserName)
			return errors.New("username already exists")
		}
		logger.GetLogger("quiver").Errorf("user %s already exists", user.UserName)
		return errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.GetLogger("quiver").Errorf("user create %s failed %s", user.UserName, err.Error())
		return err
	}

	user.CreateTime = db.NowFunc()
	user.UpdateTime = db.NowFunc()
	return db.Create(user).Error
}

func (s *UserService) GetUser(env string, userID uint64) (*models.User, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger("quiver").Errorf("user not found for id %d", userID)
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) ListUser(env string, page, size int) ([]models.User, int64, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, 0, errors.New("db not initialized")
	}

	var users []models.User
	var total int64

	offset := (page - 1) * size

	if err := db.Model(&models.User{}).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count users: %v", err)
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to query users: %v", err)
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) UpdateUser(env string, userID uint64, updates map[string]interface{}) (*models.User, error) {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return nil, errors.New("db not initialized")
	}

	var user models.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger("quiver").Errorf("user not found for user_id %d", userID)
			return nil, errors.New("user not found")
		}
		logger.GetLogger("quiver").Errorf("error finding user for user_id %d: %v", userID, err)
		return nil, err
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		logger.GetLogger("quiver").Errorf("error updating user for user_id %d: %v", userID, err)
		return nil, err
	}

	return s.GetUser(env, userID)
}

func (s *UserService) DeleteUser(env string, userID uint64) error {
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return errors.New("db not initialized")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.User{}).Error; err != nil {
			logger.GetLogger("quiver").Errorf("delete user failed: %v", err)
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.Permission{}).Error; err != nil {
			logger.GetLogger("quiver").Errorf("delete user permissions failed: %v", err)
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.AccessKey{}).Error; err != nil {
			logger.GetLogger("quiver").Errorf("delete user access keys failed: %v", err)
			return err
		}
		return nil
	})
}
