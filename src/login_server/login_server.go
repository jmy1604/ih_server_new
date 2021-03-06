package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"ih_server_new/libs/log"
	"ih_server_new/libs/timer"
	"ih_server_new/libs/utils"
	"ih_server_new/proto/gen_go/client_message"
	"ih_server_new/proto/gen_go/client_message_id"
	"ih_server_new/proto/gen_go/server_message"
	"ih_server_new/src/server_config"
	"ih_server_new/src/share_data"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)

type WaitCenterInfo struct {
	res_chan    chan *msg_server_message.C2LPlayerAccInfo
	create_time int32
}

type LoginServer struct {
	start_time         time.Time
	quit               bool
	shutdown_lock      *sync.Mutex
	shutdown_completed bool
	ticker             *timer.TickTimer
	initialized        bool

	login_http_listener net.Listener
	login_http_server   http.Server
	use_https           bool

	redis_conn *utils.RedisConn

	acc2c_wait      map[string]*WaitCenterInfo
	acc2c_wait_lock *sync.RWMutex
}

var server *LoginServer

func (this *LoginServer) Init() (ok bool) {
	this.start_time = time.Now()
	this.shutdown_lock = &sync.Mutex{}
	this.acc2c_wait = make(map[string]*WaitCenterInfo)
	this.acc2c_wait_lock = &sync.RWMutex{}
	this.redis_conn = &utils.RedisConn{}
	share_data.UidPlayerListInit()
	//account_mgr_init()

	this.initialized = true

	return true
}

func (this *LoginServer) Start(use_https bool) bool {
	if !this.redis_conn.Connect(config.RedisServerIP) {
		return false
	}

	go server_list.Run()

	if use_https {
		go this.StartHttps(server_config.GetConfPathFile("server.crt"), server_config.GetConfPathFile("server.key"))
	} else {
		go this.StartHttp()
	}

	this.use_https = use_https
	log.Event("服务器已启动", nil, log.Property{"IP", config.ListenClientIP})
	log.Trace("**************************************************")

	this.Run()

	return true
}

func (this *LoginServer) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}

		this.shutdown_completed = true
	}()

	this.ticker = timer.NewTickTimer(1000)
	this.ticker.Start()
	defer this.ticker.Stop()

	go this.redis_conn.Run(100)

	for {
		select {
		case d, ok := <-this.ticker.Chan:
			{
				if !ok {
					return
				}

				begin := time.Now()
				this.OnTick(d)
				time_cost := time.Now().Sub(begin).Seconds()
				if time_cost > 1 {
					log.Trace("耗时 %v", time_cost)
					if time_cost > 30 {
						log.Error("耗时 %v", time_cost)
					}
				}
			}
		}
	}
}

func (this *LoginServer) Shutdown() {
	if !this.initialized {
		return
	}

	this.shutdown_lock.Lock()
	defer this.shutdown_lock.Unlock()

	if this.quit {
		return
	}
	this.quit = true

	this.redis_conn.Close()

	log.Trace("关闭游戏主循环")

	begin := time.Now()

	if this.ticker != nil {
		this.ticker.Stop()
	}

	for {
		if this.shutdown_completed {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	this.login_http_listener.Close()
	center_conn.ShutDown()
	hall_agent_manager.net.Shutdown()

	db.Save()
	db.Close()

	log.Trace("关闭游戏主循环耗时 %v 秒", time.Now().Sub(begin).Seconds())
}

func (this *LoginServer) OnTick(t timer.TickTime) {
}

func (this *LoginServer) add_to_c_wait(acc string, c_wait *WaitCenterInfo) {
	this.acc2c_wait_lock.Lock()
	defer this.acc2c_wait_lock.Unlock()

	this.acc2c_wait[acc] = c_wait
}

func (this *LoginServer) remove_c_wait(acc string) {
	this.acc2c_wait_lock.Lock()
	defer this.acc2c_wait_lock.Unlock()

	delete(this.acc2c_wait, acc)
}

func (this *LoginServer) get_c_wait_by_acc(acc string) *WaitCenterInfo {
	this.acc2c_wait_lock.RLock()
	defer this.acc2c_wait_lock.RUnlock()

	return this.acc2c_wait[acc]
}

func (this *LoginServer) pop_c_wait_by_acc(acc string) *WaitCenterInfo {
	this.acc2c_wait_lock.Lock()
	defer this.acc2c_wait_lock.Unlock()

	cur_wait := this.acc2c_wait[acc]
	if nil != cur_wait {
		delete(this.acc2c_wait, acc)
		return cur_wait
	}

	return nil
}

//=================================================================================

type LoginHttpHandle struct{}

func (this *LoginServer) StartHttp() bool {
	var err error
	this.reg_http_mux()

	this.login_http_listener, err = net.Listen("tcp", config.ListenClientIP)
	if nil != err {
		log.Error("LoginServer StartHttp Failed %s", err.Error())
		return false
	}

	login_http_server := http.Server{
		Handler:     &LoginHttpHandle{},
		ReadTimeout: 6 * time.Second,
	}

	err = login_http_server.Serve(this.login_http_listener)
	if err != nil {
		log.Error("启动Login Http Server %s", err.Error())
		return false
	}

	return true
}

func (this *LoginServer) StartHttps(crt_file, key_file string) bool {
	this.reg_http_mux()

	this.login_http_server = http.Server{
		Addr:        config.ListenClientIP,
		Handler:     &LoginHttpHandle{},
		ReadTimeout: 6 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	err := this.login_http_server.ListenAndServeTLS(crt_file, key_file)
	if err != nil {
		log.Error("启动https server error[%v]", err.Error())
		return false
	}

	return true
}

var login_http_mux map[string]func(http.ResponseWriter, *http.Request)

func (this *LoginServer) reg_http_mux() {
	login_http_mux = make(map[string]func(http.ResponseWriter, *http.Request))
	login_http_mux["/client"] = client_http_handler
}

func (this *LoginHttpHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var act_str, url_str string
	url_str = r.URL.String()
	idx := strings.Index(url_str, "?")
	if -1 == idx {
		act_str = url_str
	} else {
		act_str = string([]byte(url_str)[:idx])
	}
	log.Info("ServeHTTP actstr(%s)", act_str)
	if h, ok := login_http_mux[act_str]; ok {
		h(w, r)
	}
	return
}

type JsonRequestData struct {
	MsgId   int32  // 消息ID
	MsgData []byte // 消息体
}

type JsonResponseData struct {
	Code    int32  // 错误码
	MsgId   int32  // 消息ID
	MsgData []byte // 消息体
}

func _check_register(account, password string) (err_code int32) {
	if b, err := regexp.MatchString(`^[\.a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, account); !b {
		if err != nil {
			log.Error("account[%v] not valid account, err %v", account, err.Error())
		} else {
			log.Error("account[%v] not match", account)
		}
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_IS_INVALID)
		return
	}

	if password == "" {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_PASSWORD_INVALID)
		return
	}

	err_code = 1
	return
}

func _generate_account_uuid(account string) string {
	uid, err := uuid.NewV1()
	if err != nil {
		log.Error("Account %v generate uuid error %v", account, err.Error())
		return ""
	}
	return uid.String()
}

func register_handler(account, password string, is_guest bool) (err_code int32, resp_data []byte) {
	if len(account) > 128 {
		log.Error("Account[%v] length %v too long", account, len(account))
		return -1, nil
	}

	if len(password) > 32 {
		log.Error("Account[%v] password[%v] length %v too long", account, password, len(password))
		return -1, nil
	}

	account_row := acc_mgr.Get(account)
	if account_row != nil {
		log.Error("Account[%v] already exists", account)
		return int32(msg_client_message.E_ERR_ACCOUNT_ALREADY_REGISTERED), nil
	}

	var o bool
	account_row, o = account_table.Select("account", account)
	if o {
		log.Error("Account[%v] already exists", account)
		return int32(msg_client_message.E_ERR_ACCOUNT_ALREADY_REGISTERED), nil
	}

	account_row = account_table.NewRow(account)

	if !is_guest {
		err_code = _check_register(account, password)
		if err_code < 0 {
			return
		}
	}

	uid := _generate_account_uuid(account)
	if uid == "" {
		err_code = -1
		return
	}

	account_row.Lock()
	account_row.Set_unique_id(uid)
	account_row.Set_password(password)
	account_row.Set_register_time(int32(time.Now().Unix()))
	if is_guest {
		account_row.Set_channel("guest")
	}
	account_table.Insert(account_row)
	account_row.Unlock()

	var response msg_client_message.S2CRegisterResponse = msg_client_message.S2CRegisterResponse{
		Account:  account,
		Password: password,
		IsGuest:  is_guest,
	}

	var err error
	resp_data, err = proto.Marshal(&response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("login_handler marshal response error: %v", err.Error())
		return
	}

	log.Debug("Account[%v] password[%v] registered", account, password)

	err_code = 1
	return
}

func bind_new_account_handler(server_id int32, account, password, new_account, new_password, new_channel string) (err_code int32, resp_data []byte) {
	if len(new_account) > 128 {
		log.Error("Account[%v] length %v too long", new_account, len(new_account))
		return -1, nil
	}

	if new_channel != "facebook" && len(new_password) > 32 {
		log.Error("Account[%v] password[%v] length %v too long", new_account, new_password, len(new_password))
		return -1, nil
	}

	if account == new_account {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NAME_MUST_DIFFRENT_TO_OLD)
		log.Error("Account %v can not bind same new account", account)
		return
	}

	row := acc_mgr.Get(account)
	if row == nil {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NOT_REGISTERED)
		log.Error("Account %v not registered, cant bind new account", account)
		return
	}

	ban_row := ban_mgr.Get(row.Get_unique_id())
	if ban_row != nil && ban_row.Get_start_time() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	if row.Get_password() != password {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_PASSWORD_INVALID)
		log.Error("Account %v password %v invalid, cant bind new account", account, password)
		return
	}

	if row.Get_channel() != "guest" {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NOT_GUEST)
		log.Error("Account %v not guest", account)
		return
	}

	if row.Get_bind_new_account() != "" {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_ALREADY_BIND)
		log.Error("Account %v already bind", account)
		return
	}

	if acc_mgr.Get(new_account) != nil {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NEW_BIND_ALREADY_EXISTS)
		log.Error("New Account %v to bind already exists", new_account)
		return
	}

	if new_channel != "" {
		if new_channel == "facebook" {
			err_code = _verify_facebook_login(new_account, new_password)
			if err_code < 0 {
				return
			}
		} else {
			err_code = -1
			log.Error("Account %v bind a unsupported channel %v account %v", account, new_channel, new_account)
			return
		}
	} else {
		err_code = _check_register(new_account, new_password)
		if err_code < 0 {
			return
		}
	}

	row.Set_bind_new_account(new_account)
	register_time := row.Get_register_time()
	uid := row.Get_unique_id()

	var last_server_id int32
	last_server_id = row.Get_last_server_id()

	new_row := account_table.NewRow(new_account)
	if new_channel == "" {
		new_row.Set_password(new_password)
	}
	new_row.Set_register_time(register_time)
	new_row.Set_unique_id(uid)
	new_row.Set_before_bind_account(account)
	new_row.Set_last_server_id(last_server_id)

	// sync to db
	account_table.UpdateWithFieldName(row, []string{"bind_new_account"})
	account_table.Insert(new_row)

	hall_agent := hall_agent_manager.GetAgentByID(server_id)
	if nil == hall_agent {
		err_code = int32(msg_client_message.E_ERR_PLAYER_SELECT_SERVER_NOT_FOUND)
		log.Error("login_http_handler get hall_agent failed")
		return
	}

	req := &msg_server_message.L2HBindNewAccountRequest{
		UniqueId:   uid,
		Account:    account,
		NewAccount: new_account,
	}
	hall_agent.Send(uint16(msg_server_message.MSGID_L2H_BIND_NEW_ACCOUNT_REQUEST), req)

	response := &msg_client_message.S2CGuestBindNewAccountResponse{
		Account:     account,
		NewAccount:  new_account,
		NewPassword: new_password,
		NewChannel:  new_channel,
	}

	var err error
	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("login_handler marshal response error: %v", err.Error())
		return
	}

	log.Debug("Account[%v] bind new account[%v]", account, new_account)
	err_code = 1
	return
}

func _verify_facebook_login(user_id, input_token string) int32 {

	var resp *http.Response
	var err error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	type _facebook_data struct {
		AppID     string `json:"app_id"`
		IsValid   bool   `json:"is_valid"`
		UserID    string `json:"user_id"`
		IssuedAt  int    `json:"issued_at"`
		ExpiresAt int    `json:"expires_at"`
	}

	type facebook_data struct {
		Data _facebook_data `json:"data"`
	}

	var verified bool
	for i := 0; i < len(config.Facebook); i++ {
		url_str := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%v&access_token=%v|%v", input_token, config.Facebook[i].FacebookAppID, config.Facebook[i].FacebookAppSecret)
		log.Debug("verify facebook url: %v", url_str)

		client := &http.Client{Transport: tr}
		resp, err = client.Get(url_str)
		if nil != err {
			log.Error("Facebook verify error %s", err.Error())
			continue
		}

		if resp.StatusCode != 200 {
			log.Error("Facebook verify response code %v", resp.StatusCode)
			continue
		}

		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if nil != err {
			log.Error("Read facebook verify result err(%s) !", err.Error())
			continue
		}

		log.Debug("facebook verify result data: %v", string(data))

		var fdata facebook_data
		err = json.Unmarshal(data, &fdata)
		if nil != err {
			log.Error("Facebook verify ummarshal err(%s)", err.Error())
			continue
		}

		if !fdata.Data.IsValid {
			log.Error("Facebook verify input_token[%v] failed", input_token)
			continue
		}

		if fdata.Data.UserID != user_id {
			log.Error("Facebook verify client user_id[%v] different to result user_id[%v]", user_id, fdata.Data.UserID)
			continue
		}

		verified = true
		break
	}

	if !verified {
		return -1
	}

	log.Debug("Facebook verify user_id[%v] and input_token[%v] success", user_id, input_token)

	return 1
}

func login_handler(account, password, channel, client_os string, is_verify bool) (err_code int32, resp_data []byte) {
	var is_new bool
	acc_row := acc_mgr.Get(account)
	if config.VerifyAccount {
		if channel == "" {
			if acc_row == nil {
				err_code = int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR)
				log.Error("Account %v not exist", account)
				return
			}
			if acc_row.Get_password() != password {
				err_code = int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR)
				log.Error("Account %v password %v invalid", account, password)
				return
			}
		} else if channel == "facebook" {
			err_code = _verify_facebook_login(account, password)
			if err_code < 0 {
				return
			}
			if acc_row == nil {
				acc_row = account_table.NewRow(account) //dbc.Accounts.AddRow(account)
				acc_row.Set_channel("facebook")
				is_new = true
			}
		} else if channel == "guest" {
			if acc_row == nil {
				acc_row = account_table.NewRow(account) //dbc.Accounts.AddRow(account)
				acc_row.Set_channel("guest")
				is_new = true
			} else {
				if acc_row.Get_password() != password {
					err_code = int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR)
					log.Error("Account %v password %v invalid", account, password)
					return
				}
			}
		} else {
			log.Error("Account %v use unsupported channel %v login", account, channel)
			return -1, nil
		}
	} else {
		if acc_row == nil {
			acc_row = account_table.NewRow(account) //dbc.Accounts.AddRow(account)
			is_new = true
		}
	}

	if acc_row.Get_unique_id() == "" {
		uid := _generate_account_uuid(account)
		if uid != "" {
			acc_row.Set_unique_id(uid)
		}
	}

	now_time := time.Now()
	last_time := acc_row.Get_last_get_player_list_time()
	if int32(now_time.Unix())-last_time >= 5*60 {
		share_data.LoadUidPlayerList(server.redis_conn, acc_row.Get_unique_id())
		acc_row.Set_last_get_player_list_time(int32(now_time.Unix()))
	}

	ban_row := ban_mgr.Get(acc_row.Get_unique_id())
	if ban_row != nil && ban_row.Get_start_time() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	// --------------------------------------------------------------------------------------------
	// 选择默认服
	var select_server_id int32
	select_server_id = acc_row.Get_last_server_id()
	if select_server_id <= 0 {
		server := server_list.RandomOneServer(client_os)
		if server == nil {
			err_code = int32(msg_client_message.E_ERR_INTERNAL)
			log.Error("Server List random null !!!")
			return
		}
		select_server_id = server.Id
		acc_row.Set_last_server_id(select_server_id)
	}

	var hall_ip, token string
	err_code, hall_ip, token = _select_server(acc_row.Get_unique_id(), account, select_server_id)
	if err_code < 0 {
		return
	}
	// --------------------------------------------------------------------------------------------

	acc_row.Set_token(token)

	// sync to db
	if is_new {
		account_table.Insert(acc_row)
	} else {
		account_table.UpdateWithFieldName(acc_row, []string{"last_server_id", "token"})
	}

	response := &msg_client_message.S2CLoginResponse{
		Acc:    account,
		Token:  token,
		GameIP: hall_ip,
	}

	if server_list.Servers == nil {
		response.Servers = make([]*msg_client_message.ServerInfo, 0)
	} else {
		servers := server_list.GetServers(client_os)
		l := len(servers)
		response.Servers = make([]*msg_client_message.ServerInfo, l)
		for i := 0; i < l; i++ {
			response.Servers[i] = &msg_client_message.ServerInfo{
				Id:   servers[i].Id,
				Name: servers[i].Name,
				IP:   servers[i].IP,
			}
		}
	}

	response.InfoList = share_data.GetUidPlayerList(acc_row.Get_unique_id())
	response.LastServerId = select_server_id
	if channel == "guest" {
		response.BoundAccount = acc_row.Get_bind_new_account()
	}

	var err error
	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("login_handler marshal response error: %v", err.Error())
		return
	}

	log.Debug("Account[%v] channel[%v] logined", account, channel)

	return
}

func _select_server(unique_id, account string, server_id int32) (err_code int32, hall_ip, access_token string) {
	sinfo := server_list.GetServerById(server_id)
	if sinfo == nil {
		err_code = int32(msg_client_message.E_ERR_PLAYER_SELECT_SERVER_NOT_FOUND)
		log.Error("select_server_handler player[%v] select server[%v] not found")
		return
	}

	hall_agent := hall_agent_manager.GetAgentByID(server_id)
	if nil == hall_agent {
		err_code = int32(msg_client_message.E_ERR_PLAYER_SELECT_SERVER_NOT_FOUND)
		log.Error("login_http_handler account %v get hall_agent failed by server_id %v", account, server_id)
		return
	}

	access_token = share_data.GenerateAccessToken(unique_id)
	hall_agent.Send(uint16(msg_server_message.MSGID_L2H_SYNC_ACCOUNT_TOKEN), &msg_server_message.L2HSyncAccountToken{
		UniqueId: unique_id,
		Account:  account,
		Token:    access_token,
	})

	hall_ip = sinfo.IP

	err_code = 1

	return
}

func select_server_handler(account, token string, server_id int32) (err_code int32, resp_data []byte) {
	row := acc_mgr.Get(account)
	if row == nil {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_NOT_REGISTERED)
		log.Error("select_server_handler: account[%v] not register", account)
		return
	}

	ban_row := ban_mgr.Get(row.Get_unique_id())
	if ban_row != nil && ban_row.Get_start_time() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	if /*token != acc.get_token()*/ token != row.Get_token() {
		err_code = int32(msg_client_message.E_ERR_PLAYER_TOKEN_ERROR)
		log.Error("select_server_handler player[%v] token[%v] invalid, need[%v]", account, token, row.Get_token())
		return
	}

	err_code, hall_ip, access_token := _select_server(row.Get_unique_id(), account, server_id)
	if err_code < 0 {
		return
	}

	if server.use_https {
		hall_ip = "https://" + hall_ip
	} else {
		hall_ip = "http://" + hall_ip
	}

	response := &msg_client_message.S2CSelectServerResponse{
		Acc:   account,
		Token: access_token,
		IP:    hall_ip,
	}

	var err error
	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("select_server_handler marshal response error: %v", err.Error())
		return
	}

	row.Set_last_server_id(server_id)

	log.Trace("Account %v selected server %v", account, server_id)

	return
}

func set_password_handler(account, password, new_password string) (err_code int32, resp_data []byte) {
	row := acc_mgr.Get(account)
	if row == nil {
		err_code = int32(msg_client_message.E_ERR_PLAYER_NOT_EXIST)
		log.Error("set_password_handler account[%v] not found", account)
		return
	}

	ban_row := ban_mgr.Get(row.Get_unique_id())
	if ban_row != nil && ban_row.Get_start_time() > 0 {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_BE_BANNED)
		log.Error("Account %v has been banned, cant login", account)
		return
	}

	if row.Get_password() != password {
		err_code = int32(msg_client_message.E_ERR_ACCOUNT_PASSWORD_INVALID)
		log.Error("set_password_handler account[%v] password is invalid", account)
		return
	}

	row.Set_password(new_password)

	response := &msg_client_message.S2CSetLoginPasswordResponse{
		Account:     account,
		Password:    password,
		NewPassword: new_password,
	}

	var err error
	resp_data, err = proto.Marshal(response)
	if err != nil {
		err_code = int32(msg_client_message.E_ERR_INTERNAL)
		log.Error("set_password_handler marshal response error: %v", err.Error())
		return
	}

	return
}

func response_error(err_code int32, w http.ResponseWriter) {
	err_response := JsonResponseData{
		Code: err_code,
	}
	data, err := json.Marshal(err_response)
	if nil != err {
		log.Error("login_http_handler json mashal error")
		return
	}
	w.Write(data)
}

func register_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			return
		}
	}()

	if !config.VerifyAccount {
		response_error(-1, w)
		log.Error("no need verify account and no need register")
		return
	}

	account := r.URL.Query().Get("account")

	password := r.URL.Query().Get("password")
	is_guest := r.URL.Query().Get("is_guest")
	ig, err := strconv.Atoi(is_guest)
	if err != nil {
		response_error(-1, w)
		log.Error("is_guest %v set invalid", is_guest)
		return
	}

	if ig == 0 && password == "" {
		response_error(-1, w)
		log.Error("password can not set to empty")
		return
	}

	err_code, data := register_handler(account, password, func() bool {
		if ig > 0 {
			return true
		}
		return false
	}())

	if err_code < 0 {
		response_error(err_code, w)
		log.Error("login_http_handler err_code[%v]", err_code)
		return
	}

	if data == nil {
		response_error(-1, w)
		log.Error("cant get response data failed")
		return
	}

	http_res := &JsonResponseData{Code: 0, MsgId: int32(msg_client_message_id.MSGID_S2C_REGISTER_RESPONSE), MsgData: data}
	data, err = json.Marshal(http_res)
	if nil != err {
		response_error(-1, w)
		log.Error("login_http_handler json mashal error")
		return
	}
	w.Write(data)

	log.Debug("New account %v registered, is_guest %v", account, is_guest)
}

func bind_new_account_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			return
		}
	}()

	if !config.VerifyAccount {
		response_error(-1, w)
		log.Error("no need bind new account")
		return
	}

	server_id_str := r.URL.Query().Get("server_id")
	server_id, err := strconv.Atoi(server_id_str)
	if err != nil {
		log.Error("server_id convert err %v", err.Error())
		return
	}
	account := r.URL.Query().Get("account")
	password := r.URL.Query().Get("password")

	new_account := r.URL.Query().Get("new_account")
	new_password := r.URL.Query().Get("new_password")
	new_channel := r.URL.Query().Get("new_channel")

	err_code, data := bind_new_account_handler(int32(server_id), account, password, new_account, new_password, new_channel)
	if err_code < 0 {
		response_error(err_code, w)
		log.Error("login_http_handler err_code[%v]", err_code)
		return
	}

	if data == nil {
		response_error(-1, w)
		log.Error("cant get response data failed")
		return
	}

	http_res := &JsonResponseData{Code: 0, MsgId: int32(msg_client_message_id.MSGID_S2C_GUEST_BIND_NEW_ACCOUNT_RESPONSE), MsgData: data}
	data, err = json.Marshal(http_res)
	if nil != err {
		response_error(-1, w)
		log.Error("login_http_handler json mashal error")
		return
	}
	w.Write(data)

	log.Debug("Account %v bind new account %v", account, new_account)
}

func login_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			return
		}
	}()

	// account
	account := r.URL.Query().Get("account")
	if "" == account {
		response_error(int32(msg_client_message.E_ERR_PLAYER_ACC_OR_PASSWORD_ERROR), w)
		log.Error("login_http_handler get msg_id failed")
		return
	}

	// password
	password := r.URL.Query().Get("password")

	// channel
	channel := r.URL.Query().Get("channel")

	// client os
	client_os := r.URL.Query().Get("client_os")

	log.Debug("account: %v, password: %v, channel: %v", account, password, channel)

	var err_code int32
	var data []byte
	err_code, data = login_handler(account, password, channel, client_os, false)

	if err_code < 0 {
		response_error(err_code, w)
		log.Error("login_http_handler err_code[%v]", err_code)
		return
	}

	if data == nil {
		response_error(-1, w)
		log.Error("cant get response data failed")
		return
	}

	http_res := &JsonResponseData{Code: 0, MsgId: int32(msg_client_message_id.MSGID_S2C_LOGIN_RESPONSE), MsgData: data}
	var err error
	data, err = json.Marshal(http_res)
	if nil != err {
		response_error(-1, w)
		log.Error("login_http_handler json mashal error")
		return
	}
	w.Write(data)
	log.Debug("Account %v logined, channel %v", account, channel)
}

func select_server_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			return
		}
	}()

	account := r.URL.Query().Get("account")
	if "" == account {
		response_error(-1, w)
		log.Error("login_http_handler get account is empty")
		return
	}

	token := r.URL.Query().Get("token")
	if "" == token {
		response_error(-1, w)
		log.Error("login_http_handler get token is empty")
		return
	}

	server_id_str := r.URL.Query().Get("server_id")
	if "" == server_id_str {
		response_error(-1, w)
		log.Error("login_http_handler get server_id is empty")
		return
	}

	server_id, err := strconv.Atoi(server_id_str)
	if err != nil {
		response_error(-1, w)
		log.Error("login_http_handler transfer server_id[%v] error[%v]", server_id_str, err.Error())
		return
	}
	log.Debug("account: %v, token: %v, server_id: %v", account, token, server_id)

	var err_code int32
	var data []byte
	err_code, data = select_server_handler(account, token, int32(server_id))

	if err_code < 0 {
		response_error(err_code, w)
		log.Error("login_http_handler err_code[%v]", err_code)
		return
	}

	if data == nil {
		response_error(-1, w)
		log.Error("cant get response data")
		return
	}

	http_res := &JsonResponseData{Code: 0, MsgId: int32(msg_client_message_id.MSGID_S2C_SELECT_SERVER_RESPONSE), MsgData: data}
	data, err = json.Marshal(http_res)
	if nil != err {
		response_error(-1, w)
		log.Error("login_http_handler json mashal error")
		return
	}
	w.Write(data)
}

func set_password_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			return
		}
	}()

	account := r.URL.Query().Get("account")
	if "" == account {
		response_error(-1, w)
		log.Error("set_password_http_handler get account is empty")
		return
	}

	password := r.URL.Query().Get("password")

	new_password := r.URL.Query().Get("new_password")
	if "" == new_password {
		response_error(-1, w)
		log.Error("set_password_http_handler get new password is empty")
		return
	}

	if password == new_password {
		response_error(-1, w)
		log.Error("set_password_http_handler set password must different to old password")
		return
	}

	log.Debug("account: %v, password: %v, new_password: %v", account, password, new_password)

	var err_code int32
	var data []byte
	err_code, data = set_password_handler(account, password, new_password)

	if err_code < 0 {
		response_error(err_code, w)
		log.Error("set_password_http_handler err_code[%v]", err_code)
		return
	}

	if data == nil {
		response_error(-1, w)
		log.Error("cant get response data")
		return
	}

	http_res := &JsonResponseData{Code: 0, MsgId: int32(msg_client_message_id.MSGID_S2C_SET_LOGIN_PASSWORD_RESPONSE), MsgData: data}
	var err error
	data, err = json.Marshal(http_res)
	if nil != err {
		response_error(-1, w)
		log.Error("set_password_http_handler json mashal error")
		return
	}
	w.Write(data)
}

func _send_error(w http.ResponseWriter, msg_id, ret_code int32) {
	m := &msg_client_message.S2C_ONE_MSG{ErrorCode: ret_code}
	res2cli := &msg_client_message.S2C_MSG_DATA{MsgList: []*msg_client_message.S2C_ONE_MSG{m}}
	final_data, err := proto.Marshal(res2cli)
	if nil != err {
		log.Error("client_msg_handler marshal 1 client msg failed err(%s)", err.Error())
		return
	}

	data := final_data
	data = append(data, byte(0))

	iret, err := w.Write(data)
	if nil != err {
		log.Error("client_msg_handler write data 1 failed err[%s] ret %d", err.Error(), iret)
		return
	}
}

func client_http_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
			debug.PrintStack()
		}
	}()

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_http_handler ReadAll err[%s]", err.Error())
		return
	}

	var msg msg_client_message.C2S_ONE_MSG
	err = proto.Unmarshal(data, &msg)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_http_handler proto Unmarshal err[%s]", err.Error())
		return
	}

	var err_code, msg_id int32
	if msg.MsgCode == int32(msg_client_message_id.MSGID_C2S_LOGIN_REQUEST) {
		var login_msg msg_client_message.C2SLoginRequest
		err = proto.Unmarshal(msg.GetData(), &login_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SLoginRequest unmarshal err %v", err.Error())
			return
		}
		if login_msg.GetAcc() == "" {
			_send_error(w, 0, -1)
			log.Error("Acc is empty")
			return
		}
		msg_id = int32(msg_client_message_id.MSGID_S2C_LOGIN_RESPONSE)
		err_code, data = login_handler(login_msg.GetAcc(), login_msg.GetPassword(), login_msg.GetChannel(), login_msg.GetClientOS(), login_msg.GetIsAppleVerifyUse())
	} else if msg.MsgCode == int32(msg_client_message_id.MSGID_C2S_SELECT_SERVER_REQUEST) {
		var select_msg msg_client_message.C2SSelectServerRequest
		err = proto.Unmarshal(msg.GetData(), &select_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SSelectServerRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message_id.MSGID_S2C_SELECT_SERVER_RESPONSE)
		err_code, data = select_server_handler(select_msg.GetAcc(), select_msg.GetToken(), select_msg.GetServerId())
	} else if msg.MsgCode == int32(msg_client_message_id.MSGID_C2S_REGISTER_REQUEST) {
		var register_msg msg_client_message.C2SRegisterRequest
		err = proto.Unmarshal(msg.GetData(), &register_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SRegisterRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message_id.MSGID_S2C_REGISTER_RESPONSE)
		err_code, data = register_handler(register_msg.GetAccount(), register_msg.GetPassword(), register_msg.GetIsGuest())
	} else if msg.MsgCode == int32(msg_client_message_id.MSGID_C2S_SET_LOGIN_PASSWORD_REQUEST) {
		var pass_msg msg_client_message.C2SSetLoginPasswordRequest
		err = proto.Unmarshal(msg.GetData(), &pass_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SSetLoginPasswordRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message_id.MSGID_S2C_SET_LOGIN_PASSWORD_RESPONSE)
		err_code, data = set_password_handler(pass_msg.GetAccount(), pass_msg.GetPassword(), pass_msg.GetNewPassword())
	} else if msg.MsgCode == int32(msg_client_message_id.MSGID_C2S_GUEST_BIND_NEW_ACCOUNT_REQUEST) {
		var bind_msg msg_client_message.C2SGuestBindNewAccountRequest
		err = proto.Unmarshal(msg.GetData(), &bind_msg)
		if err != nil {
			_send_error(w, 0, -1)
			log.Error("Msg C2SGuestBindNewAccountRequest unmarshal err %v", err.Error())
			return
		}
		msg_id = int32(msg_client_message_id.MSGID_S2C_GUEST_BIND_NEW_ACCOUNT_RESPONSE)
		err_code, data = bind_new_account_handler(bind_msg.GetServerId(), bind_msg.GetAccount(), bind_msg.GetPassword(), bind_msg.GetNewAccount(), bind_msg.GetNewPassword(), bind_msg.GetNewChannel())
	} else {
		if msg.MsgCode > 0 {
			_send_error(w, msg.MsgCode, int32(msg_client_message.E_ERR_PLAYER_MSG_ID_NOT_FOUND))
			log.Error("Unsupported msg %v", msg.MsgCode)
		} else {
			_send_error(w, msg.MsgCode, int32(msg_client_message.E_ERR_PLAYER_MSG_ID_INVALID))
			log.Error("Invalid msg %v", msg.MsgCode)
		}
		return
	}

	var resp_msg msg_client_message.S2C_ONE_MSG
	resp_msg.MsgCode = msg_id
	resp_msg.ErrorCode = err_code
	resp_msg.Data = data
	data, err = proto.Marshal(&resp_msg)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_msg_handler marshal 2 client msg failed err(%s)", err.Error())
		return
	}

	iret, err := w.Write(data)
	if nil != err {
		_send_error(w, 0, -1)
		log.Error("client_msg_handler write data 2 failed err[%s] ret %d", err.Error(), iret)
		return
	}
}
