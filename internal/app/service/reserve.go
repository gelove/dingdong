package service

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

var lastNotify time.Time

func MockCartMap() map[string]interface{} {
	first := make(map[string]interface{})
	cartStr := `{"products":[{"type":1,"id":"612cc0982c34fab505117d4e","price":"828.00","count":1,"description":"","sizes":[],"cart_id":"612cc0982c34fab505117d4e","parent_id":"","parent_batch_type":-1,"category_path":"","manage_category_path":"411,412,413","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"洋河蓝色经典梦之蓝M6+52度白酒 550ml/瓶","product_type":0,"small_image":"https://ddfs-public.ddimg.mobi/img/blind/product-management/202108/1242efbb2a37470aa081683513fb3677.jpg?width=800&height=800","total_price":"828.00","origin_price":"828.00","total_origin_price":"828.00","no_supplementary_price":"828.00","no_supplementary_total_price":"828.00","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"瓶","net_weight":"550","net_weight_unit":"ml","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":1,"is_presale":0}]}`
	json.MustDecodeFromString(cartStr, &first)
	productList := first["products"].([]interface{})
	products := make([]map[string]interface{}, 0, len(productList))
	for _, v := range productList {
		product := v.(map[string]interface{})
		product["total_money"] = product["total_price"]
		product["total_origin_money"] = product["total_origin_price"]
		products = append(products, product)
	}
	cartMap := make(map[string]interface{})
	cartMap["products"] = products
	return cartMap
}

func GetMultiReserveTime(cartMap map[string]interface{}) (*reserve_time.GoTimes, error) {
	url := "https://maicai.api.ddxq.mobi/order/getMultiReserveTime"

	products := cartMap["products"].([]map[string]interface{})
	productsList := [][]map[string]interface{}{products}
	productsJson := json.MustEncodeToString(productsList)

	conf := config.Get()
	headers := session.GetHeaders()
	// 响应压缩有乱码 暂不压缩
	// headers["accept-encoding"] = "gzip, deflate, br"
	params := session.GetParams(headers)
	params["group_config_id"] = ""
	params["isBridge"] = "false"
	params["address_id"] = conf.Params["address_id"]
	params["products"] = productsJson

	form, err := session.Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := reserve_time.Result{}
	_, err = session.Client().R().
		SetHeaders(headers).
		SetFormData(form).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodPost, url)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	// log.Println("resp =>", resp.String())
	// log.Println("result =>", json.MustEncodeToString(result))
	if !result.Success {
		return nil, errs.WithMessage(code.InvalidResponse, result.Msg)
	}
	if len(result.Data) == 0 || len(result.Data[0].Times) == 0 || len(result.Data[0].Times[0].Times) == 0 {
		return nil, errs.New(code.NoReserveTime)
	}

	times := result.Data[0].Times[0].Times
	validTimes := filterValidTimes(times)
	if len(validTimes) == 0 {
		return nil, errs.New(code.NoReserveTimeAndRetry)
	}
	log.Println("发现可用的配送时段, 请尽快下单!")
	return validTimes[0], nil
}

func GetMultiReserveTimeAndNotify(cartMap map[string]interface{}) {
	_, err := GetMultiReserveTime(cartMap)
	if err != nil {
		log.Println(err)
		return
	}

	conf := config.Get()
	if conf.NotifyNeeded {
		notify()
	}
}

func notify() {
	now := time.Now()
	conf := config.Get()
	interval := time.Duration(conf.NotifyInterval) * time.Minute
	if now.Before(lastNotify.Add(interval)) {
		return
	}
	lastNotify = now
	list, err := GetHomeFlowDetail()
	if err != nil {
		log.Printf("获取首页流量失败 => %+v", err)
	}
	productNames := make([]string, 0, 10)
	for i, item := range list {
		if i >= 10 {
			continue
		}
		letter := []rune(item.Name)
		if len(letter) > 5 {
			productNames = append(productNames, string(letter[:5]))
			continue
		}
		productNames = append(productNames, string(letter))
	}
	ellipsis := ""
	if len(list) >= 10 {
		ellipsis = "..."
	}
	products := strings.Join(productNames, " ")
	wg := new(sync.WaitGroup)
	for _, v := range conf.Users {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			Push(key, products, ellipsis)
		}(v)
	}
	wg.Wait()
}

func filterValidTimes(times []*reserve_time.GoTimes) []*reserve_time.GoTimes {
	var validTimes []*reserve_time.GoTimes
	for _, v := range times {
		if v.DisableType != 0 {
			continue
		}
		validTimes = append(validTimes, v)
	}
	return validTimes
}
