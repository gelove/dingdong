package service

import (
	"log"
	"net/http"
	"time"

	"dingdong/internal/app/dto"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

func MockCartMap() map[string]interface{} {
	first := make(map[string]interface{})
	// cartStr := `{"products":[{"type":1,"id":"612cc0982c34fab505117d4e","price":"828.00","count":1,"description":"","sizes":[],"cart_id":"612cc0982c34fab505117d4e","parent_id":"","parent_batch_type":-1,"category_path":"","manage_category_path":"411,412,413","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"洋河蓝色经典梦之蓝M6+52度白酒 550ml/瓶","product_type":0,"small_image":"https://ddfs-public.ddimg.mobi/img/blind/product-management/202108/1242efbb2a37470aa081683513fb3677.jpg?width=800&height=800","total_price":"828.00","origin_price":"828.00","total_origin_price":"828.00","no_supplementary_price":"828.00","no_supplementary_total_price":"828.00","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"瓶","net_weight":"550","net_weight_unit":"ml","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":1,"is_presale":0}]}`
	cartStr := `{"products":[{"price":"1199","name":"泸州老窖 国窖1573浓香型52度白酒 500ml/瓶","spec":"传承古老匠心技艺，佳节礼赠宴请（非质量原因，一经售出，不退不换）","sizes":[],"status":1,"type":0,"activity":[],"oid":48010,"id":"612cc085095fb38bb1a61a0a","origin_price":"1199","vip_price":"","stock_number":9000,"marketing_tags":[],"is_promotion":0,"buy_limit":0,"share_marketing_tags":[],"is_booking_new":null,"today_stockout_new":null,"product_name":"泸州老窖 国窖1573浓香型52度白酒 500ml/瓶","small_image":"https://ddfs-public.ddimg.mobi/img/blind/product-management/202108/e891f048fb7145949504a2e543f0e52a.jpg?width=800&amp;height=800","category_id":"","total_sales":12998,"month_sales":0,"station_stock":-1,"mark_discount":0,"mark_new":0,"mark_self":0,"category_path":"","stockout_reserved":false,"sale_point_msg":[],"is_presale":0,"presale_delivery_date_display":"","is_gift":0,"is_bulk":0,"net_weight":"500","net_weight_unit":"ml","is_onion":0,"is_invoice":1,"badge_img":"","badge_position":1,"is_vod":false,"decision_information":["保真正品","高端礼赠","节日宴请"],"user_bought":0,"sub_list":[],"desc_tags":[{"type":123,"name":"保真正品"},{"type":123,"name":"高端礼赠"},{"type":123,"name":"节日宴请"}],"alg_tags":[],"image_preferential_choice":"","today_stockout":"","is_booking":0,"min_order_quantity":null,"sale_unit":"瓶","temperature_layer":"","storage_value_id":0,"attribute_tags":[],"algo_id":"8c2e391001c1dc968e0f7bd746bacdca","recommend_cate":1,"recommend_reason":"","share_text":"新用户领108元红包","scene_id":-1,"scene_tag":"","recommended_reason":["保真正品","高端礼赠","节日宴请"],"wx_applet_code_path":"pages/productDetail/productDetail","wx_applet_url":"/pages/productDetail/productDetail?product_id=612cc085095fb38bb1a61a0a","web_url":"/pages/vipPackage/vip/vip?id=612cc085095fb38bb1a61a0a","feature_tag":0,"feature_tag_url":{"type":"0","name":"","url_map":{}}}]}`
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

func MockMultiReserveTime() *reserve_time.GoTimes {
	reserveTime := &reserve_time.GoTimes{}
	halfPastTwoPM := date.TodayUnix(14, 30, 0)
	now := time.Now().Unix()
	if now >= date.TodayUnix(0, 0, 0) && now <= date.TodayUnix(6, 20, 0) {
		reserveTime.StartTimestamp = date.TodayUnix(6, 30, 0)
		reserveTime.EndTimestamp = halfPastTwoPM
		return reserveTime
	}
	if now >= date.TodayUnix(8, 20, 0) && now <= date.TodayUnix(8, 50, 0) {
		reserveTime.StartTimestamp = halfPastTwoPM
		reserveTime.EndTimestamp = date.TodayUnix(22, 30, 0)
		return reserveTime
	}
	reserveTime.StartTimestamp = (now/60 + 5) * 60 // 叮咚是在当前时间直接加5分钟
	reserveTime.EndTimestamp = date.TodayUnix(22, 30, 0)
	return reserveTime
}

func GetMultiReserveTime(cartMap map[string]interface{}) (*reserve_time.GoTimes, error) {
	api := "https://maicai.api.ddxq.mobi/order/getMultiReserveTime"

	products := cartMap["products"].([]map[string]interface{})
	productsList := [][]map[string]interface{}{products}
	productsJson := json.MustEncodeToString(productsList)

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	params["address_id"] = session.Address().Id
	params["group_config_id"] = ""
	params["isBridge"] = "false"
	params["products"] = productsJson

	form, err := session.Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := new(reserve_time.Result)
	errMsg := new(dto.ErrorMessage)
	resp, err := session.Client().R().
		SetHeaders(headers).
		SetFormData(form).
		SetResult(result).
		SetError(errMsg).
		// SetRetryCount(50).
		Send(http.MethodPost, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !resp.IsSuccess() {
		return nil, errs.WithMessage(code.ResponseError, "获取叮咚运力失败 => "+resp.String())
	}
	if !result.Success {
		return nil, errs.WithMessage(code.ResponseError, "获取叮咚运力失败 => "+json.MustEncodeToString(result))
	}
	if len(result.Data) == 0 || len(result.Data[0].Times) == 0 || len(result.Data[0].Times[0].Times) == 0 {
		return nil, errs.New(code.NoReserveTime)
	}

	times := result.Data[0].Times[0].Times
	validTimes := filterValidTimes(times)
	if len(validTimes) == 0 {
		return nil, errs.New(code.NoReserveTimeAndRetry)
	}
	validTime := validTimes[0]
	log.Printf("发现可用运力[%s-%s](%d-%d), 请尽快下单!", validTime.StartTime, validTime.EndTime, validTime.StartTimestamp, validTime.EndTimestamp)
	return validTime, nil
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
