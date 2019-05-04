package login_db

import (
	"github.com/huoshan017/mysql-go/manager"
)

type TablesManager struct {
	db_T_Account_Table *T_Account_Table
	db_T_Ban_Player_Table *T_Ban_Player_Table
}

func (this *TablesManager) Init(db *mysql_manager.DB) {
	this.db_T_Account_Table = &T_Account_Table{}
	this.db_T_Account_Table.Init(db)
	this.db_T_Ban_Player_Table = &T_Ban_Player_Table{}
	this.db_T_Ban_Player_Table.Init(db)
}

func (this *TablesManager) Get_T_Account_Table() *T_Account_Table{
	return this.db_T_Account_Table
}

func (this *TablesManager) Get_T_Ban_Player_Table() *T_Ban_Player_Table{
	return this.db_T_Ban_Player_Table
}

func NewTablesManager(db *mysql_manager.DB) *TablesManager {
	tm := &TablesManager{}
	tm.Init(db)
	return tm
}
