package schedule

import (
	"Dataconsistency/appv3/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	updateQueue = make(chan int, 10000) // 缓存更新队列
	workersOnce sync.Once               // 确保worker只启动一次
)

// Run 启动缓存更新worker
func Run() {
	workersOnce.Do(func() {
		startCacheWorkers(10)
	})
}

// 启动缓存更新worker
func startCacheWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			log.Printf("缓存更新worker %d 已启动", workerID)
			for userID := range updateQueue {
				updateUserCache(userID)
			}
		}(i)
	}
}

// 更新单个用户缓存（纯异步更新）
func updateUserCache(userID int) {
	key := fmt.Sprintf("user_%d", userID)

	// 从数据库获取最新数据
	info := db.Info{}
	data, err := info.GetInfo(userID)
	if err != nil {
		log.Printf("异步更新缓存获取数据失败(userID=%d): %v", userID, err)
		return
	}

	// 设置缓存
	if data != nil && data.ID > 0 {
		dataBytes, _ := json.Marshal(data)
		err = db.Rdb.Set(context.Background(), key, string(dataBytes), 30*time.Second).Err()
	} else {
		// 防止缓存穿透：存储空值
		err = db.Rdb.Set(context.Background(), key, "null", 30*time.Second).Err()
	}

	if err != nil {
		log.Printf("异步更新Redis缓存失败(userID=%d): %v", userID, err)
		return
	}

	log.Printf("异步更新缓存成功(userID=%d)", userID)
}

// AddToQueue 添加用户ID到更新队列（新增的导出函数）
func AddToQueue(userID int) {
	select {
	case updateQueue <- userID:
		log.Printf("已加入缓存更新队列(userID=%d)", userID)
	default:
		log.Printf("更新队列已满，跳过缓存更新(userID=%d)", userID)
	}
}
