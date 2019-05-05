package main

import (
	"ih_server_new/src/db_gen/login_db"
	"sync"
)

type AccountMgr struct {
	accounts sync.Map
}

var acc_mgr AccountMgr

func (this *AccountMgr) GetLocal(acc string) *login_db.T_Account {
	v, o := this.accounts.Load(acc)
	if !o {
		return nil
	}
	return v.(*login_db.T_Account)
}

func (this *AccountMgr) Get(acc string) *login_db.T_Account {
	v, o := this.accounts.Load(acc)
	if !o {
		v, o = account_table.Select("account", acc)
		if !o {
			return nil
		}
		this.accounts.Store(acc, v)
	}
	return v.(*login_db.T_Account)
}

func (this *AccountMgr) SetLocal(acc string, account_row *login_db.T_Account) bool {
	_, o := this.accounts.LoadOrStore(acc, account_row)
	if o {
		return false
	}
	return true
}
