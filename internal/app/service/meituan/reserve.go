package meituan

import (
	"fmt"
	"net/http"

	"github.com/imroc/req/v3"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

func GetMultiReserveTime() (*dto.MeiTuanReserveResult, error) {
	station := config.GetMeiTuan().Station
	api := fmt.Sprintf("https://mall.meituan.com/api/c/poi/%s/notice?uuid=1805465c9dcc8-93f09a12558fe0-0-5a900-1805465c9dc54&xuuid=1805465c9dcc8-93f09a12558fe0-0-5a900-1805465c9dc54&__reqTraceID=c88c301d-8a04-6c76-358b-c2471ce5f824&platform=ios&utm_medium=wxapp&brand=xiaoxiangmaicai&tenantId=1&utm_term=5.33.1&msgOpSource=2&poi=%s&stockPois=%s&ci=1&bizId=2&openId=oV_5G4wXrnpDWzzPA2OpxkkVlZrY&address_id=1950000008&sysName=iOS&sysVerion=15.4&app_tag=union&uci=1&userid=126711747", station, station, station)
	result := new(dto.MeiTuanReserveResult)
	errMsg := new(dto.ErrorMessage)
	resp, err := req.C().R().
		SetResult(result).
		SetError(errMsg).
		SetRetryCount(3).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !resp.IsSuccess() {
		return nil, errs.WithMessage(code.ResponseError, "获取美团运力失败 => "+resp.String())
	}
	if result.Code != 0 {
		return nil, errs.WithMessage(code.ResponseError, "获取美团运力失败 => "+json.MustEncodeToString(result))
	}
	return result, nil
}
