package main

import (
	"ih_server_new/src/db_gen/game_db"
	"sync"
)

type BanPlayerMgr struct {
	ban_players sync.Map
}

var ban_player_mgr BanPlayerMgr

func (this *BanPlayerMgr) GetLocal(unique_id string) *game_db.T_Ban_Player {
	v, o := this.ban_players.Load(unique_id)
	if !o {
		return nil
	}
	return v.(*game_db.T_Ban_Player)
}

func (this *BanPlayerMgr) SetLocal(unique_id string, ban_player *game_db.T_Ban_Player) {
	this.ban_players.Store(unique_id, ban_player)
}

func (this *BanPlayerMgr) Get(unique_id string) *game_db.T_Ban_Player {
	v, o := this.ban_players.Load(unique_id)
	if !o {
		v = ban_player_table.SelectByPrimaryField(unique_id)
		if v == nil {
			return nil
		}
		this.ban_players.Store(unique_id, v)
	}
	return v.(*game_db.T_Ban_Player)
}
