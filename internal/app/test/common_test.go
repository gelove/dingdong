package test

import (
	"testing"

	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/pkg/js"
	"dingdong/pkg/json"
)

func TestJsCall(t *testing.T) {
	headers := session.GetHeaders()
	headers["accept-language"] = "en-us"
	params := session.GetParams(headers)
	params["tab_type"] = "1"
	params["page"] = "1"
	res, err := js.Call(jsFile, "sign", json.MustEncodeToString(params))
	if err != nil {
		t.Error("js parser error =>", err)
		return
	}
	// nars 对应就可以
	// sesi 可以不用管, 依赖 nars 与 一个随机字符串, 每次计算应都不同, 但是在JS虚拟机中伪随机数似乎不变, 每次都会得到同一个值
	t.Log("value =>", res.String())
}

func TestJsonGet(t *testing.T) {
	conf := config.Get()
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
	t.Log(json.MustEncodeToString(out))
}

func TestSnapUpTime(t *testing.T) {
	t.Log(date.FirstSnapUpTime())
	t.Log(date.SecondSnapUpTime())
}
