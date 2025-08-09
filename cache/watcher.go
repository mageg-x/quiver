package cache

import (
	"errors"
	"fmt"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
	"quiver/utils"
	"sync"
	"time"
)

type CacheableModel interface {
	GetID() uint64
	GetUpdateTime() time.Time
	CacheKey(env string) string // 返回要删除的缓存 key
	TableName() string
}

type Watcher struct {
	watchTables map[string]WatchTable
	mu          sync.Mutex
}
type OnKeyUpdateCallback func(evn, key string)
type WatchTable struct {
	Env           string
	Name          string
	UpdateTime    time.Time
	Interval      time.Duration
	LastCheckTime time.Time
	Running       bool
	KeyUpdateCB   OnKeyUpdateCallback
}

var (
	globalWatcher *Watcher
)

func (w *Watcher) watch() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			utils.WithTryLock(&w.mu, func() {
				for key := range w.watchTables {
					item := w.watchTables[key]
					if item.LastCheckTime.Add(item.Interval).After(time.Now()) || item.Running {
						continue
					}

					// 标记为正在运行（避免重复触发）
					item.Running = true
					item.LastCheckTime = time.Now()
					w.watchTables[key] = item // 写回 map

					// 启动异步检查
					go func(k string, table *WatchTable) {
						var newUpdateTime time.Time
						defer utils.WithLock(&w.mu, func() {
							table.Running = false
							table.UpdateTime = newUpdateTime
							w.watchTables[k] = *table
						})
						switch table.Name {
						case "accesskey":
							newUpdateTime = checkUpdate[*models.AccessKey](table.Env, table)
						case "permission":
							newUpdateTime = checkUpdate[*models.Permission](table.Env, table)
						case "namespace_release":
							newUpdateTime = checkUpdate[*models.NamespaceRelease](table.Env, table)
						}
					}(key, &item)
				}
			})
		}
	}()
}
func checkUpdate[T CacheableModel](env string, table *WatchTable) time.Time {
	db := database.GetDB(env)
	if db == nil || table == nil {
		logger.GetLogger("quiver").Errorf("db is nil")
		return time.Unix(0, 0)
	}

	const batchSize = 1000
	var lastID uint64
	updateTime := table.UpdateTime

	for {
		var items []T
		err := db.Where("update_time > ? AND id > ?", updateTime, lastID).
			Order("update_time ASC, id ASC").
			Limit(batchSize).
			Find(&items).Error

		if err != nil {
			//logger.GetLogger("quiver").Errorf("query error: %v", err)
			return updateTime
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			lastID = max(item.GetID(), lastID)
			if item.GetUpdateTime().After(updateTime) {
				updateTime = item.GetUpdateTime()
			}
			key := item.CacheKey(env)
			logger.GetLogger("quiver").Infof("ready to delete key: %s", key)
			if err := Delete(key); err != nil {
				logger.GetLogger("quiver").Errorf("delete %s error: %v", key, err)
			}

			if table != nil && table.KeyUpdateCB != nil {
				go func() {
					table.KeyUpdateCB(env, key)
				}()
			}
		}

		if len(items) < batchSize {
			break
		}
	}

	// 更新传入的指针
	return updateTime
}

// StartWatch 全局方法封装
func startWatch() {
	globalWatcher = &Watcher{
		watchTables: make(map[string]WatchTable),
	}
	globalWatcher.watch()
}

func AddWatch(items []WatchTable) error {
	if globalWatcher == nil {
		return errors.New("watcher not initialized")
	}

	utils.WithLock(&globalWatcher.mu, func() {
		for _, item := range items {
			key := fmt.Sprintf("%s:%s", item.Env, item.Name)
			globalWatcher.watchTables[key] = WatchTable{
				Env:         item.Env,
				Name:        item.Name,
				UpdateTime:  time.Now(),
				Interval:    item.Interval,
				KeyUpdateCB: item.KeyUpdateCB,
			}
		}
	})
	return nil
}
