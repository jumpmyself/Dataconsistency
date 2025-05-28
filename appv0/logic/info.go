package logic

import (
	"Dataconsistency/appv0/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetInfo 缓存回溯方法查询用户信息
func GetInfo(w http.ResponseWriter, r *http.Request) {
	//必须是get请求（go自带的http无法识别不同请求）
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 GET 请求", http.StatusMethodNotAllowed)
		return
	}
	//1.查看缓存是否存在
	userID := 1 // 应该从请求参数中获取真实ID
	key := fmt.Sprintf("user_%d", userID)
	cache, err := db.Rdb.Get(context.Background(), key).Result()
	if err == nil {
		w.Header().Set("Cache-Control", "application/json")
		response := map[string]interface{}{
			"source": "cache",
			"data":   cache,
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
		return
	}
	//2.如果不存在查询数据库
	t := db.Info{}
	ret, err := t.GetInfo(1)
	if err != nil {
		return
	}

	//3.如果不存在设置ret
	if ret != nil && ret.ID > 0 {
		if err := db.Rdb.Set(context.Background(), key, fmt.Sprint(ret), time.Second*5).Err(); err != nil {
			fmt.Printf("缓存设置失败:%v\n", err)
		}
	} else {
		//处理数据不存在的情况（防止缓存穿透）
		// 缓存空值并设置较短过期时间
		_ = db.Rdb.Set(context.Background(), key, "null", 30*time.Second).Err()
	}

	//4.返回结果
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"source": "database",
		"data":   ret,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
