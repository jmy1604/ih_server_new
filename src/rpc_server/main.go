package main

import (
	"fmt"
	"ih_server_new/libs/log"
	"ih_server_new/src/db_gen/rpc_db"
	"ih_server_new/src/server_config"
	"ih_server_new/src/share_data"
	"time"

	"github.com/huoshan017/mysql-go/manager"
)

var config server_config.RpcServerConfig
var server_list share_data.ServerList
var db mysql_manager.DB
var google_pay_table *rpc_db.T_Google_Pay_Table
var apple_pay_table *rpc_db.T_Apple_Pay_Table

func main() {
	defer func() {
		log.Trace("关闭rpc_server服务器", nil)
		if err := recover(); err != nil {
			log.Stack(err)
		}
		server.Shutdown()
		time.Sleep(time.Second * 5)
		log.Close()
	}()

	if !server_config.ServerConfigLoad("rpc_server.json", &config) {
		fmt.Printf("载入RPC Server配置失败")
		return
	}

	if !server_list.ReadConfig(server_config.GetConfPathFile("server_list.json")) {
		return
	}

	var err error
	if config.MYSQL_NAME != "" {
		if !db.LoadConfig(server_config.GetDBDefineFile("rpc_db.json")) {
			fmt.Printf("载入db定义配置失败")
			return
		}

		log.Trace("连接数据库", config.MYSQL_NAME, log.Property{"地址", config.MYSQL_IP})
		if !db.Connect(config.MYSQL_IP, config.MYSQL_ACCOUNT, config.MYSQL_PWD, config.MYSQL_NAME) {
			log.Error("连接数据库失败 %v", err)
			return
		}

		log.Trace("连接数据库成功", nil)
		db.Run()

		tb_mgr := rpc_db.NewTablesManager(&db)
		google_pay_table = tb_mgr.Get_T_Google_Pay_Table()
		apple_pay_table = tb_mgr.Get_T_Apple_Pay_Table()

		if !signal_mgr.Init() {
			log.Error("signal_mgr init failed")
			return
		}
	}

	err = server.Init()
	if err != nil {
		log.Error("RPC Server init error[%v]", err.Error())
		return
	}

	if config.MYSQL_NAME != "" {
		if signal_mgr.IfClosing() {
			return
		}
	}

	if config.GmServerUseHttps {
		go gm_service.StartHttps(server_config.GetConfPathFile("server.crt"), server_config.GetConfPathFile("server.key"))
	} else {
		go gm_service.StartHttp()
	}

	fmt.Println("启动服务...")

	server.Start()

	fmt.Println("服务已停止!")
}
