package common

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/pkg/errs"
	"dingdong/pkg/json"
	"dingdong/pkg/textual"
	"dingdong/pkg/uri"
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

func TestErrors(t *testing.T) {
	err := errs.Wrap(errs.ReserveTimeIsDisabled, "test")
	t.Log(err)
	t.Logf("%+v", err)
	if !errs.Is(err, errs.ReserveTimeIsDisabled) {
		t.Error("error is not equal")
		return
	}
	t.Log("error is equal")

	err = errs.NoReserveTime
	// 判断是否为同一种错误类型
	if !errs.As(err, &errs.ReserveTimeIsDisabled) {
		t.Error("error is not equal")
		return
	}
	t.Log("error type is equal")
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

func TestSort(t *testing.T) {
	list := []string{"api_version", "app_client_id", "app_type", "buildVersion", "channel", "city_number", "countryCode", "device_id", "device_model", "device_name", "device_token", "idfa", "ip", "languageCode", "latitude", "localeIdentifier", "longitude", "os_version", "seqid", "station_id", "time", "uid", "address_id"}
	sort.Strings(list)
	t.Logf("%#v", list)
}

func TestJsonGet(t *testing.T) {
	conf := config.GetDingDong()
	bs := json.MustEncode(conf)
	t.Log(json.Get(bs, "headers", "cookie").ToString())
	t.Log(json.Get(bs, "headers").Get("cookie").ToString())
}

func TestUrlEncode(t *testing.T) {
	str := `{"ETA_time_default_selection":"C1.2"}`
	str = uri.QueryEscape(str)
	if str != "%7B%22ETA_time_default_selection%22%3A%22C1.2%22%7D" {
		t.Error("url encode error")
		return
	}

	str = "iPhone13,2iPhone 12"
	str = uri.QueryEscape(str)
	if str != "iPhone13%2C2iPhone%2012" {
		t.Error("url encode error")
		return
	}

	str = "BP1NV/qf1rudTeHnT1hrHfvJ+rhLNEieZtrTeMs2yM6qOgBRKaTr0oa+RmY8y5YD8nrPlKXUY5xfE+q6YudWclw=="
	str = uri.QueryEscape(str)
	if str != "BP1NV/qf1rudTeHnT1hrHfvJ%2BrhLNEieZtrTeMs2yM6qOgBRKaTr0oa%2BRmY8y5YD8nrPlKXUY5xfE%2Bq6YudWclw%3D%3D" {
		t.Error("url encode error")
		return
	}

	expect := "%5B%22%5B%7B%5C%22sale_batches%5C%22%3A%7B%5C%22batch_type%5C%22%3A-1%7D%2C%5C%22is_coupon_gift%5C%22%3A0%2C%5C%22id%5C%22%3A%5C%226262964b929609dd27fc214d%5C%22%2C%5C%22price%5C%22%3A%5C%228.80%5C%22%2C%5C%22is_booking%5C%22%3A0%2C%5C%22count%5C%22%3A1%2C%5C%22small_image%5C%22%3A%5C%22https%3A%5C%5C%5C/%5C%5C%5C/imgnew.ddimg.mobi%5C%5C%5C/product%5C%5C%5C/b42d4fa545584bada451a43f1fedc0e3.jpg?width%3D800%26height%3D800%5C%22%2C%5C%22type%5C%22%3A1%2C%5C%22origin_price%5C%22%3A%5C%228.80%5C%22%2C%5C%22product_type%5C%22%3A0%2C%5C%22product_name%5C%22%3A%5C%22%E5%8F%AE%E5%92%9A%E4%BF%9D%E4%BE%9B%E5%AE%9A%E5%88%B6%E7%9A%87%E5%90%8E%E5%90%90%E5%8F%B8%20110g%5C%5C%5C/%E8%A2%8B%5C%22%7D%5D%22%5D"
	str = `[{"sale_batches":{"batch_type":-1},"is_coupon_gift":0,"id":"6262964b929609dd27fc214d","price":"8.80","is_booking":0,"count":1,"small_image":"https://imgnew.ddimg.mobi/product/b42d4fa545584bada451a43f1fedc0e3.jpg?width=800&height=800","type":1,"origin_price":"8.80","product_type":0,"product_name":"叮咚保供定制皇后吐司 110g/袋"}]`
	str = textual.ReplaceBatch(str, [2]string{`"`, `\"`}, [2]string{`/`, `\\\/`})
	str = uri.QueryEscape(fmt.Sprintf(`["%s"]`, str))
	if str != expect {
		t.Error("url encode error")
		return
	}
}

func TestJsonEncode(t *testing.T) {
	m := map[string]string{
		"reserved_time":           "",
		"real_match_supply_order": "",
		"is_supply_order":         "",
		"package_type":            "",
		"package_id":              "",
		"products":                "",
	}
	// 每次顺序不同
	t.Log(json.MustEncodeFast(m))
}

func TestJson(t *testing.T) {
	list := []map[string]any{
		{
			"sale_batches":   map[string]any{"batch_type": -1},
			"is_coupon_gift": 0,
			"id":             "62577c4c756fd5f16d9c623c",
			"price":          "16.80",
			"is_booking":     0,
			"count":          2,
			"small_image":    "https://imgnew.ddimg.mobi/product/bbc92baf0a0941b2a9755da5b4c1667d.jpg?width=800&height=800",
			"type":           1,
			"origin_price":   "16.80",
			"product_type":   0,
			"product_name":   "白玉小笼包（12只）300g/盒",
		},
	}
	str := json.MustEncodeWithUnescapeHTML(list)
	t.Logf("%s", json.MustEncodeWithUnescapeHTML([]string{str}))
}
