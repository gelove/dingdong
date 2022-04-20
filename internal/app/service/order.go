package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"dingdong/internal/app/dto"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"

	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/pkg/json"
)

func GetCheckOrder(cartMap map[string]interface{}, reserveTime *reserve_time.GoTimes) (map[string]interface{}, error) {
	url := "https://maicai.api.ddxq.mobi/order/checkOrder"
	conf := config.Get()

	packagesInfo := make(map[string]interface{})
	for k, v := range cartMap {
		packagesInfo[k] = v
	}
	packagesInfo["reserved_time"] = map[string]interface{}{
		"reserved_time_start": reserveTime.StartTimestamp,
		"reserved_time_end":   reserveTime.EndTimestamp,
	}
	packagesJson := json.MustEncodeToString([]interface{}{packagesInfo})

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["packages"] = packagesJson
	params["address_id"] = conf.Params["address_id"]
	params["user_ticket_id"] = "default"
	params["freight_ticket_id"] = "default"
	params["is_use_point"] = "0"
	params["is_use_balance"] = "0"
	params["is_buy_vip"] = "0"
	params["coupons_id"] = ""
	params["is_buy_coupons"] = "0"
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
		Send(http.MethodPost, url)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return nil, errs.WithMessage(code.InvalidResponse, result.Msg)
	}
	body, err := resp.ToBytes()
	if err != nil {
		return nil, errs.Wrap(code.ParseFailed, err)
	}
	order := json.Get(body, "data", "order")
	// log.Println("data.order", order.ToString())

	res := map[string]interface{}{
		"price":                  order.Get("total_money").ToString(),
		"freight_money":          order.Get("freight_money").ToString(),                                  // 运费
		"freight_discount_money": order.Get("freight_discount_money").ToString(),                         // 运费折扣
		"order_freight":          order.Get("freights", "0", "freight", "freight_real_money").ToString(), // 订单真实运费
		"user_ticket_id":         order.Get("default_coupon", "_id").ToString(),
	}
	log.Println("订单总金额 =>", res["price"])
	return res, nil
}

func AddNewOrder(cartMap map[string]interface{}, reserveTime *reserve_time.GoTimes, checkOrderMap map[string]interface{}) error {
	url := "https://maicai.api.ddxq.mobi/order/addNewOrder"
	conf := config.Get()

	paymentOrder := map[string]interface{}{
		"reserved_time_start":    reserveTime.StartTimestamp,
		"reserved_time_end":      reserveTime.EndTimestamp,
		"parent_order_sign":      cartMap["parent_order_sign"],
		"address_id":             conf.Params["address_id"],
		"pay_type":               6,
		"product_type":           1,
		"form_id":                strings.ReplaceAll(uuid.New().String(), "-", ""),
		"receipt_without_sku":    nil,
		"vip_money":              "",
		"vip_buy_user_ticket_id": "",
		"coupons_money":          "",
		"coupons_id":             "",
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
		"first_selected_big_time": 0,
		"receipt_without_sku":     0,
	}
	for k, v := range cartMap {
		if v == nil {
			continue
		}
		packages[k] = v
	}

	packageOrder := map[string]interface{}{
		"payment_order": paymentOrder,
		"packages":      []interface{}{packages},
	}
	packageOrderJson := json.MustEncodeToString(packageOrder)

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["package_order"] = packageOrderJson
	params["showData"] = "true"
	params["showMsg"] = "false"
	params["ab_config"] = `{"key_onion":"C"}`
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
		Send(http.MethodPost, url)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return errs.WithMessage(code.InvalidResponse, result.Msg)
	}
	log.Println("恭喜你，已成功下单 =>", resp.String())
	return nil
}
