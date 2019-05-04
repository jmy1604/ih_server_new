package login_db

import (
	"sync"
	"github.com/huoshan017/mysql-go/base"
	"github.com/huoshan017/mysql-go/manager"
)

type t_ban_player_field_pair_func func (t *T_Ban_Player) *mysql_base.FieldValuePair

var t_ban_player_fields_map = map[string]t_ban_player_field_pair_func{
	"unique_id": func (t *T_Ban_Player) *mysql_base.FieldValuePair{
		return t.GetFVP_unique_id()
	},
	"account": func (t *T_Ban_Player) *mysql_base.FieldValuePair{
		return t.GetFVP_account()
	},
	"player_id": func (t *T_Ban_Player) *mysql_base.FieldValuePair{
		return t.GetFVP_player_id()
	},
	"start_time": func (t *T_Ban_Player) *mysql_base.FieldValuePair{
		return t.GetFVP_start_time()
	},
	"start_time_str": func (t *T_Ban_Player) *mysql_base.FieldValuePair{
		return t.GetFVP_start_time_str()
	},
	"duration": func (t *T_Ban_Player) *mysql_base.FieldValuePair{
		return t.GetFVP_duration()
	},
}

type T_Ban_Player struct {
	unique_id string
	account string
	player_id int32
	start_time int32
	start_time_str string
	duration int32
	locker sync.RWMutex
}

func Create_T_Ban_Player() *T_Ban_Player {
	return &T_Ban_Player{
	}
}

func (this *T_Ban_Player) Get_unique_id() string {
	return this.unique_id
}

func (this *T_Ban_Player) Set_unique_id(v string) {
	this.unique_id = v
}

func (this *T_Ban_Player) GetWithLock_unique_id() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.unique_id
}

func (this *T_Ban_Player) SetWithLock_unique_id(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.unique_id = v
}

func (this *T_Ban_Player) GetFVP_unique_id() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "unique_id", Value: this.Get_unique_id() }
}

func (this *T_Ban_Player) Get_account() string {
	return this.account
}

func (this *T_Ban_Player) Set_account(v string) {
	this.account = v
}

func (this *T_Ban_Player) GetWithLock_account() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.account
}

func (this *T_Ban_Player) SetWithLock_account(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.account = v
}

func (this *T_Ban_Player) GetFVP_account() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "account", Value: this.Get_account() }
}

func (this *T_Ban_Player) Get_player_id() int32 {
	return this.player_id
}

func (this *T_Ban_Player) Set_player_id(v int32) {
	this.player_id = v
}

func (this *T_Ban_Player) GetWithLock_player_id() int32 {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.player_id
}

func (this *T_Ban_Player) SetWithLock_player_id(v int32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.player_id = v
}

func (this *T_Ban_Player) GetFVP_player_id() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "player_id", Value: this.Get_player_id() }
}

func (this *T_Ban_Player) Get_start_time() int32 {
	return this.start_time
}

func (this *T_Ban_Player) Set_start_time(v int32) {
	this.start_time = v
}

func (this *T_Ban_Player) GetWithLock_start_time() int32 {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.start_time
}

func (this *T_Ban_Player) SetWithLock_start_time(v int32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.start_time = v
}

func (this *T_Ban_Player) GetFVP_start_time() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "start_time", Value: this.Get_start_time() }
}

func (this *T_Ban_Player) Get_start_time_str() string {
	return this.start_time_str
}

func (this *T_Ban_Player) Set_start_time_str(v string) {
	this.start_time_str = v
}

func (this *T_Ban_Player) GetWithLock_start_time_str() string {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.start_time_str
}

func (this *T_Ban_Player) SetWithLock_start_time_str(v string) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.start_time_str = v
}

func (this *T_Ban_Player) GetFVP_start_time_str() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "start_time_str", Value: this.Get_start_time_str() }
}

func (this *T_Ban_Player) Get_duration() int32 {
	return this.duration
}

func (this *T_Ban_Player) Set_duration(v int32) {
	this.duration = v
}

func (this *T_Ban_Player) GetWithLock_duration() int32 {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.duration
}

func (this *T_Ban_Player) SetWithLock_duration(v int32) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.duration = v
}

func (this *T_Ban_Player) GetFVP_duration() *mysql_base.FieldValuePair {
	return &mysql_base.FieldValuePair{ Name: "duration", Value: this.Get_duration() }
}

func (this *T_Ban_Player) GetFVPList(fields_name []string) []*mysql_base.FieldValuePair {
	var field_list []*mysql_base.FieldValuePair
	for _, field_name := range fields_name {
		fun := t_ban_player_fields_map[field_name]
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

func (this *T_Ban_Player) _format_field_list() []*mysql_base.FieldValuePair {
	var field_list []*mysql_base.FieldValuePair
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "unique_id", Value: this.Get_unique_id() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "account", Value: this.Get_account() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "player_id", Value: this.Get_player_id() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "start_time", Value: this.Get_start_time() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "start_time_str", Value: this.Get_start_time_str() })
	field_list = append(field_list, &mysql_base.FieldValuePair{ Name: "duration", Value: this.Get_duration() })
	return field_list
}

func (this *T_Ban_Player) Lock() {
	this.locker.Lock()
}

func (this *T_Ban_Player) Unlock() {
	this.locker.Unlock()
}

func (this *T_Ban_Player) RLock() {
	this.locker.RLock()
}

func (this *T_Ban_Player) RUnlock() {
	this.locker.RUnlock()
}

type T_Ban_Player_AtomicExecFunc func(*T_Ban_Player)

func (this *T_Ban_Player) AtomicExecute(exec_func T_Ban_Player_AtomicExecFunc) {
	this.locker.Lock()
	defer this.locker.Unlock()
	exec_func(this)
}

func (this *T_Ban_Player) AtomicExecuteReadOnly(exec_func T_Ban_Player_AtomicExecFunc) {
	this.locker.RLock()
	defer this.locker.RUnlock()
	exec_func(this)
}

type T_Ban_Player_Table struct {
	db *mysql_manager.DB
}

func (this *T_Ban_Player_Table) Init(db *mysql_manager.DB) {
	this.db = db
}

func (this *T_Ban_Player_Table) Select(key string, value interface{}) (*T_Ban_Player, bool) {
	var field_list = []string{"unique_id", "account", "player_id", "start_time", "start_time_str", "duration"}
	var t = Create_T_Ban_Player()
	var dest_list = []interface{}{&t.unique_id, &t.account, &t.player_id, &t.start_time, &t.start_time_str, &t.duration}
	if !this.db.Select("t_ban_player", key, value, field_list, dest_list) {
		return nil, false
	}
	return t, true
}

func (this *T_Ban_Player_Table) SelectMulti(key string, value interface{}) ([]*T_Ban_Player, bool) {
	var field_list = []string{"unique_id", "account", "player_id", "start_time", "start_time_str", "duration"}
	var result_list mysql_base.QueryResultList
	if !this.db.SelectRecords("t_ban_player", key, value, field_list, &result_list) {
		return nil, false
	}
	var r []*T_Ban_Player
	for {
		var t = Create_T_Ban_Player()
		var dest_list = []interface{}{&t.unique_id, &t.account, &t.player_id, &t.start_time, &t.start_time_str, &t.duration}
		if !result_list.Get(dest_list...) {
			break
		}
		r = append(r, t)
	}
	return r, true
}

func (this *T_Ban_Player_Table) SelectPrimaryField() ([]string) {
	var result_list mysql_base.QueryResultList
	if !this.db.SelectFieldNoKey("t_ban_player", "unique_id", &result_list) {
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

func (this *T_Ban_Player_Table) Insert(t *T_Ban_Player) {
	var field_list = t._format_field_list()
	if field_list != nil {
		this.db.Insert("t_ban_player", field_list)
	}
}

func (this *T_Ban_Player_Table) InsertIgnore(t *T_Ban_Player) {
	var field_list = t._format_field_list()
	if field_list != nil {
		this.db.InsertIgnore("t_ban_player", field_list)
	}
}

func (this *T_Ban_Player_Table) Delete(unique_id string) {
	this.db.Delete("t_ban_player", "unique_id", unique_id)
}

func (this *T_Ban_Player_Table) NewRow(unique_id string) *T_Ban_Player {
	return &T_Ban_Player{ unique_id: unique_id, }
}

func (this *T_Ban_Player_Table) UpdateAll(t *T_Ban_Player) {
	var field_list = t._format_field_list()
	if field_list != nil {
		this.db.Update("t_ban_player", "unique_id", t.Get_unique_id(), field_list)
	}
}

func (this *T_Ban_Player_Table) UpdateWithFVPList(unique_id string, field_list []*mysql_base.FieldValuePair) {
	this.db.Update("t_ban_player", "unique_id", unique_id, field_list)
}

func (this *T_Ban_Player_Table) UpdateWithFieldName(t *T_Ban_Player, fields_name []string) {
	var field_list = t.GetFVPList(fields_name)
	if field_list != nil {
		this.UpdateWithFVPList(t.Get_unique_id(), field_list)
	}
}

func (this *T_Ban_Player_Table) TransactionInsert(transaction *mysql_manager.Transaction, t *T_Ban_Player) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Insert("t_ban_player", field_list)
	}
}

func (this *T_Ban_Player_Table) TransactionDelete(transaction *mysql_manager.Transaction, unique_id string) {
	transaction.Delete("t_ban_player", "unique_id", unique_id)
}

func (this *T_Ban_Player_Table) TransactionUpdateAll(transaction *mysql_manager.Transaction, t*T_Ban_Player) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Update("t_ban_player", "unique_id", t.Get_unique_id(), field_list)
	}
}

func (this *T_Ban_Player_Table) TransactionUpdateWithFVPList(transaction *mysql_manager.Transaction, unique_id string, field_list []*mysql_base.FieldValuePair) {
	transaction.Update("t_ban_player", "unique_id", unique_id, field_list)
}

func (this *T_Ban_Player_Table) TransactionUpdateWithFieldName(transaction *mysql_manager.Transaction, t *T_Ban_Player, fields_name []string) {
	field_list := t.GetFVPList(fields_name)
	if field_list != nil {
		transaction.Update("t_ban_player", "unique_id", t.Get_unique_id(), field_list)
	}
}

func TransactionInsert_T_Ban_Player(transaction *mysql_manager.Transaction, t *T_Ban_Player) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Insert("t_ban_player", field_list)
	}
}

func TransactionDelete_T_Ban_Player(transaction *mysql_manager.Transaction, unique_id string) {
	transaction.Delete("t_ban_player", "unique_id", unique_id)
}

func TransactionUpdateAll_T_Ban_Player(transaction *mysql_manager.Transaction, t *T_Ban_Player) {
	field_list := t._format_field_list()
	if field_list != nil {
		transaction.Update("t_ban_player", "unique_id", t.Get_unique_id(), field_list)
	}
}

func TransactionUpdate_T_Ban_Player(transaction *mysql_manager.Transaction, t *T_Ban_Player, fields_name []string) {
	field_list := t.GetFVPList(fields_name)
	if field_list != nil {
		transaction.Update("t_ban_player", "unique_id", t.Get_unique_id(), field_list)
	}
}
