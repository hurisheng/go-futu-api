package futuapi

import (
	"fmt"

	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetAccList"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetFunds"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetHistoryOrderFillList"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetHistoryOrderList"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetMaxTrdQtys"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetOrderList"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_GetPositionList"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_ModifyOrder"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_PlaceOrder"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_SubAccPush"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_UnlockTrade"
	"github.com/hurisheng/go-futu-api/protobuf/Trd_UpdateOrder"
)

// TrdGetAccList 获取交易账户列表
func (api *FutuAPI) TrdGetAccList(req *Trd_GetAccList.Request) (<-chan *Trd_GetAccList.Response, error) {
	out := make(chan *Trd_GetAccList.Response)
	if err := api.send(ProtoIDTrdGetAccList, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetAccList error: %w", err)
	}
	return out, nil
}

// TrdUnlockTrade 解锁或锁定交易
func (api *FutuAPI) TrdUnlockTrade(req *Trd_UnlockTrade.Request) (<-chan *Trd_UnlockTrade.Response, error) {
	out := make(chan *Trd_UnlockTrade.Response)
	if err := api.send(ProtoIDTrdUnlockTrade, req, out); err != nil {
		return nil, fmt.Errorf("TrdUnlockTrade error: %w", err)
	}
	return out, nil
}

// TrdSubAccPush 订阅接收交易账户的推送数据
func (api *FutuAPI) TrdSubAccPush(req *Trd_SubAccPush.Request) (<-chan *Trd_SubAccPush.Response, error) {
	out := make(chan *Trd_SubAccPush.Response)
	if err := api.send(ProtoIDTrdSubAccPush, req, out); err != nil {
		return nil, fmt.Errorf("TrdSubAccPush error: %w", err)
	}
	return out, nil
}

// TrdGetFunds 获取账户资金
func (api *FutuAPI) TrdGetFunds(req *Trd_GetFunds.Request) (<-chan *Trd_GetFunds.Response, error) {
	out := make(chan *Trd_GetFunds.Response)
	if err := api.send(ProtoIDTrdGetFunds, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetFunds error: %w", err)
	}
	return out, nil
}

// TrdGetPositionList 获取持仓列表
func (api *FutuAPI) TrdGetPositionList(req *Trd_GetPositionList.Request) (<-chan *Trd_GetPositionList.Response, error) {
	out := make(chan *Trd_GetPositionList.Response)
	if err := api.send(ProtoIDTrdGetPositionList, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetPositionList error: %w", err)
	}
	return out, nil
}

// TrdGetMaxTrdQtys 获取最大交易数量
func (api *FutuAPI) TrdGetMaxTrdQtys(req *Trd_GetMaxTrdQtys.Request) (<-chan *Trd_GetMaxTrdQtys.Response, error) {
	out := make(chan *Trd_GetMaxTrdQtys.Response)
	if err := api.send(ProtoIDTrdGetMaxTrdQtys, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetMaxTrdQtys error: %w", err)
	}
	return out, nil
}

// TrdGetOrderList 获取订单列表
func (api *FutuAPI) TrdGetOrderList(req *Trd_GetOrderList.Request) (<-chan *Trd_GetOrderList.Response, error) {
	out := make(chan *Trd_GetOrderList.Response)
	if err := api.send(ProtoIDTrdGetOrderList, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetOrderList error: %w", err)
	}
	return out, nil
}

// TrdPlaceOrder 下单
func (api *FutuAPI) TrdPlaceOrder(req *Trd_PlaceOrder.Request) (<-chan *Trd_PlaceOrder.Response, error) {
	out := make(chan *Trd_PlaceOrder.Response)
	if err := api.send(ProtoIDTrdPlaceOrder, req, out); err != nil {
		return nil, fmt.Errorf("TrdPlaceOrder error: %w", err)
	}
	return out, nil
}

// TrdModifyOrder 修改订单(改价、改量、改状态等)
func (api *FutuAPI) TrdModifyOrder(req *Trd_ModifyOrder.Request) (<-chan *Trd_ModifyOrder.Response, error) {
	out := make(chan *Trd_ModifyOrder.Response)
	if err := api.send(ProtoIDTrdModifyOrder, req, out); err != nil {
		return nil, fmt.Errorf("TrdModifyOrder error: %w", err)
	}
	return out, nil
}

// TrdUpdateOrder 推送订单更新
func (api *FutuAPI) TrdUpdateOrder() (<-chan *Trd_UpdateOrder.Response, error) {
	out := make(chan *Trd_UpdateOrder.Response)
	if err := api.subscribe(ProtoIDTrdUpdateOrder, out); err != nil {
		return nil, fmt.Errorf("TrdUpdateOrder error: %w", err)
	}
	return out, nil
}

// TrdGetOrderFillList 获取成交列表
func (api *FutuAPI) TrdGetOrderFillList(req *Trd_GetOrderList.Request) (<-chan *Trd_GetOrderList.Response, error) {
	out := make(chan *Trd_GetOrderList.Response)
	if err := api.send(ProtoIDTrdGetOrderFillList, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetOrderFillList error: %w", err)
	}
	return out, nil
}

// TrdUpdateOrderFill 推送新成交
func (api *FutuAPI) TrdUpdateOrderFill() (<-chan *Trd_UpdateOrder.Response, error) {
	out := make(chan *Trd_UpdateOrder.Response)
	if err := api.subscribe(ProtoIDTrdUpdateOrderFill, out); err != nil {
		return nil, fmt.Errorf("TrdUpdateOrderFill error: %w", err)
	}
	return out, nil
}

// TrdGetHistoryOrderList 获取历史订单列表
func (api *FutuAPI) TrdGetHistoryOrderList(req *Trd_GetHistoryOrderList.Request) (<-chan *Trd_GetHistoryOrderList.Response, error) {
	out := make(chan *Trd_GetHistoryOrderList.Response)
	if err := api.send(ProtoIDTrdGetHistoryOrderList, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetHistoryOrderList error: %w", err)
	}
	return out, nil
}

// TrdGetHistoryOrderFillList 获取历史成交列表
func (api *FutuAPI) TrdGetHistoryOrderFillList(req *Trd_GetHistoryOrderFillList.Request) (<-chan *Trd_GetHistoryOrderFillList.Response, error) {
	out := make(chan *Trd_GetHistoryOrderFillList.Response)
	if err := api.send(ProtoIDTrdGetHistoryOrderFillList, req, out); err != nil {
		return nil, fmt.Errorf("TrdGetHistoryOrderFillList error: %w", err)
	}
	return out, nil
}
