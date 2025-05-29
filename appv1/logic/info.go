package logic

import (
	"Dataconsistency/appv1/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("user_%d", userID)
	cache, err := db.Rdb.Get(context.Background(), key).Result()
	if err == nil {
		// 缓存命中：返回缓存数据并明确标注来源
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
	//2.如果不存在查询数据库
	t := db.Info{}
	ret, err := t.GetInfo(userID)
	if err != nil {
		return
	}

	//3.并设置ret
	if ret != nil && ret.ID > 0 {
		// 序列化数据为JSON
		dataBytes, _ := json.Marshal(ret)
		_ = db.Rdb.Set(context.Background(), key, string(dataBytes), 10*time.Second).Err()
	} else {
		// 防止缓存穿透：存储空值（短时间）
		_ = db.Rdb.Set(context.Background(), key, "null", 30*time.Second).Err()
	}

	//4.返回结果
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"source": "database",
		"data":   ret,
	})

}

// SetInfoW2  先写数据库再删缓存（删缓存的方式、直接删除）
// TODO 直接删除大key缓存会阻塞redis
func SetInfoW2(w http.ResponseWriter, r *http.Request) {
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
	//
	info := db.Info{}
	err = info.Save(userID, UserName)
	if err != nil {
		http.Error(w, `{"error": "数据库写入失败"}`, http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("user_%d", userID)
	if err := db.Rdb.Del(context.Background(), key).Err(); err != nil {
		http.Error(w, `{"error": "Redis 删除失败"}`, http.StatusInternalServerError)
		return
	}

	// 向客户端发送响应
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusOK,
		"id":   userID,
		"data": "先写数据库再直接删除缓存",
	})
}

// SetInfoW3  先写数据库再删缓存（删缓存的方式、采用设置过期时间的方式）
func SetInfoW3(w http.ResponseWriter, r *http.Request) {
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
	info := db.Info{}
	err = info.Save(userID, UserName)
	if err != nil {
		http.Error(w, `{"error": "数据库写入失败"}`, http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("user_%d", userID)
	if err := db.Rdb.PExpire(context.Background(), key, time.Millisecond*200).Err(); err != nil {
		http.Error(w, `{"error": "Redis expire删除失败"}`, http.StatusInternalServerError)
		return
	}

	// 向客户端发送响应
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusOK,
		"id":   userID,
		"data": "先写数据库、再设置过期时间删除redis",
	})
}
