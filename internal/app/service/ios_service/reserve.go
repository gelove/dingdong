package ios_service

import (
	"net/http"
	"time"

	"dingdong/internal/app/dto"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/pkg/ddmc/ios_session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/pkg/json"
)

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

func GetMultiReserveTime(cartMap map[string]any) (*reserve_time.GoTimes, error) {
	api := "https://maicai.api.ddxq.mobi/order/getMultiReserveTime"

	products := cartMap["products"].([]any)

	list := make([]map[string]any, 0, len(products))
	fields := []string{"sale_batches", "is_coupon_gift", "id", "price", "is_booking", "count", "small_image", "type", "origin_price", "product_type", "product_name"}
	for _, product := range products {
		val, ok := product.(map[string]any)
		if !ok {
			continue
		}
		item := make(map[string]any)
		for _, v := range fields {
			item[v] = val[v]
		}
		if item["is_coupon_gift"] == nil {
			item["is_coupon_gift"] = 0
		}
		list = append(list, item)
	}
	listStr := json.MustEncodeFast(list)
	productsJson := json.MustEncodeFast([]string{listStr})
	// log.Printf("%s", productsJson)
	now := time.Now().Unix()
	// now = 1651584614 // mock 测试用
	headers := ios_session.TakeHeaders(now, []string{"accept", "accept-encoding", "accept-language", "content-type", "cookie", "ddmc-api-version", "ddmc-app-client-id", "ddmc-build-version", "ddmc-channel", "ddmc-city-number", "ddmc-country-code", "ddmc-device-id", "ddmc-device-model", "ddmc-device-name", "ddmc-device-token", "ddmc-idfa", "ddmc-ip", "ddmc-language-code", "ddmc-latitude", "ddmc-locale-identifier", "ddmc-longitude", "ddmc-os-version", "ddmc-station-id", "ddmc-uid", "user-agent", "x-tingyun", "x-tingyun-id"})
	params := ios_session.TakeParams(now, []string{"api_version", "app_client_id", "app_type", "buildVersion", "channel", "city_number", "countryCode", "device_id", "device_model", "device_name", "device_token", "idfa", "ip", "languageCode", "latitude", "localeIdentifier", "longitude", "os_version", "station_id", "uid", "address_id", "seqid"})
	// params["seqid"] = "3358795165" // mock 测试用
	params["ab_config"] = `{"ETA_time_default_selection":"C1.1"}`
	params["address_id"] = ios_session.Address().Id
	params["products"] = productsJson

	signed, err := ios_session.Sign(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = signed["sign"]
	headers["sign"] = signed["sign"]
	headers["nars"] = signed["nars"]
	headers["sesi"] = signed["sesi"]

	result := new(reserve_time.Result)
	errMsg := new(dto.ErrorMessage)
	resp, err := ios_session.Client().R().
		SetHeaders(headers).
		SetFormData(params).
		SetResult(result).
		SetError(errMsg).
		// SetRetryCount(3).
		Send(http.MethodPost, api)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	if !resp.IsSuccess() {
		return nil, errs.Wrap(errs.GetReserveTimeFailed, "[叮咚] => "+resp.String())
	}
	if !result.Success {
		return nil, errs.Wrap(errs.GetReserveTimeFailed, "[叮咚] => "+resp.String())
	}
	if len(result.Data) == 0 || len(result.Data[0].Times) == 0 || len(result.Data[0].Times[0].Times) == 0 {
		return nil, errs.WithStack(errs.NoReserveTime)
	}

	times := result.Data[0].Times[0].Times
	validTimes := filterValidTimes(times)
	if len(validTimes) == 0 {
		return nil, errs.WithStack(errs.NoReserveTimeAndRetry)
	}
	validTime := validTimes[0]
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
