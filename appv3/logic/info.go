package logic

import (
	"Dataconsistency/appv3/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// GetInfo 缓存回溯方法查询用户信息
func GetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 GET 请求", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("user_%d", userID)
	cache, err := db.Rdb.Get(context.Background(), key).Result()

	// 处理缓存命中
	if err == nil {
		// 处理空值缓存（防止缓存穿透）
		if cache == "null" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"source": "redis",
				"data":   nil,
			})
			return
		}

		// 处理有效缓存
		var cachedData db.Info
		if err := json.Unmarshal([]byte(cache), &cachedData); err == nil {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"source": "redis",
				"data":   cachedData,
			})
			return
		}
	}

	// 缓存未命中，查询数据库
	t := db.Info{}
	ret, err := t.GetInfo(userID)
	if err != nil {
		http.Error(w, `{"error": "数据库查询失败"}`, http.StatusInternalServerError)
		return
	}

	// 设置缓存
	if ret != nil && ret.ID > 0 {
		dataBytes, _ := json.Marshal(ret)
		_ = db.Rdb.Set(context.Background(), key, string(dataBytes), 30*time.Second).Err()
	} else {
		// 防止缓存穿透：存储空值
		_ = db.Rdb.Set(context.Background(), key, "null", 30*time.Second).Err()
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"source": "database",
		"data":   ret,
	})
}

// SetInfoW0 只更新数据库并触发异步缓存更新
func SetInfoW0(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 Get 请求", http.StatusMethodNotAllowed)
		return
	}

	UserName := r.URL.Query().Get("name")
	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	// 更新数据库
	info := db.Info{}
	err = info.Save(userID, UserName)
	if err != nil {
		http.Error(w, `{"error": "数据库写入失败"}`, http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusOK,
		"id":   userID,
		"data": "数据库更新完成，缓存异步更新已触发",
	})
}
