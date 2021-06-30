package futuapi

import (
	"context"
	"testing"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
)

func TestConnect(t *testing.T) {
	api := NewFutuAPI()
	defer api.Close(context.Background())

	api.SetRecvNotify(true)
	nCh, err := api.SysNotify(context.Background())
	if err != nil {
		t.Error(err)
	}

	if err := api.Connect(context.Background(), ":11111"); err != nil {
		t.Error(err)
		return
	}

	if sub, err := api.QuerySubscription(context.Background(), true); err != nil {
		t.Error(err)
	} else {
		t.Error(sub)
	}

	tCh, err := api.UpdateTicker(context.Background())
	if err != nil {
		t.Error(err)
	}
	if err := api.Subscribe(context.Background(), []*Security{{qotcommon.QotMarket_QotMarket_HK_Security, "00700"}}, []qotcommon.SubType{qotcommon.SubType_SubType_Ticker}, true, true, true, true); err != nil {
		t.Error(err)
	}
	select {
	case notify := <-nCh:
		t.Error(notify)
	case ticker := <-tCh:
		t.Error(ticker)
	}

	if sub, err := api.QuerySubscription(context.Background(), true); err != nil {
		t.Error(err)
	} else {
		t.Error(sub)
	}

	secs, err := api.GetUserSecurity(context.Background(), "全部")
	if err != nil {
		t.Error(err)
	} else {
		t.Error(secs)
	}
}
