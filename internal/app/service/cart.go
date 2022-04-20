package service

import (
	"log"
	"net/http"

	"dingdong/internal/app/dto"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
)

func AllCheck() error {
	url := "https://maicai.api.ddxq.mobi/cart/allCheck"

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["is_check"] = "1"
	query, err := session.Sign(params)
	if err != nil {
		return errs.Wrap(code.SignFailed, err)
	}

	result := dto.Result{}
	_, err = session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodGet, url)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return errs.WithMessage(code.InvalidResponse, result.Msg)
	}
	log.Println("勾选购物车全选按钮成功")
	return nil
}

func GetCart() (map[string]interface{}, error) {
	url := "https://maicai.api.ddxq.mobi/cart/index"

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["is_load"] = "1"
	params["ab_config"] = `{"key_onion":"D","key_cart_discount_price":"C"}`
	query, err := session.Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := dto.Result{}
	_, err = session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		// SetRetryCount(50).
		Send(http.MethodGet, url)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return nil, errs.WithMessage(code.InvalidResponse, result.Msg)
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, errs.WithMessage(code.AssertFailed, "获取购物车数据失败")
	}
	// 有效可购的商品
	list, ok := data["new_order_product_list"].([]interface{})
	if !ok || len(list) == 0 {
		return nil, errs.New(code.NoValidProduct)
	}

	first, ok := list[0].(map[string]interface{})
	if !ok {
		return nil, errs.WithMessage(code.AssertFailed, "获取购物车产品数据失败")
	}
	// coupon_rebate_money 优惠券返现金额 我的请求中没有这个字段
	res := make(map[string]interface{})
	for k, v := range first {
		if k == "products" || v == nil {
			continue
		}
		res[k] = v
	}
	productList, ok := first["products"].([]interface{})
	if !ok {
		return nil, errs.New(code.AssertFailed)
	}
	products := make([]map[string]interface{}, 0, len(productList))
	for _, v := range productList {
		product, ok := v.(map[string]interface{})
		if !ok {
			return nil, errs.New(code.AssertFailed)
		}
		product["total_money"] = product["total_price"]
		product["total_origin_money"] = product["total_origin_price"]
		products = append(products, product)
	}
	for k, v := range products {
		log.Printf("[%v] %s 数量: %v 总价: %s", k, v["product_name"], v["count"], v["total_price"])
	}
	res["products"] = products

	if v, ok := data["parent_order_info"].(map[string]interface{}); ok {
		res["parent_order_sign"] = v["parent_order_sign"]
	}
	return res, nil
}
