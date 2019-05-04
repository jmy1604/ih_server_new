package main

import (
	"fmt"
	"ih_server_new/libs/log"
	"ih_server_new/src/db_gen/login_db"
	"ih_server_new/src/server_config"
	"ih_server_new/src/share_data"

	"github.com/huoshan017/mysql-go/manager"
)

var config server_config.LoginServerConfig
var shutingdown bool
var db mysql_manager.DB
var tb_mgr *login_db.TablesManager
var account_table *login_db.T_Account_Table
var ban_player_table *login_db.T_Ban_Player_Table
var server_list share_data.ServerList

func main() {
	defer func() {
		log.Trace("关闭服务器", nil)
		if err := recover(); err != nil {
			log.Stack(err)
		}
		if server != nil {
			server.Shutdown()
		}
		log.Close()
	}()

	if !server_config.ServerConfigLoad("login_server.json", &config) {
		fmt.Printf("载入LoginServer配置失败")
		return
	}

	log.Trace("配置:服务器ID", config.ServerId)
	log.Trace("配置:服务器名称", config.ServerName)
	log.Trace("配置:服务器地址(对Client)", config.ListenClientIP)
	log.Trace("配置:服务器地址(对Game)", config.ListenGameIP)

	server_list.ReadConfig(server_config.GetConfPathFile("server_list.json"))

	if !global_config_load() {
		log.Error("global_config_load failed !")
		return
	}

	db_define := server_config.GetDBDefineFile("login_db.json")
	if !db.LoadConfig(db_define) {
		log.Error("載入數據庫定義%v配置失敗", db_define)
		return
	}

	log.Trace("连接数据库", config.MYSQL_NAME, log.Property{"地址", config.MYSQL_IP})
	if !db.Connect(config.MYSQL_IP, config.MYSQL_ACCOUNT, config.MYSQL_PWD, config.MYSQL_NAME) {
		log.Error("连接数据库(dbname: %v, ip: %v, user: %v, password: %v)失败 ", config.MYSQL_NAME, config.MYSQL_IP, config.MYSQL_ACCOUNT, config.MYSQL_PWD)
		return
	}

	db.Run()

	tb_mgr = login_db.NewTablesManager(&db)
	account_table = tb_mgr.Get_T_Account_Table()
	ban_player_table = tb_mgr.Get_T_Ban_Player_Table()

	if !signal_mgr.Init() {
		log.Error("signal_mgr init failed")
		return
	}

	server = new(LoginServer)
	if !server.Init() {
		return
	}

	if signal_mgr.IfClosing() {
		return
	}

	if !hall_agent_manager.Init() {
		return
	}

	center_conn.Init()
	go center_conn.Start()

	err := hall_agent_manager.Start()
	if err != nil {
		return
	}

	server.Start(config.UseHttps)
}
