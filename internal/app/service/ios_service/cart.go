package ios_service

import (
	"log"
	"net/http"
	"time"

	"dingdong/internal/app/dto"
	"dingdong/internal/app/pkg/ddmc/ios_session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/pkg/json"
)

func AllCheck() error {
	api := "https://maicai.api.ddxq.mobi/cart/allCheck"

	now := time.Now().Unix()
	headers := ios_session.TakeHeaders(now, []string{"accept", "accept-language", "content-type", "cookie", "ddmc-api-version", "ddmc-app-client-id", "ddmc-build-version", "ddmc-channel", "ddmc-city-number", "ddmc-country-code", "ddmc-device-id", "ddmc-device-model", "ddmc-device-name", "ddmc-device-token", "ddmc-idfa", "ddmc-ip", "ddmc-language-code", "ddmc-latitude", "ddmc-locale-identifier", "ddmc-longitude", "ddmc-os-version", "ddmc-station-id", "ddmc-uid", "user-agent", "im_secret"})
	params := ios_session.TakeParams(now, []string{"ab_config", "api_version", "app_client_id", "app_type", "buildVersion", "channel", "city_number", "countryCode", "device_id", "device_model", "device_name", "device_token", "idfa", "ip", "languageCode", "latitude", "localeIdentifier", "longitude", "os_version", "seqid", "station_id", "uid"})
	params["is_filter"] = "0"
	params["is_check"] = "1"
	params["is_load"] = "1"

	result := dto.Result{}
	errMsg := new(dto.ErrorMessage)
	resp, err := ios_session.Client().R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&result).
		SetError(errMsg).
		SetRetryCount(3).
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

func MockCartMap() map[string]any {
	data := make(map[string]any)
	// cartStr := `{"products":[{"type":1,"id":"612cc0982c34fab505117d4e","price":"828.00","count":1,"description":"","sizes":[],"cart_id":"612cc0982c34fab505117d4e","parent_id":"","parent_batch_type":-1,"category_path":"","manage_category_path":"411,412,413","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"洋河蓝色经典梦之蓝M6+52度白酒 550ml/瓶","product_type":0,"small_image":"https://ddfs-public.ddimg.mobi/img/blind/product-management/202108/1242efbb2a37470aa081683513fb3677.jpg?width=800&height=800","total_price":"828.00","origin_price":"828.00","total_origin_price":"828.00","no_supplementary_price":"828.00","no_supplementary_total_price":"828.00","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"瓶","net_weight":"550","net_weight_unit":"ml","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":1,"is_presale":0}]}`
	cartStr := `{"products":[{"type":1,"id":"6262964b929609dd27fc214d","price":"8.80","count":1,"description":"","sizes":[],"cart_id":"6262964b929609dd27fc214d","parent_id":"","parent_batch_type":-1,"category_path":"","manage_category_path":"1916,1936,1937","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"叮咚保供定制皇后吐司 110g/袋","product_type":0,"small_image":"https://imgnew.ddimg.mobi/product/b42d4fa545584bada451a43f1fedc0e3.jpg?width=800&height=800","total_price":"8.80","origin_price":"8.80","total_origin_price":"8.80","no_supplementary_price":"8.80","no_supplementary_total_price":"8.80","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"袋","net_weight":"110","net_weight_unit":"g","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":1,"is_presale":0}]}`
	json.MustDecodeFromString(cartStr, &data)
	productList := data["products"].([]any)
	products := make([]any, 0, len(productList))
	for _, v := range productList {
		product := v.(map[string]any)
		product["total_money"] = product["total_price"]
		product["total_origin_money"] = product["total_origin_price"]
		products = append(products, product)
	}
	cartMap := make(map[string]any)
	cartMap["products"] = products
	return cartMap
}

func GetCart() (map[string]any, error) {
	api := "https://maicai.api.ddxq.mobi/cart/index"

	now := time.Now().Unix()
	headers := ios_session.TakeHeaders(now, nil)
	params := ios_session.TakeParams(now, nil)

	result := new(dto.Result)
	errMsg := new(dto.ErrorMessage)
	resp, err := ios_session.Client().R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&result).
		SetError(errMsg).
		SetRetryCount(3).
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

	products := make([]any, 0, len(productList))
	for _, v := range productList {
		product, ok := v.(map[string]any)
		if !ok {
			continue
		}
		product["total_money"] = product["total_price"]
		product["total_origin_money"] = product["total_origin_price"]
		product["is_coupon_gift"] = product["is_gift"]
		products = append(products, product)
	}
	for k, v := range products {
		val := v.(map[string]any)
		log.Printf("[%v] %s 数量: %v 总价: %s", k, val["product_name"], val["count"], val["total_price"])
	}
	res["products"] = products

	res["real_match_supply_order"] = res["is_supply_order"]
	if v, ok := data["parent_order_info"].(map[string]any); ok {
		res["parent_order_sign"] = v["parent_order_sign"]
	}
	return res, nil
}
