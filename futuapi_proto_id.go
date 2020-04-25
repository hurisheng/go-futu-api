package futuapi

// ProtoID for all the commands.
const (
	ProtoIDInitConnect    = 1001
	ProtoIDGetGlobalState = 1002
	ProtoIDNotify         = 1003
	ProtoIDKeepAlive      = 1004

	ProtoIDTrdGetAccList  = 2001
	ProtoIDTrdUnlockTrade = 2005
	ProtoIDTrdSubAccPush  = 2008

	ProtoIDTrdGetFunds        = 2101
	ProtoIDTrdGetPositionList = 2102
	ProtoIDTrdGetMaxTrdQtys   = 2111

	ProtoIDTrdGetOrderList            = 2201
	ProtoIDTrdPlaceOrder              = 2202
	ProtoIDTrdModifyOrder             = 2205
	ProtoIDTrdUpdateOrder             = 2208
	ProtoIDTrdGetOrderFillList        = 2211
	ProtoIDTrdUpdateOrderFill         = 2218
	ProtoIDTrdGetHistoryOrderList     = 2221
	ProtoIDTrdGetHistoryOrderFillList = 2222

	ProtoIDQotSub                 = 3001
	ProtoIDQotRegQotPush          = 3002
	ProtoIDQotGetSubInfo          = 3003
	ProtoIDQotGetBasicQot         = 3004
	ProtoIDQotUpdateBasicQot      = 3005
	ProtoIDQotGetKL               = 3006
	ProtoIDQotUpdateKL            = 3007
	ProtoIDQotGetRT               = 3008
	ProtoIDQotUpdateRT            = 3009
	ProtoIDQotGetTicker           = 3010
	ProtoIDQotUpdateTicker        = 3011
	ProtoIDQotGetOrderBook        = 3012
	ProtoIDQotUpdateOrderBook     = 3013
	ProtoIDQotGetBroker           = 3014
	ProtoIDQotUpdateBroker        = 3015
	ProtoIDQotUpdatePriceReminder = 3019

	ProtoIDQotRequestHistoryKL      = 3103
	ProtoIDQotRequestHistoryKLQuota = 3104
	ProtoIDQotRequestRehab          = 3105

	ProtoIDQotGetStaticInfo          = 3202
	ProtoIDQotGetSecuritySnapshot    = 3203
	ProtoIDQotGetPlateSet            = 3204
	ProtoIDQotGetPlateSecurity       = 3205
	ProtoIDQotGetReference           = 3206
	ProtoIDQotGetOwnerPlate          = 3207
	ProtoIDQotGetHoldingChangeList   = 3208
	ProtoIDQotGetOptionChain         = 3209
	ProtoIDQotGetWarrant             = 3210
	ProtoIDQotGetCapitalFlow         = 3211
	ProtoIDQotGetCapitalDistribution = 3212
	ProtoIDQotGetUserSecurity        = 3213
	ProtoIDQotModifyUserSecurity     = 3214
	ProtoIDQotStockFilter            = 3215
	ProtoIDQotGetCodeChange          = 3216
	ProtoIDQotGetIpoList             = 3217
	ProtoIDQotGetFutureInfo          = 3218
	ProtoIDQotRequestTradeDate       = 3219
	ProtoIDQotSetPriceReminder       = 3220
	ProtoIDQotGetPriceReminder       = 3221
)
