package service

import (
	"net/http"

	"dingdong/internal/app/dto/flow_detail"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

func GetHomeFlowDetail() ([]flow_detail.Item, error) {
	api := "https://maicai.api.ddxq.mobi/homeApi/homeFlowDetail"

	headers := session.GetHeaders()
	headers["accept-language"] = "en-us"
	params := session.GetParams(headers)
	params["tab_type"] = "1"
	params["page"] = "1"
	query, err := session.Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := flow_detail.Result{}
	_, err = session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		SetRetryCount(10).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return nil, errs.WithMessage(code.ResponseError, "获取首页失败 => "+json.MustEncodeToString(result))
	}
	return result.Data.List, nil
}
