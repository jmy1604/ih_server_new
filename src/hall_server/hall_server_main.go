package main

import (
	"fmt"
	"ih_server_new/libs/log"
	"ih_server_new/src/db_gen/game_db"
	"ih_server_new/src/server_config"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"

	"github.com/huoshan017/mysql-go/manager"
)

var config server_config.GameServerConfig
var shutingdown bool
var dbc DBC
var db_mgr mysql_manager.DB
var battle_record_table *game_db.T_Battle_Record_Table
var global_table *game_db.T_Global_Table
var player_table *game_db.T_Player_Table
var ban_player_table *game_db.T_Ban_Player_Table

func after_center_match_conn() {
	if signal_mgr.IfClosing() {
		return
	}
}

func main() {
	defer func() {
		log.Trace("关闭服务器", nil)
		if err := recover(); err != nil {
			log.Stack(err)
			debug.PrintStack()
		}
		time.Sleep(time.Second * 5)
		hall_server.Shutdown()
	}()

	if !server_config.ServerConfigLoad("hall_server.json", &config) {
		fmt.Printf("载入GameServer配置失败")
		return
	}

	log.Trace("配置:服务器监听客户端地址", config.ListenClientInIP)
	log.Trace("配置:最大客户端连接数)", config.MaxClientConnections)
	log.Trace("连接数据库", config.MYSQL_NAME, log.Property{"地址", config.MYSQL_IP})

	err := dbc.Conn(config.MYSQL_NAME, config.MYSQL_IP, config.MYSQL_ACCOUNT, config.MYSQL_PWD, func() string {
		if config.MYSQL_COPY_PATH == "" {
			return config.GetDBBackupPath()
		} else {
			return config.MYSQL_COPY_PATH
		}
	}())
	if err != nil {
		log.Error("连接数据库失败 %v", err)
		return
	} else {
		log.Trace("连接数据库成功", nil)
		go dbc.Loop()
	}

	db_define := server_config.GetDBDefineFile("game_db.json")
	if !db_mgr.LoadConfig(db_define) {
		log.Error("载入数据库定义配置%v失败", db_define)
		return
	}

	if !db_mgr.Connect(config.MYSQL_IP, config.MYSQL_ACCOUNT, config.MYSQL_PWD, config.MYSQL_NAME) {
		log.Error("连接数据库失败")
		return
	}

	db_mgr.Run()

	tb_mgr := game_db.NewTablesManager(&db_mgr)
	battle_record_table = tb_mgr.Get_T_Battle_Record_Table()
	global_table = tb_mgr.Get_T_Global_Table()
	ban_player_table = tb_mgr.Get_T_Ban_Player_Table()

	if !signal_mgr.Init() {
		log.Error("signal_mgr init failed")
		return
	}

	// 配置加载
	if !global_config.Init("global.json") {
		log.Error("global_config_load failed !")
		return
	} else {
		log.Trace("global_config_load succeed !")
	}

	if !msg_handler_mgr.Init() {
		log.Error("msg_handler_mgr init failed !")
		return
	} else {
		log.Trace("msg_handler_mgr init succeed !")
	}

	if !player_mgr.Init() {
		log.Error("player_mgr init failed !")
		return
	} else {
		log.Trace("player_mgr init succeed !")
	}

	if !login_token_mgr.Init() {
		log.Error("启动login_token_mgr失败")
		return
	}

	if err := table_init(); err != nil {
		log.Error("%v", err.Error())
		return
	}

	// pprof
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	// 排行榜
	rank_list_mgr.Init()

	// 好友推荐
	friend_recommend_mgr.Init()

	// 月卡管理
	charge_month_card_manager.Init()

	// 远征战力匹配管理
	top_power_match_manager = NewTopPowerMatchManager(&TopPowerRankItem{}, 100000)

	if nil != dbc.Preload() {
		log.Error("dbc Preload Failed !!")
		return
	} else {
		log.Info("dbc Preload succeed !!")
	}

	if !login_conn_mgr.Init() {
		log.Error("login_conn_mgr init failed")
		return
	}

	// 初始化CenterServer
	center_conn.Init()

	// 初始化大厅
	if !hall_server.Init() {
		log.Error("hall_server init failed !")
		return
	} else {
		log.Trace("hall_server init succeed !")
	}

	if signal_mgr.IfClosing() {
		return
	}

	// 连接CenterServer
	log.Trace("连接中心服务器！！")
	go center_conn.Start()
	center_conn.WaitConnectFinished()

	after_center_match_conn()

	hall_server.Start(config.UseHttps)
}
