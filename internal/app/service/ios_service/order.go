package ios_service

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/ddmc/ios_session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/pkg/json"
)

func CheckOrder(cartMap map[string]any) (map[string]string, error) {
	api := "https://maicai.api.ddxq.mobi/order/checkOrder"

	pack := json.NewOrderMap(cartMap, []string{"reserved_time", "real_match_supply_order", "is_supply_order", "package_type", "package_id", "products"})

	products := cartMap["products"].([]any)
	productFields := []string{"id", "is_booking", "total_money", "is_invoice", "total_origin_money", "category_path", "count", "type", "batch_type", "is_coupon_gift", "price", "order_sort", "instant_rebate_money", "activity_id", "conditions_num", "price_type", "product_type", "origin_price"}
	list := make([]any, 0, len(products))
	for _, product := range products {
		val, ok := product.(map[string]any)
		if !ok {
			continue
		}
		item := make(map[string]any)
		for _, v := range productFields {
			item[v] = val[v]
		}
		if v, ok := val["is_invoice"].(float64); ok {
			item["is_invoice"] = v > 0
		}
		item["batch_type"] = 0
		if v, ok := val["sale_batches"].(map[string]any); ok {
			item["batch_type"] = v["batch_type"]
		}
		if v, ok := val["order_sort"].(float64); ok {
			item["order_sort"] = strconv.Itoa(int(v))
		}
		if v, ok := val["price_type"].(float64); ok {
			item["price_type"] = strconv.Itoa(int(v))
		}
		list = append(list, item)
	}
	pack.Map["products"] = list
	pack.Map["reserved_time"] = map[string]any{
		"time_biz_type": 0,
	}
	packages := json.OrderMaps{pack}.MarshalToString()

	now := time.Now().Unix()
	// now = 1651583468 // mock 测试
	headers := ios_session.TakeHeaders(now, []string{"accept", "accept-encoding", "accept-language", "content-type", "cookie", "ddmc-api-version", "ddmc-app-client-id", "ddmc-build-version", "ddmc-channel", "ddmc-city-number", "ddmc-country-code", "ddmc-device-id", "ddmc-device-model", "ddmc-device-name", "ddmc-device-token", "ddmc-idfa", "ddmc-ip", "ddmc-language-code", "ddmc-latitude", "ddmc-locale-identifier", "ddmc-longitude", "ddmc-os-version", "ddmc-station-id", "ddmc-uid", "user-agent", "im_secret"})

	params := ios_session.TakeParams(now, []string{"api_version", "app_client_id", "app_type", "buildVersion", "channel", "city_number", "countryCode", "device_id", "device_model", "device_name", "device_token", "idfa", "ip", "languageCode", "latitude", "localeIdentifier", "longitude", "os_version", "station_id", "uid"})
	params["seqid"] = "3358795160"
	params["ab_config"] = `{"ETA_time_default_selection":"C1.1"}`
	params["address_id"] = ios_session.Address().Id
	params["coupons_id"] = ""
	params["freight_ticket_id"] = "default"
	params["user_ticket_ids"] = "default"
	params["is_buy_coupons"] = "0"
	params["is_buy_vip"] = "0"
	params["is_use_point"] = "0"
	params["is_use_balance"] = "0"
	params["check_order_type"] = "0"
	params["packages"] = packages
	signed, err := ios_session.Sign(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = signed["sign"]
	headers["sign"] = signed["sign"]
	headers["nars"] = signed["nars"]
	headers["sesi"] = signed["sesi"]

	// body := ios_session.EncodeFormDataToString(params)
	// log.Println(body)
	result := new(dto.Result)
	errMsg := new(dto.ErrorMessage)
	resp, err := ios_session.Client().R().
		SetHeaders(headers).
		// SetBody(body).
		SetFormData(params).
		SetResult(result).
		SetError(errMsg).
		// SetRetryCount(50).
		Send(http.MethodPost, api)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	if !resp.IsSuccess() {
		return nil, errs.Wrap(errs.CheckOrderFailed, resp.String())
	}
	if !result.Success {
		return nil, errs.Wrap(errs.CheckOrderFailed, resp.String())
	}
	data, err := resp.ToBytes()
	if err != nil {
		return nil, errs.WithStack(err)
	}

	order := json.Get(data, "data", "order")
	// log.Println("data.order", order.ToString())
	freight := order.Get("freights", 0, "freight")
	res := map[string]string{
		"price":                  order.Get("total_money").ToString(),              // 总价
		"freight_discount_money": freight.Get("discount_freight_money").ToString(), // 运费折扣
		"freight_money":          freight.Get("freight_money").ToString(),          // 运费
		"order_freight":          freight.Get("freight_real_money").ToString(),     // 订单真实运费
	}
	if order.Get("default_coupon", "_id").ToString() != "" {
		res["user_ticket_id"] = order.Get("default_coupon", "_id").ToString() // 优惠券id
	}
	log.Println("[叮咚]订单总金额 =>", res["price"])
	return res, nil
}

func AddNewOrder(cartMap map[string]any, reserveTime *reserve_time.GoTimes, checkOrderMap map[string]string) error {
	api := "https://maicai.api.ddxq.mobi/order/addNewOrder"

	now := time.Now().Unix()
	headers := ios_session.TakeHeaders(now, []string{"ddmc-os-version", "ddmc-city-number", "user-agent", "ddmc-locale-identifier", "ddmc-device-token", "cookie", "ddmc-api-version", "ddmc-build-version", "ddmc-idfa", "ddmc-longitude", "ddmc-latitude", "ddmc-app-client-id", "content-length", "ddmc-uid", "ddmc-device-name", "accept-language", "ddmc-device-model", "ddmc-channel", "ddmc-device-id", "ddmc-country-code", "ddmc-ip", "accept-encoding", "content-type", "ddmc-language-code", "ddmc-station-id", "accept"})

	paymentOrder := map[string]any{
		"current_position":    []string{headers["ddmc-latitude"], headers["ddmc-longitude"]},
		"receipt_without_sku": "0",
		"is_use_balance":      "0",
		"order_type":          1,
		"parent_order_sign":   cartMap["parent_order_sign"],
		"used_point_num":      0,
		"pay_type":            4,
		"address_id":          ios_session.Address().Id,
		"user_ticket_id":      checkOrderMap["user_ticket_id"],
		"price":               checkOrderMap["price"],
		"order_freight":       checkOrderMap["order_freight"],
	}
	conf := config.GetDingDong()
	if conf.PayType > 0 {
		paymentOrder["pay_type"] = conf.PayType
	}

	pack := map[string]any{
		"soon_arrival":            0,
		"reserved_time_start":     reserveTime.StartTimestamp,
		"eta_trace_id":            "",
		"package_id":              cartMap["package_id"],
		"package_type":            cartMap["package_type"],
		"reserved_time_end":       reserveTime.EndTimestamp,
		"time_biz_type":           0,
		"real_match_supply_order": false,
		"first_selected_big_time": "1",
	}
	items := cartMap["products"].([]any)
	products := make([]map[string]any, 0, len(items))
	for _, val := range items {
		item, ok := val.(map[string]any)
		if !ok {
			continue
		}
		product := map[string]any{
			"id":         item["id"],
			"count":      item["count"],
			"price":      item["price"],
			"parent_id":  item["parent_id"],
			"cart_id":    item["cart_id"],
			"order_sort": item["order_sort"],
			"batch_type": 0,
		}
		if v, ok := item["sale_batches"].(map[string]any); ok {
			product["batch_type"] = v["batch_type"]
		}
		products = append(products, product)
	}
	pack["products"] = products

	packageOrder := map[string]any{
		"payment_order": paymentOrder,
		"packages":      []any{pack},
	}
	packageOrderJson := json.MustEncodeFast(packageOrder)

	params := ios_session.TakeParams(now, []string{"api_version", "app_client_id", "app_type", "buildVersion", "channel", "city_number", "countryCode", "device_id", "device_model", "device_name", "device_token", "idfa", "ip", "languageCode", "latitude", "localeIdentifier", "longitude", "os_version", "seqid", "station_id", "uid"})
	params["ab_config"] = `{"key_no_condition_barter":true}`
	params["clientDetail"] = ""
	params["package_order"] = packageOrderJson
	// log.Printf("AddNewOrder params => %#v", params)
	signed, err := ios_session.Sign(params)
	if err != nil {
		return err
	}
	params["sign"] = signed["sign"]
	headers["sign"] = signed["sign"]
	headers["nars"] = signed["nars"]
	headers["sesi"] = signed["sesi"]

	result := new(dto.Result)
	errMsg := new(dto.ErrorMessage)
	resp, err := ios_session.Client().R().
		SetHeaders(headers).
		SetFormData(params).
		SetResult(result).
		SetError(errMsg).
		// SetRetryCount(50).
		Send(http.MethodPost, api)
	if err != nil {
		return errs.WithStack(err)
	}
	if !resp.IsSuccess() {
		return errs.Wrap(errs.SubmitOrderFailed, resp.String())
	}
	if !result.Success {
		if result.Code == 5004 {
			return errs.WithStack(errs.ReserveTimeIsDisabled)
		}
		return errs.Wrap(errs.SubmitOrderFailed, resp.String())
	}
	log.Println("[叮咚]恭喜你，已成功下单 =>", resp.String())
	return nil
}
