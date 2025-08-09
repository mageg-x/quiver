package database

import (
	"fmt"
	gormlogger "gorm.io/gorm/logger"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"quiver/config"
	"quiver/logger"
	"time"
)

var (
	DBMap         = make(map[string]*gorm.DB)
	dbMapMutex    sync.RWMutex // 使用 RWMutex 来支持并发读和独占写
	connectingMux sync.Mutex   //  一次只能是一个连接在进行
)

// Connect 连接数据库，根据传入的环境加载对应配置
func Connect(env string) (*gorm.DB, error) {
	dbConfig, err := config.GetDatabaseConfig(env)
	if err != nil {
		logger.GetLogger("quiver").Errorf("Failed to get database config for environment '%s': %v", env, err)
		return nil, fmt.Errorf("failed to get database config: %w", err)
	}

	logger.GetLogger("quiver").Infof("Using database config: %+v for evn: %s", dbConfig, env)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})

	if err != nil {
		logger.GetLogger("quiver").Errorf("Failed to connect to database %+v for env : %s error: %v",
			dbConfig, env, err)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.GetLogger("quiver").Errorf("Failed to get raw database instance:%+v for env : %s error: %v",
			dbConfig, env, err)
		return nil, fmt.Errorf("failed to get raw database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.GetLogger("quiver").Infof("Database connected successfully %+v for env : %s", dbConfig, env)

	return db, nil
}

// GetDB 获取指定环境下的数据库实例（自动连接如果未连接）
func GetDB(env string) *gorm.DB {
	// 首先尝试读锁检查是否存在数据库连接
	dbMapMutex.Lock()
	defer dbMapMutex.Unlock()
	db, exists := DBMap[env]
	//dbMapMutex.Unlock()

	if exists && db != nil {
		//logger.GetLogger("quiver").Infof("Using existing database %v connection for env: %s", db, env)
		return db
	} else {
		logger.GetLogger("quiver").Warnf("Database connection not found for env: %s, connecting...", env)
	}

	// 创建新的数据库连接, 比较耗时，放在异步处理
	logger.GetLogger("quiver").Infof("goroutine started for env: %s", env)
	if connectingMux.TryLock() {
		defer connectingMux.Unlock()
		db, err := Connect(env)
		if err != nil || db == nil {
			logger.GetLogger("quiver").Errorf("Failed to connect to database for env: %s, error: %v", env, err)
			return nil
		}
		logger.GetLogger("quiver").Infof("Success to connect to database for env: %s", env)
		DBMap[env] = db
		logger.GetLogger("quiver").Infof("save db %v for env: %s", db, env)
		return db
	}

	return nil
}
