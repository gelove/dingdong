package applet_service

import (
	"log"
	"net/http"

	"dingdong/internal/app/dto"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
)

func AllCheck() error {
	api := "https://maicai.api.ddxq.mobi/cart/allCheck"

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["is_check"] = "1"
	params["is_load"] = "1"
	params["ab_config"] = `{"key_onion":"D","key_cart_discount_price":"C"}`
	query, err := session.Sign(params)
	if err != nil {
		return err
	}

	result := dto.Result{}
	resp, err := session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodGet, api)
	if err != nil {
		return errs.WithStack(err)
	}
	if !resp.IsSuccess() {
		return errs.Wrap(errs.CheckAllFailed, resp.String())
	}
	if !result.Success {
		return errs.Wrap(errs.CheckAllFailed, resp.String())
	}
	return nil
}

func GetCart() (map[string]any, error) {
	api := "https://maicai.api.ddxq.mobi/cart/index"

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["is_load"] = "1"
	params["ab_config"] = `{"key_onion":"D","key_cart_discount_price":"C"}`
	query, err := session.Sign(params)
	if err != nil {
		return nil, err
	}

	result := dto.Result{}
	resp, err := session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	if !resp.IsSuccess() {
		return nil, errs.Wrap(errs.GetCartFailed, resp.String())
	}
	if !result.Success {
		return nil, errs.Wrap(errs.GetCartFailed, resp.String())
	}
	data, ok := result.Data.(map[string]any)
	if !ok {
		return nil, errs.WithStack(errs.GetCartFailed)
	}
	// 有效可购的商品
	list, ok := data["new_order_product_list"].([]any)
	if !ok || len(list) == 0 {
		return nil, errs.WithStack(errs.NoValidProduct)
	}

	first, ok := list[0].(map[string]any)
	if !ok {
		return nil, errs.WithStack(errs.GetCartFailed)
	}
	// coupon_rebate_money 优惠券返现金额 我的请求中没有这个字段
	res := make(map[string]any)
	for k, v := range first {
		if k == "products" || v == nil {
			continue
		}
		res[k] = v
	}
	productList, ok := first["products"].([]any)
	if !ok {
		return nil, errs.WithStack(errs.GetCartFailed)
	}
	products := make([]map[string]any, 0, len(productList))
	for _, v := range productList {
		product, ok := v.(map[string]any)
		if !ok {
			return nil, errs.WithStack(errs.GetCartFailed)
		}
		product["total_money"] = product["total_price"]
		product["total_origin_money"] = product["total_origin_price"]
		products = append(products, product)
	}
	for k, v := range products {
		log.Printf("[%v] %s 数量: %v 总价: %s", k, v["product_name"], v["count"], v["total_price"])
	}
	res["products"] = products

	if v, ok := data["parent_order_info"].(map[string]any); ok {
		res["parent_order_sign"] = v["parent_order_sign"]
	}
	return res, nil
}
