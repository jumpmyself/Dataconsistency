package schedule

import (
	"Dataconsistency/appv3/db"
	"context"
	"fmt"
	"time"
)

func Run() {
	go updateCache()
	go updateCache15()
}

// 异步更新
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

// 异步更新多个
func updateCachev1() {
	t := time.NewTicker(5 * time.Second)
	for range t.C {
		fmt.Println(time.Now().Unix())
		info := db.Info{}
		data1, _ := info.GetInfo(1)
		data2, _ := info.GetInfo(2)
		data3, _ := info.GetInfo(3)
		data4, _ := info.GetInfo(4)
		data5, _ := info.GetInfo(5)

		key1 := fmt.Sprintf("user_%v", 1)
		key2 := fmt.Sprintf("user_%v", 2)
		key3 := fmt.Sprintf("user_%v", 3)
		key4 := fmt.Sprintf("user_%v", 4)
		key5 := fmt.Sprintf("user_%v", 5)
		_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
		fmt.Println("缓存更新完成")
		_ = db.Rdb.Set(context.Background(), key2, fmt.Sprint(data2), time.Second*5).Err()
		fmt.Println("缓存更新完成")
		_ = db.Rdb.Set(context.Background(), key3, fmt.Sprint(data3), time.Second*5).Err()
		fmt.Println("缓存更新完成")
		_ = db.Rdb.Set(context.Background(), key4, fmt.Sprint(data4), time.Second*5).Err()
		fmt.Println("缓存更新完成")
		_ = db.Rdb.Set(context.Background(), key5, fmt.Sprint(data5), time.Second*5).Err()
		fmt.Println("缓存更新完成")

	}
}

// 使用协程的方式异步更新多个
func updateCachev2() {
	t := time.NewTicker(5 * time.Second)
	for range t.C {
		fmt.Println(time.Now().Unix())
		info := db.Info{}

		go func() {
			data1, _ := info.GetInfo(1)
			key1 := fmt.Sprintf("user_%v", 1)
			_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
		}()
		fmt.Println(time.Now().Unix())
		go func() {
			data1, _ := info.GetInfo(2)
			key1 := fmt.Sprintf("user_%v", 2)
			_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
		}()
		fmt.Println(time.Now().Unix())
		go func() {
			data1, _ := info.GetInfo(3)
			key1 := fmt.Sprintf("user_%v", 3)
			_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
		}()
		fmt.Println(time.Now().Unix())
		go func() {
			data1, _ := info.GetInfo(4)
			key1 := fmt.Sprintf("user_%v", 4)
			_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
		}()
		fmt.Println(time.Now().Unix())
		go func() {
			data1, _ := info.GetInfo(5)
			key1 := fmt.Sprintf("user_%v", 5)
			_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
		}()
		fmt.Println(time.Now().Unix())
	}
}

// 使用select的方式监听ticker
func updateCachev3() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			fmt.Println(time.Now().Unix())
			info := db.Info{}

			go func() {
				data1, _ := info.GetInfo(1)
				key1 := fmt.Sprintf("user_%v", 1)
				_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
			}()
			fmt.Println(time.Now().Unix())
			go func() {
				data1, _ := info.GetInfo(2)
				key1 := fmt.Sprintf("user_%v", 2)
				_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
			}()
			fmt.Println(time.Now().Unix())
			go func() {
				data1, _ := info.GetInfo(3)
				key1 := fmt.Sprintf("user_%v", 3)
				_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
			}()
			fmt.Println(time.Now().Unix())
			go func() {
				data1, _ := info.GetInfo(4)
				key1 := fmt.Sprintf("user_%v", 4)
				_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
			}()
			fmt.Println(time.Now().Unix())
			go func() {
				data1, _ := info.GetInfo(5)
				key1 := fmt.Sprintf("user_%v", 5)
				_ = db.Rdb.Set(context.Background(), key1, fmt.Sprint(data1), time.Second*5).Err()
			}()
			fmt.Println(time.Now().Unix())
		}
	}
}

func updateCache15() {
	//更新缓存数据
	key := fmt.Sprintf("book_%v", 1)
	t := time.NewTicker(15 * time.Second)
	for range t.C {
		info := db.Info{}
		data, _ := info.GetInfo(1)

		_ = db.Rdb.Set(context.Background(), key, fmt.Sprint(data), time.Second*5).Err()
		fmt.Println("缓存更新完成")
	}

}

func updateCachev4() {
	t1 := time.NewTicker(5 * time.Second)
	t2 := time.NewTicker(15 * time.Second)
	t3 := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-t1.C:
			fmt.Println("这是5s的任务在启动")
			go func() {

			}()
		case <-t2.C:
			fmt.Println("这是15s的任务在启动")
			go func() {

			}()
		case <-t3.C:
			fmt.Println("这是30s的任务在启动")
			go func() {

			}()
		}

	}

}
