package cache

import (
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"path/filepath"
	"quiver/logger"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto"
	"github.com/google/uuid"
)

// QCache 是一个两级缓存：内存 + 磁盘（临时，不持久化）
type QCache struct {
	memory *ristretto.Cache
	disk   *badger.DB
	dir    string
	mu     sync.RWMutex
	closed bool
}

// Options 包含创建 QCache 所需的配置选项。
type Options struct {
	MaxMemCost int64  // 内存缓存的最大开销
	DiskDir    string // 磁盘缓存的目录路径
}

type CachedItem struct {
	Value    []byte `json:"value"`
	ExpireAt int64  `json:"expireAt"` // Unix 时间戳（秒）
}

var (
	globalCache *QCache
	initOnce    sync.Once
)

func Create(option Options) (*QCache, error) {
	memoryCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: option.MaxMemCost / 10,
		MaxCost:     option.MaxMemCost,
		BufferItems: 64,
	})

	if err != nil {
		logger.GetLogger("quiver").Errorf("new ristretto cache failed: %v", err)
		return nil, err
	} else {
		logger.GetLogger("quiver").Infof("new mem cache success")
	}

	cache := QCache{
		memory: memoryCache,
		disk:   nil,
		dir:    "",
	}

	if option.DiskDir != "" {
		// 初始化 badger
		cache.dir = filepath.Join(option.DiskDir, uuid.New().String())
		if err := os.MkdirAll(cache.dir, 0755); err != nil {
			logger.GetLogger("quiver").Errorf("create cahce dir %s failed: %v", cache.dir, err)
			return nil, err
		}
		db, err := cache.ReOpen()
		if err != nil {
			logger.GetLogger("quiver").Errorf("open cache disk db %s failed: %v", cache.dir, err)
			return nil, err
		}
		cache.disk = db

		logger.GetLogger("quiver").Infof("new disk cache success: %s", cache.dir)
	}

	return &cache, nil
}

func (q *QCache) Get(key string) ([]byte, bool, error) {
	// 1. 先查内存
	if val, found := q.memory.Get(key); found {
		if item, ok := val.(*CachedItem); ok {
			if item.ExpireAt != 0 && time.Now().Unix() > item.ExpireAt {
				// 内存过期，删除
				q.memory.Del(key)
				_ = q.disk.Update(func(txn *badger.Txn) error {
					return txn.Delete([]byte(key))
				})
				return nil, false, nil
			}
			return item.Value, true, nil
		}
	}

	// 2. 内存未命中，查磁盘
	if q.disk == nil || q.disk.IsClosed() {
		return nil, false, nil
	}

	var item *CachedItem
	err := q.disk.View(func(txn *badger.Txn) error {
		itemObj, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return itemObj.Value(func(val []byte) error {
			// 反序列化
			item = &CachedItem{}
			err = msgpack.Unmarshal(val, item)
			if err != nil {
				// 数据已经被破坏,删除掉
				_ = q.Delete(key)
			}
			return err
		})
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	// 3. 检查磁盘读出的数据是否过期
	if item.ExpireAt != 0 && time.Now().Unix() > item.ExpireAt {
		// 磁盘数据也过期，删除
		if err := q.disk.Update(func(txn *badger.Txn) error {
			return txn.Delete([]byte(key))
		}); err != nil {
			logger.GetLogger("quiver").Warnf("delete expired key %s failed: %v", key, err)
			return nil, false, err
		}
		return nil, false, nil
	}

	// 4. 写回内存（带 TTL 重建）
	q.memory.Set(key, item, int64(len(item.Value)))

	return item.Value, true, nil
}

func (q *QCache) Set(key string, value []byte, ttl time.Duration) error {
	expireAt := time.Now().Add(ttl).Unix()
	if ttl == 0 {
		expireAt = 0
	}
	item := &CachedItem{
		Value:    value,
		ExpireAt: expireAt,
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	// 写入内存（带成本）
	q.memory.Set(key, item, int64(len(value)))

	if q.disk == nil || q.disk.IsClosed() {
		return nil
	}

	// 序列化
	data, err := msgpack.Marshal(item)
	if err != nil || len(data) == 0 {
		logger.GetLogger("quiver").Errorf("json marshal failed: %v", err)
		return err
	}

	// 写入磁盘（存储序列化后的字节）
	return q.disk.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}
func (q *QCache) Delete(key string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.memory.Del(key)

	if q.disk == nil || q.disk.IsClosed() {
		return nil
	}

	return q.disk.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (q *QCache) ReOpen() (*badger.DB, error) {
	if q.dir == "" {
		return nil, errors.New("cache dir is empty")
	}
	if q.disk != nil {
		if err := q.disk.Close(); err != nil {
			logger.GetLogger("quiver").Errorf("close badger cache failed: %v", err)
		}
		if err := os.RemoveAll(q.dir); err != nil {
			logger.GetLogger("quiver").Errorf("remove cache dir failed: %v", err)
		}

		q.closed = true
		q.disk = nil
	}

	opt := badger.DefaultOptions(q.dir).
		WithValueDir(q.dir).
		WithSyncWrites(false). // 提高性能
		WithLogger(nil)        // 可选：禁用日志

	db, err := badger.Open(opt)
	if err != nil {
		// 尝试删除缓存目录
		err := os.RemoveAll(q.dir)
		if err != nil {
			logger.GetLogger("quiver").Errorf("failed to remove cache dir %s: %v", q.dir, err)
			return nil, err
		}
		// 尝试换个目录 获取 q.dir 的 父目录
		q.dir = filepath.Join(filepath.Dir(q.dir), uuid.New().String())
		logger.GetLogger("quiver").Infof("trying to change disk cache dir form %s to %s : %s", opt.ValueDir, q.dir)
		opt.Dir = q.dir
		opt.ValueDir = q.dir
		db, err = badger.Open(opt)

		if err != nil {
			logger.GetLogger("quiver").Errorf("open badger failed after cleanup: %v", err)
			return nil, err
		}
	}
	q.disk = db
	return db, nil
}

func (q *QCache) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}

	q.memory.Close()

	if q.disk == nil || q.disk.IsClosed() {
		q.closed = true
		return nil
	}

	if err := q.disk.Close(); err != nil {
		logger.GetLogger("quiver").Errorf("close badger cache failed: %v", err)
	}
	if err := os.RemoveAll(q.dir); err != nil {
		logger.GetLogger("quiver").Errorf("remove cache dir failed: %v", err)
	}

	q.closed = true
	return nil
}

// Init 全局方法封装
func Init(option Options) error {
	var err error
	initOnce.Do(func() {
		globalCache, err = Create(option)
		startWatch()
	})
	return err
}

func Get(key string) ([]byte, bool, error) {
	if globalCache == nil {
		return nil, false, errors.New("cache not initialized")
	}
	return globalCache.Get(key)
}

func Set(key string, value []byte, ttl time.Duration) error {
	if globalCache == nil {
		return errors.New("cache not initialized")
	}
	return globalCache.Set(key, value, ttl)
}

func Close() error {
	if globalCache == nil {
		return nil
	}
	return globalCache.Close()
}

func Delete(key string) error {
	if globalCache == nil {
		return errors.New("cache not initialized")
	}
	return globalCache.Delete(key)
}
