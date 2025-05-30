package appv3

import (
	"Dataconsistency/appv3/db"
	"Dataconsistency/appv3/logic"
	"Dataconsistency/appv3/schedule"
	"fmt"
	"net/http"
)

func Run() {
	//1.加载配置文件、链接数据库
	db.NewDb()
	db.NewRdb()
	//2.启动缓存更新worker
	schedule.Run()
	//3.注册路由
	http.HandleFunc("/get_name", logic.GetInfo)
	http.HandleFunc("/set_name_w1", logic.SetInfoW0)

	//3.启动http服务
	fmt.Println("server started at http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		panic(fmt.Sprintf("fail to start server: %v", err))
	}

}
