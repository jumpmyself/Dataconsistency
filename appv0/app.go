package appv0

import (
	"Dataconsistency/appv0/db"
	"Dataconsistency/appv0/logic"
	"fmt"
	"net/http"
)

func Run() {
	//1.加载配置文件、链接数据库
	db.NewDb()
	db.NewRdb()

	//2.注册路由
	http.HandleFunc("/get_name", logic.GetInfo)

	//3.启动http服务
	fmt.Println("server started at http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		panic(fmt.Sprintf("fail to start server: %v", err))
	}

}
