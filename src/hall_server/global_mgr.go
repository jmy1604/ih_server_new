package main

import (
	"ih_server_new/src/share_data"
	"sync"
)

type GlobalMgr struct {
	player_id_locker        sync.RWMutex
	guild_id_locker         sync.RWMutex
	battle_record_id_locker sync.RWMutex
	mail_id_locker          sync.RWMutex
}

var global_mgr GlobalMgr

func (this *GlobalMgr) GetNextPlayerId() int32 {
	this.player_id_locker.Lock()
	defer this.player_id_locker.Unlock()

	curr_player_id := global_table.GetRow().Get_curr_player_id()
	new_id := share_data.GeneratePlayerId(config.ServerId, curr_player_id+1)
	global_table.GetRow().Set_curr_player_id(curr_player_id + 1)
	global_table.UpdateWithFieldName(global_table.GetRow(), []string{"curr_player_id"})
	return new_id
}

func (this *GlobalMgr) GetNextGuildId() int32 {
	this.guild_id_locker.Lock()
	defer this.guild_id_locker.Unlock()

	curr_guild_id := global_table.GetRow().Get_curr_guild_id()
	new_id := share_data.GenerateGuildId(config.ServerId, curr_guild_id+1)
	global_table.GetRow().Set_curr_guild_id(curr_guild_id + 1)
	global_table.UpdateWithFieldName(global_table.GetRow(), []string{"curr_guild_id"})
	return new_id
}

func (this *GlobalMgr) GetNextBattleRecordId() int32 {
	this.battle_record_id_locker.Lock()
	defer this.battle_record_id_locker.Unlock()

	curr_record_id := global_table.GetRow().Get_curr_record_id()
	next_record_id := curr_record_id + 1
	global_table.GetRow().Set_curr_record_id(next_record_id)
	global_table.UpdateWithFieldName(global_table.GetRow(), []string{"curr_record_id"})
	return next_record_id
}
