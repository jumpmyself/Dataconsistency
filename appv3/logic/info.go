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

// GetInfoR0 缓存回溯后门 - 强制刷新指定用户的缓存
func GetInfoR0(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 GET 请求", http.StatusMethodNotAllowed)
		return
	}

	// 身份验证 - 仅限管理端使用
	if !isAdminRequest(r) {
		http.Error(w, `{"error": "未授权访问"}`, http.StatusUnauthorized)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	cacheParam := r.URL.Query().Get("cache")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("user_%d", userID)

	// 直接查询数据库（跳过缓存）
	info := db.Info{}
	data, err := info.GetInfo(userID)
	if err != nil {
		http.Error(w, `{"error": "数据库查询失败"}`, http.StatusInternalServerError)
		return
	}

	// 刷新缓存（无论cache参数如何）
	if data != nil && data.ID > 0 {
		dataBytes, err := json.Marshal(data)
		if err == nil {
			// 统一设置较长的缓存时间（区别于常规业务缓存）
			_ = db.Rdb.Set(context.Background(), key, string(dataBytes), 5*time.Minute).Err()
		}
	} else {
		// 防止缓存穿透
		_ = db.Rdb.Set(context.Background(), key, "null", 5*time.Minute).Err()
	}

	// 标识这是强制更新的结果
	source := "强制更新缓存（数据库查询）"
	if cacheParam == "true" {
		source = "缓存强制刷新完成"
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"source": source,
		"data":   data,
	})
}

// 简单管理端验证（应根据实际安全要求加强）
func isAdminRequest(r *http.Request) bool {
	token := r.Header.Get("X-Admin-Token")
	return token == "secure_admin_token_123" // 应替换为实际验证逻辑
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
