package main

import (
	"errors"
	"fmt"
	"ih_server_new/libs/log"
	"ih_server_new/libs/rpc"
	"ih_server_new/src/rpc_proto"
	"time"
)

// 大厅到大厅的调用
type H2H_CallProc struct {
}

func (this *H2H_CallProc) Do(args *rpc_proto.H2R_Transfer, reply *rpc_proto.H2R_TransferResult) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()
	rpc_client := GetRpcClientByPlayerId(args.ReceivePlayerId)
	if rpc_client == nil {
		err_str := fmt.Sprintf("!!!!!! Not found rpc client by player id %v", args.ReceivePlayerId)
		return errors.New(err_str)
	}

	log.Debug("H2H_CallProc @@@@@@@ call method[%v] args[%v] reply[%v]", args.Method, args.Args, reply.Result)

	var result interface{}
	err := rpc_client.Call(args.Method, args.Args, result)
	if err != nil {
		return err
	}
	log.Debug("H2H_CallProc @@@@@@@ call method[%v] result[%v]", args.Method, result)
	reply.Result = result
	return nil
}

// ping 大厅
type H2R_PingProc struct {
}

func (this *H2R_PingProc) Do(args *rpc_proto.H2R_Ping, result *rpc_proto.H2R_Pong) error {
	// 不做任何处理
	return nil
}

/* 监听RPC调用 */
type H2R_ListenRPCProc struct {
}

func (this *H2R_ListenRPCProc) Do(args *rpc_proto.H2R_ListenIPNoitfy, result *rpc_proto.H2R_ListenIPResult) error {
	defer func() {
		if err := recover(); err != nil {
			log.Stack(err)
		}
	}()

	log.Info("get notify listen rpc ip: %v", args.ListenIP)
	// 再连接到HallServer

	if !server.connect_hall(args.ListenIP, args.ServerId) {
		err_str := fmt.Sprintf("不能连接到大厅[IP:%v, Id:%v]", args.ListenIP, args.ServerId)
		return errors.New(err_str)
	}

	time.Sleep(time.Second * 1)
	return nil
}

// 全局调用
type H2R_GlobalProc struct {
}

func (this *H2R_GlobalProc) ChargeSave(args *rpc_proto.H2R_ChargeSave, result *rpc_proto.H2R_ChargeSaveResult) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	if args.Channel == 1 {
		row := google_pay_table.NewRow(args.OrderId)
		row.Set_bundle_id(args.BundleId)
		row.Set_account(args.Account)
		row.Set_player_id(args.PlayerId)
		row.Set_pay_time(args.PayTime)
		row.Set_pay_time_str(args.PayTimeStr)
		google_pay_table.InsertIgnore(row)
	} else if args.Channel == 2 {
		row := apple_pay_table.NewRow(args.OrderId)
		row.Set_bundle_id(args.BundleId)
		row.Set_account(args.Account)
		row.Set_player_id(args.PlayerId)
		row.Set_pay_time(args.PayTime)
		row.Set_pay_time_str(args.PayTimeStr)
		apple_pay_table.InsertIgnore(row)
	} else {
		err_str := fmt.Sprintf("@@@ H2R_GlobalProc::ChargeSave Player[%v,%v], Unknown Channel %v", args.Account, args.PlayerId, args.Channel)
		return errors.New(err_str)
	}

	log.Trace("@@@ Charge Save %v", args)

	return nil
}

// 初始化
func (this *RpcServer) init_proc_service() bool {
	this.rpc_service = &rpc.Service{}

	if !this.rpc_service.Register(&H2H_CallProc{}) {
		return false
	}

	if !this.rpc_service.Register(&H2R_ListenRPCProc{}) {
		return false
	}

	if !this.rpc_service.Register(&H2R_GlobalProc{}) {
		return false
	}

	if !this.rpc_service.Register(&G2G_CommonProc{}) {
		return false
	}

	// 注册用户自定义RPC数据类型
	rpc_proto.RegisterRpcUserType()

	if this.rpc_service.Listen(config.ListenIP) != nil {
		return false
	}
	return true
}

// 反初始化
func (this *RpcServer) uninit_proc_service() {
	if this.rpc_service != nil {
		this.rpc_service.Close()
		this.rpc_service = nil
	}
}
