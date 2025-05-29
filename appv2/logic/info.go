package logic

import (
	"Dataconsistency/appv2/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func setCache(key string) error {
	return db.Rdb.PExpire(context.TODO(), key, time.Millisecond*1).Err()
}

// SetInfoW1  延时双删（问题在于两次删除中间要间隔多久）
func SetInfoW1(w http.ResponseWriter, r *http.Request) {
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
	key := fmt.Sprintf("user_%d", userID)
	err = setCache(key)
	if err != nil {
		log.Printf("延时双删第一次删除失败:key=%s,err=%v", key, err)
		return
	}

	info := db.Info{}
	err = info.Save(userID, UserName)
	if err != nil {
		http.Error(w, `{"error": "数据库写入失败"}`, http.StatusInternalServerError)
		return
	}
	time.Sleep(time.Microsecond * 10)
	err = setCache(key)
	if err != nil {
		http.Error(w, `{"error": "延时双删第二次删除失败"}`, http.StatusInternalServerError)
		return
	}

	// 向客户端发送响应
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusOK,
		"id":   userID,
		"data": "延时双删完成",
	})
}
