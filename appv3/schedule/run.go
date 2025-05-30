package schedule

import (
	"Dataconsistency/appv3/db"
	"context"
	"fmt"
	"time"
)

func Run() {
	go updateCache()
}

func updateCache() {
	//更新缓存数据
	key := fmt.Sprintf("book_%v", 1)
	t := time.NewTicker(5 * time.Second)
	for range t.C {
		info := db.Info{}
		data, _ := info.GetInfo(1)

		_ = db.Rdb.Set(context.Background(), key, fmt.Sprint(data), time.Second*5).Err()
		fmt.Println("缓存更新完成")
	}
}
