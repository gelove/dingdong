package service

import (
	"net/http"
	"time"

	"dingdong/internal/app/dto/flow_detail"
	"dingdong/internal/app/pkg/ddmc/ios_session"
	"dingdong/internal/app/pkg/errs"
)

func GetHomeFlowDetail() ([]flow_detail.Item, error) {
	api := "https://maicai.api.ddxq.mobi/homeApi/homeFlowDetail"

	now := time.Now().Unix()
	headers := ios_session.TakeHeaders(now, nil)
	headers["accept-language"] = "en-us"
	params := ios_session.TakeParams(now, nil)
	params["tab_type"] = "1"
	params["page"] = "1"
	query, err := ios_session.Sign(params)
	if err != nil {
		return nil, err
	}

	result := flow_detail.Result{}
	resp, err := ios_session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		SetRetryCount(10).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	if !result.Success {
		return nil, errs.Wrap(errs.GetFlowDetailFailed, resp.String())
	}
	return result.Data.List, nil
}
