package login_db

import (
	"sync"
	"github.com/huoshan017/mysql-go/base"
	"github.com/huoshan017/mysql-go/manager"
)

type t_account_field_pair_func func (t *T_Account) *mysql_base.FieldValuePair

var t_account_fields_map = map[string]t_account_field_pair_func{
	"account": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_account()
	},
	"unique_id": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_unique_id()
	},
	"password": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_password()
	},
	"register_time": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_register_time()
	},
	"channel": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_channel()
	},
	"token": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_token()
	},
	"last_server_id": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_last_server_id()
	},
	"bind_new_account": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_bind_new_account()
	},
	"before_bind_account": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_before_bind_account()
	},
	"last_get_player_list_time": func (t *T_Account) *mysql_base.FieldValuePair{
		return t.GetFVP_last_get_player_list_time()
	},
}

type T_Account struct {
	account string
	unique_id string
	password string
	register_time int32
	channel string
	token string
	last_server_id int32
	bind_new_account string
	before_bind_account string
	last_get_player_list_time int32
	locker sync.RWMutex
}

func Create_T_Account() *T_Account {
	return &T_Account{
	}
}

func (this *T_Account) Get_account() string {
	return this.account
}

func (this *T_Account) Set_account(v string) {
	this.account = v
}

func (this *T_Account) GetWithLock_account() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.account
}

func (this *T_Account) SetWithLock_account(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.account = v
}

func (this *T_Account) GetFVP_account() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "account", Value: this.Get_account() }
}

func (this *T_Account) Get_unique_id() string {
	return this.unique_id
}

func (this *T_Account) Set_unique_id(v string) {
	this.unique_id = v
}

func (this *T_Account) GetWithLock_unique_id() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.unique_id
}

func (this *T_Account) SetWithLock_unique_id(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.unique_id = v
}

func (this *T_Account) GetFVP_unique_id() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "unique_id", Value: this.Get_unique_id() }
}

func (this *T_Account) Get_password() string {
	return this.password
}

func (this *T_Account) Set_password(v string) {
	this.password = v
}

func (this *T_Account) GetWithLock_password() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.password
}

func (this *T_Account) SetWithLock_password(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.password = v
}

func (this *T_Account) GetFVP_password() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "password", Value: this.Get_password() }
}

func (this *T_Account) Get_register_time() int32 {
	return this.register_time
}

func (this *T_Account) Set_register_time(v int32) {
	this.register_time = v
}

func (this *T_Account) GetWithLock_register_time() int32 {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.register_time
}

func (this *T_Account) SetWithLock_register_time(v int32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.register_time = v
}

func (this *T_Account) GetFVP_register_time() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "register_time", Value: this.Get_register_time() }
}

func (this *T_Account) Get_channel() string {
	return this.channel
}

func (this *T_Account) Set_channel(v string) {
	this.channel = v
}

func (this *T_Account) GetWithLock_channel() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.channel
}

func (this *T_Account) SetWithLock_channel(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.channel = v
}

func (this *T_Account) GetFVP_channel() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "channel", Value: this.Get_channel() }
}

func (this *T_Account) Get_token() string {
	return this.token
}

func (this *T_Account) Set_token(v string) {
	this.token = v
}

func (this *T_Account) GetWithLock_token() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.token
}

func (this *T_Account) SetWithLock_token(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.token = v
}

func (this *T_Account) GetFVP_token() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "token", Value: this.Get_token() }
}

func (this *T_Account) Get_last_server_id() int32 {
	return this.last_server_id
}

func (this *T_Account) Set_last_server_id(v int32) {
	this.last_server_id = v
}

func (this *T_Account) GetWithLock_last_server_id() int32 {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.last_server_id
}

func (this *T_Account) SetWithLock_last_server_id(v int32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.last_server_id = v
}

func (this *T_Account) GetFVP_last_server_id() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "last_server_id", Value: this.Get_last_server_id() }
}

func (this *T_Account) Get_bind_new_account() string {
	return this.bind_new_account
}

func (this *T_Account) Set_bind_new_account(v string) {
	this.bind_new_account = v
}

func (this *T_Account) GetWithLock_bind_new_account() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.bind_new_account
}

func (this *T_Account) SetWithLock_bind_new_account(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.bind_new_account = v
}

func (this *T_Account) GetFVP_bind_new_account() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "bind_new_account", Value: this.Get_bind_new_account() }
}

func (this *T_Account) Get_before_bind_account() string {
	return this.before_bind_account
}

func (this *T_Account) Set_before_bind_account(v string) {
	this.before_bind_account = v
}

func (this *T_Account) GetWithLock_before_bind_account() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.before_bind_account
}

func (this *T_Account) SetWithLock_before_bind_account(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.before_bind_account = v
}

func (this *T_Account) GetFVP_before_bind_account() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "before_bind_account", Value: this.Get_before_bind_account() }
}

func (this *T_Account) Get_last_get_player_list_time() int32 {
	return this.last_get_player_list_time
}

func (this *T_Account) Set_last_get_player_list_time(v int32) {
	this.last_get_player_list_time = v
}

func (this *T_Account) GetWithLock_last_get_player_list_time() int32 {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.last_get_player_list_time
}

func (this *T_Account) SetWithLock_last_get_player_list_time(v int32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.last_get_player_list_time = v
}

func (this *T_Account) GetFVP_last_get_player_list_time() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "last_get_player_list_time", Value: this.Get_last_get_player_list_time() }
}

func (this *T_Account) GetFVPList(fields_name []string) []*mysql_base.FieldValuePair {
	var field_list []*mysql_base.FieldValuePair
	for _, field_name := range fields_name {
		fun := t_account_fields_map[field_name]
		if fun == nil {
			continue
		}
		value_pair := fun(this)
		if value_pair != nil {
			field_list = append(field_list, value_pair)
		}
	}
	return field_list
}

func (this *T_Account) _format_field_list() []*mysql_base.FieldValuePair {
	var field_list []*mysql_base.FieldValuePair
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "account", Value: this.Get_account() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "unique_id", Value: this.Get_unique_id() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "password", Value: this.Get_password() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "register_time", Value: this.Get_register_time() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "channel", Value: this.Get_channel() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "token", Value: this.Get_token() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "last_server_id", Value: this.Get_last_server_id() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "bind_new_account", Value: this.Get_bind_new_account() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "before_bind_account", Value: this.Get_before_bind_account() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "last_get_player_list_time", Value: this.Get_last_get_player_list_time() })
	return field_list
}

func (this *T_Account) Lock() {
	this.locker.Lock()
}

func (this *T_Account) Unlock() {
	this.locker.Unlock()
}

func (this *T_Account) RLock() {
	this.locker.RLock()
}

func (this *T_Account) RUnlock() {
	this.locker.RUnlock()
}

type T_Account_AtomicExecFunc func(*T_Account)

func (this *T_Account) AtomicExecute(exec_func T_Account_AtomicExecFunc) {
	this.locker.Lock()
	defer this.locker.Unlock()
	exec_func(this)
}

func (this *T_Account) AtomicExecuteReadOnly(exec_func T_Account_AtomicExecFunc) {
	this.locker.RLock()
	defer this.locker.RUnlock()
	exec_func(this)
}

type T_Account_Table struct {
	db *mysql_manager.DB
}

func (this *T_Account_Table) Init(db *mysql_manager.DB) {
	this.db = db
}

func (this *T_Account_Table) Select(key string, value interface{}) (*T_Account, bool) {
	var field_list = []string{"account", "unique_id", "password", "register_time", "channel", "token", "last_server_id", "bind_new_account", "before_bind_account", "last_get_player_list_time"}
	var t = Create_T_Account()
	var dest_list = []interface{}{&t.account, &t.unique_id, &t.password, &t.register_time, &t.channel, &t.token, &t.last_server_id, &t.bind_new_account, &t.before_bind_account, &t.last_get_player_list_time}
	if !this.db.Select("t_account", key, value, field_list, dest_list) {
		return nil, false
	}
	return t, true
}

func (this *T_Account_Table) SelectMulti(key string, value interface{}) ([]*T_Account, bool) {
	var field_list = []string{"account", "unique_id", "password", "register_time", "channel", "token", "last_server_id", "bind_new_account", "before_bind_account", "last_get_player_list_time"}
	var result_list mysql_base.QueryResultList
	if !this.db.SelectRecords("t_account", key, value, field_list, &result_list) {
		return nil, false
	}
	var r []*T_Account
	for {
		var t = Create_T_Account()
		var dest_list = []interface{}{&t.account, &t.unique_id, &t.password, &t.register_time, &t.channel, &t.token, &t.last_server_id, &t.bind_new_account, &t.before_bind_account, &t.last_get_player_list_time}
		if !result_list.Get(dest_list...) {
			break
		}
		r = append(r, t)
	}
	return r, true
}

func (this *T_Account_Table) SelectPrimaryField() ([]string) {
	var result_list mysql_base.QueryResultList
	if !this.db.SelectFieldNoKey("t_account", "account", &result_list) {
		return nil
	}
	var value_list []string
	for {
		var d string
		if !result_list.Get(&d) {
			break
		}
		value_list = append(value_list, d)
	}
	return value_list
}

func (this *T_Account_Table) Insert(t *T_Account) {
	var field_list = t._format_field_list()
	if field_list != nil {
		this.db.Insert("t_account", field_list)
	}
}

func (this *T_Account_Table) InsertIgnore(t *T_Account) {
	var field_list = t._format_field_list()
	if field_list != nil {
		this.db.InsertIgnore("t_account", field_list)
	}
}

func (this *T_Account_Table) Delete(account string) {
	this.db.Delete("t_account", "account", account)
}

func (this *T_Account_Table) NewRow(account string) *T_Account {
	return &T_Account{ account: account, }
}

func (this *T_Account_Table) UpdateAll(t *T_Account) {
	var field_list = t._format_field_list()
	if field_list != nil {
		this.db.Update("t_account", "account", t.Get_account(), field_list)
	}
}

func (this *T_Account_Table) UpdateWithFVPList(account string, field_list []*mysql_base.FieldValuePair) {
	this.db.Update("t_account", "account", account, field_list)
}

func (this *T_Account_Table) UpdateWithFieldName(t *T_Account, fields_name []string) {
	var field_list = t.GetFVPList(fields_name)
	if field_list != nil {
		this.UpdateWithFVPList(t.Get_account(), field_list)
	}
}

func (this *T_Account_Table) TransactionInsert(transaction *mysql_manager.Transaction, t *T_Account) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Insert("t_account", field_list)
	}
}

func (this *T_Account_Table) TransactionDelete(transaction *mysql_manager.Transaction, account string) {
	transaction.Delete("t_account", "account", account)
}

func (this *T_Account_Table) TransactionUpdateAll(transaction *mysql_manager.Transaction, t*T_Account) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Update("t_account", "account", t.Get_account(), field_list)
	}
}

func (this *T_Account_Table) TransactionUpdateWithFVPList(transaction *mysql_manager.Transaction, account string, field_list []*mysql_base.FieldValuePair) {
	transaction.Update("t_account", "account", account, field_list)
}

func (this *T_Account_Table) TransactionUpdateWithFieldName(transaction *mysql_manager.Transaction, t *T_Account, fields_name []string) {
	field_list := t.GetFVPList(fields_name)
	if field_list != nil {
		transaction.Update("t_account", "account", t.Get_account(), field_list)
	}
}

func TransactionInsert_T_Account(transaction *mysql_manager.Transaction, t *T_Account) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Insert("t_account", field_list)
	}
}

func TransactionDelete_T_Account(transaction *mysql_manager.Transaction, account string) {
	transaction.Delete("t_account", "account", account)
}

func TransactionUpdateAll_T_Account(transaction *mysql_manager.Transaction, t *T_Account) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Update("t_account", "account", t.Get_account(), field_list)
	}
}

func TransactionUpdate_T_Account(transaction *mysql_manager.Transaction, t *T_Account, fields_name []string) {
	field_list := t.GetFVPList(fields_name)
	if field_list != nil {
		transaction.Update("t_account", "account", t.Get_account(), field_list)
	}
}
