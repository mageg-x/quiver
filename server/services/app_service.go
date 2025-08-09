package services

import (
	"errors"
	"quiver/database"
	"quiver/logger"
	"quiver/models"

	"gorm.io/gorm"
)

// AppService 应用服务
type AppService struct{}

// NewAppService 创建应用服务实例
func NewAppService() *AppService {
	return &AppService{}
}

// CreateApp 创建应用
func (s *AppService) CreateApp(env string, app *models.App) error {
	logger.GetLogger("quiver").Infof("app create %+v for env %s", app, env)
	db := database.GetDB(env)
	if db == nil {
		logger.GetLogger("quiver").Errorf("db is nil for env %s", env)
		return errors.New("db not initialized")
	}

	// 检查应用ID是否已存在
	var existingApp models.App
	err := db.Where("app_name = ?", app.AppName).First(&existingApp).Error
	if err == nil {
		logger.GetLogger("quiver").Errorf("app_name %s already exists", app.AppName)
		return errors.New("app_name already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.GetLogger("quiver").Errorf("app create %s failed %s", app.AppName, err.Error())
		return err
	}

	logger.GetLogger("quiver").Infof("begin create  %s in db", app.AppName)
	// 创建应用
	return db.Create(app).Error
}

// GetApp 根据ID获取应用
func (s *AppService) GetApp(env string, appName string) (*models.App, error) {
	db := database.GetDB(env)

	var app models.App
	err := db.Where("app_name = ?", appName).First(&app).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("app get %s failed %s", appName, err.Error())
		return nil, err
	}

	return &app, nil
}

// ListApp 获取所有应用
func (s *AppService) ListApp(env string, page, size int) ([]models.App, int64, error) {
	db := database.GetDB(env)

	var apps []models.App
	var total int64

	// 获取总数
	if err := db.Model(&models.App{}).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count apps num %s", err.Error())
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// UpdateApp 更新应用
func (s *AppService) UpdateApp(env string, appName string, updates map[string]interface{}) (*models.App, error) {
	db := database.GetDB(env)

	// 获取当前版本号
	app, err := s.GetApp(env, appName)
	if err != nil {
		logger.GetLogger("quiver").Errorf("app get %s failed %s", appName, err.Error())
		return nil, errors.New("app_name not exist")
	}

	ver := app.Ver
	updates["ver"] = ver + 1

	// 使用乐观锁更新
	result := db.Model(&models.App{}).Where("app_name = ? AND ver = ?", appName, ver).Updates(updates)
	if result.Error != nil {
		logger.GetLogger("quiver").Errorf("app update %s failed %s", appName, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		logger.GetLogger("quiver").Errorf("app update %s failed %s", appName, "app not found or version conflict")
		return nil, errors.New("app not found or version conflict")
	}

	return s.GetApp(env, appName)
}

// DeleteApp 删除应用
func (s *AppService) DeleteApp(env string, appName string) error {
	db := database.GetDB(env)

	// 先查询app 是否存在
	app, err := s.GetApp(env, appName)
	if err != nil || app == nil {
		logger.GetLogger("quiver").Errorf("app get %s failed %s", appName, err)
		return errors.New("app_name not exist")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 更新  Item 记录的 deleted 字段
		err := BulkMarkDeleted[*models.Item](tx, map[string]interface{}{"app_id": app.AppID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items deleted for app %s failed %s", appName, err)
			return err
		}

		// 更新与  ItemRelease 记录的 deleted 字段
		err = BulkMarkDeleted[*models.ItemRelease](tx, map[string]interface{}{"app_id": app.AppID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items release deleted for app %s failed %s", appName, err)
			return err
		}

		// 删除应用相关的所有数据
		result := tx.Where("id = ?", app.AppID).Delete(&models.App{})
		if result.Error != nil {
			logger.GetLogger("quiver").Errorf("app delete %s failed %s", appName, result.Error)
			return result.Error
		}

		if result.RowsAffected == 0 {
			logger.GetLogger("quiver").Errorf("app delete %s failed %s", appName, "app not found")
			return errors.New("app_name not exist")
		}

		return nil
	})
}
