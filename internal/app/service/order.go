package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/textual"

	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/pkg/json"
)

func filterFields(data map[string]interface{}, fields []string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range data {
		if !textual.InArray(k, fields) {
			continue
		}
		res[k] = v
	}
	return res
}

func CheckOrder(cartMap map[string]interface{}, reserveTime *reserve_time.GoTimes) (map[string]interface{}, error) {
	api := "https://maicai.api.ddxq.mobi/order/checkOrder"

	packages := make(map[string]interface{})
	for key, val := range cartMap {
		if key == "products" {
			productFields := []string{"id", "category_path", "count", "price", "total_money", "instant_rebate_money", "activity_id", "conditions_num", "product_type", "sizes", "type", "total_origin_money", "price_type", "batch_type", "sub_list", "order_sort", "origin_price"}
			items := val.([]map[string]interface{})
			products := make([]map[string]interface{}, 0)
			for _, item := range items {
				products = append(products, filterFields(item, productFields))
			}
			packages[key] = products
			continue
		}
		packages[key] = val
	}
	packages["reserved_time"] = map[string]interface{}{
		"reserved_time_start": reserveTime.StartTimestamp,
		"reserved_time_end":   reserveTime.EndTimestamp,
	}
	packagesJson := json.MustEncodeToString([]interface{}{packages})

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["address_id"] = session.Address().Id
	params["user_ticket_id"] = "default"
	params["freight_ticket_id"] = "default"
	params["is_use_point"] = "0"
	params["is_use_balance"] = "0"
	params["is_buy_vip"] = "0"
	params["coupons_id"] = ""
	params["is_buy_coupons"] = "0"
	params["packages"] = packagesJson
	params["check_order_type"] = "0"
	params["is_support_merge_payment"] = "1"
	params["showData"] = "true"
	params["showMsg"] = "false"
	form, err := session.Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := dto.Result{}
	resp, err := session.Client().R().
		SetHeaders(headers).
		SetFormData(form).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodPost, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return nil, errs.WithMessage(code.ResponseError, "订单校验失败 => "+json.MustEncodeToString(result))
	}
	body, err := resp.ToBytes()
	if err != nil {
		return nil, errs.Wrap(code.ParseFailed, err)
	}

	order := json.Get(body, "data", "order")
	// log.Println("data.order", order.ToString())
	freight := order.Get("freights", 0, "freight")
	res := map[string]interface{}{
		"price":                  order.Get("total_money").ToString(),              // 总价
		"freight_discount_money": freight.Get("discount_freight_money").ToString(), // 运费折扣
		"freight_money":          freight.Get("freight_money").ToString(),          // 运费
		"order_freight":          freight.Get("freight_real_money").ToString(),     // 订单真实运费
	}
	if order.Get("default_coupon", "_id").ToString() != "" {
		res["user_ticket_id"] = order.Get("default_coupon", "_id").ToString() // 优惠券id
	}
	log.Println("订单总金额 =>", res["price"])
	return res, nil
}

func AddNewOrder(cartMap map[string]interface{}, reserveTime *reserve_time.GoTimes, checkOrderMap map[string]interface{}) error {
	api := "https://maicai.api.ddxq.mobi/order/addNewOrder"

	paymentOrder := map[string]interface{}{
		"reserved_time_start":    reserveTime.StartTimestamp,
		"reserved_time_end":      reserveTime.EndTimestamp,
		"parent_order_sign":      cartMap["parent_order_sign"],
		"address_id":             session.Address().Id,
		"pay_type":               6,
		"product_type":           1,
		"form_id":                strings.ReplaceAll(uuid.New().String(), "-", ""),
		"receipt_without_sku":    nil,
		"vip_money":              "",
		"vip_buy_user_ticket_id": "",
		"coupons_money":          "",
		"coupons_id":             "",
	}
	conf := config.GetDingDong()
	if conf.PayType > 0 {
		paymentOrder["pay_type"] = conf.PayType
	}

	for k, v := range checkOrderMap {
		if v == nil {
			continue
		}
		paymentOrder[k] = v
	}

	packages := map[string]interface{}{
		"reserved_time_start":     reserveTime.StartTimestamp,
		"reserved_time_end":       reserveTime.EndTimestamp,
		"eta_trace_id":            "",
		"soon_arrival":            "",
		"first_selected_big_time": 1,
		"real_match_supply_order": false,
		"time_biz_type":           0,
	}
	fields := []string{"products", "only_today_products", "only_tomorrow_products", "can_used_balance_money", "can_used_point_money", "can_used_point_num", "cart_count", "front_package_bg_color", "front_package_stock_color", "front_package_text", "front_package_type", "goods_real_money", "instant_rebate_money", "is_presale", "is_share_station", "is_supply_order", "real_match_supply_order", "time_biz_type", "package_id", "package_type", "total_count", "total_money", "total_origin_money", "total_rebate_money", "used_balance_money", "used_point_money", "used_point_num"}
	for key, val := range cartMap {
		if val == nil {
			continue
		}
		if !textual.InArray(key, fields) {
			continue
		}
		if key == "products" {
			productFields := []string{"id", "parent_id", "count", "cart_id", "price", "product_type", "is_booking", "product_name", "small_image", "sale_batches", "order_sort", "sizes"}
			items := val.([]map[string]interface{})
			products := make([]map[string]interface{}, 0)
			for _, item := range items {
				products = append(products, filterFields(item, productFields))
			}
			packages[key] = products
			continue
		}
		packages[key] = val
	}

	packageOrder := map[string]interface{}{
		"payment_order": paymentOrder,
		"packages":      []interface{}{packages},
	}
	packageOrderJson := json.MustEncodeToString(packageOrder)

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["package_order"] = packageOrderJson
	params["showMsg"] = "false"
	params["showData"] = "true"
	params["ab_config"] = `{"key_onion":"C"}`
	// log.Printf("AddNewOrder params => %#v", params)
	form, err := session.Sign(params)
	if err != nil {
		return errs.Wrap(code.SignFailed, err)
	}

	result := dto.Result{}
	resp, err := session.Client().R().
		SetHeaders(headers).
		SetFormData(form).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodPost, api)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		if result.Code == 5004 {
			return errs.Wrap(code.ResponseError, errs.New(code.ReserveTimeIsDisabled))
		}
		return errs.WithMessage(code.ResponseError, "提交订单失败 => "+json.MustEncodeToString(result))
	}
	log.Println("恭喜你，已成功下单 =>", resp.String())
	return nil
}
