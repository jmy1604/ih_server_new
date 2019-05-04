package main

import (
	"ih_server_new/src/db_gen/login_db"
	"sync"
)

type BanMgr struct {
	ban_players sync.Map
}

var ban_mgr BanMgr

func (this *BanMgr) Get(unique_id string) *login_db.T_Ban_Player {
	v, o := this.ban_players.Load(unique_id)
	if !o {
		return nil
	}
	return v.(*login_db.T_Ban_Player)
}

func (this *BanMgr) GetAndSet(unique_id string) *login_db.T_Ban_Player {
	v, o := this.ban_players.Load(unique_id)
	if !o {
		v, o = ban_player_table.Select("unique_id", unique_id)
		if !o {
			return nil
		}
	}
	return v.(*login_db.T_Ban_Player)
}

func (this *BanMgr) Set(unique_id string, ban_player *login_db.T_Ban_Player) bool {
	_, o := this.ban_players.LoadOrStore(unique_id, ban_player)
	if o {
		return false
	}
	return true
}
