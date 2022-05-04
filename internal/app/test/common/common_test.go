package common

import (
	"testing"
	"time"

	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

func BenchmarkSnapUpTime(b *testing.B) {
	// 启动内存统计
	b.ReportAllocs()

	b.Run("FirstSnapUpTime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = date.FirstSnapUpTime()
		}
	})

	// 性能更好
	b.Run("FirstSnapUpUnix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = date.FirstSnapUpUnix()
		}
	})
}

func TestError(t *testing.T) {
	err1 := errs.Wrap(code.Unexpected, errs.Wrap(code.ResponseError, errs.New(code.ReserveTimeIsDisabled)))
	t.Log(err1.Error())
	t.Log(err1.Message())
	err2 := errs.Wrap(code.Unexpected, errs.WithMessage(code.ResponseError, "提交订单失败"))
	t.Logf("%#v", err2)

	err := errs.New(code.ReserveTimeIsDisabled)
	second := errs.New(code.ReserveTimeIsDisabled)
	if !errs.As(err, &second) {
		t.Error("error is not equal")
		return
	}
	t.Log("error is equal")
}

func TestTimer(t *testing.T) {
	go func() {
		for {
			t.Log("timer =>", time.Now().Second())
			time.Sleep(time.Second)
		}
	}()
	<-time.After(time.Duration(60-1-time.Now().Second()) * time.Second)
	t.Log("timer finished")
}

func TestJsonGet(t *testing.T) {
	conf := config.GetDingDong()
	bs := json.MustEncode(conf)
	t.Log(json.Get(bs, "headers", "cookie").ToString())
	t.Log(json.Get(bs, "headers").Get("cookie").ToString())
}

func TestJsonEncode(t *testing.T) {
	str := `{
      "total_money": "19.00",
      "total_origin_money": "19.00",
      "goods_real_money": "19.00",
      "total_count": 2,
      "cart_count": 2,
      "is_presale": 0,
      "instant_rebate_money": "0.00",
      "total_rebate_money": "0.00",
      "used_balance_money": "0.00",
      "can_used_balance_money": "0.00",
      "used_point_num": 0,
      "used_point_money": "0.00",
      "can_used_point_num": 0,
      "can_used_point_money": "0.00",
      "is_share_station": 0,
      "only_today_products": [],
      "only_tomorrow_products": [],
      "package_type": 1,
      "package_id": 1,
      "front_package_text": "即时配送",
      "front_package_type": 0,
      "front_package_stock_color": "#2FB157",
      "front_package_bg_color": "#fbfefc",
      "is_supply_order": false,
      "eta_trace_id": "",
      "reserved_time_start": 1650594713,
      "reserved_time_end": 1650609000,
      "soon_arrival": "",
      "first_selected_big_time": 1
    }`
	out := make(map[string]interface{})
	json.MustDecodeFromString(str, &out)
	t.Log(json.MustEncodePrettyString(out))
}
