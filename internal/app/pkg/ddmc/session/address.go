package session

import (
	"net/http"

	"dingdong/internal/app/dto/address"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

func GetAddress() ([]address.Item, error) {
	api := "https://sunquan.api.ddxq.mobi/api/v1/user/address/"

	headers := GetUserHeader()
	params := GetUserParams(headers)
	params["source_type"] = "5"

	result := address.Result{}
	_, err := Client().R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&result).
		SetRetryCount(5).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	// log.Println(resp.String())
	if !result.Success {
		return nil, errs.WithMessage(code.ResponseError, "获取地址失败 => "+json.MustEncodeToString(result))
	}
	if len(result.Data.Valid) == 0 {
		return nil, errs.New(code.NoValidAddress)
	}

	return result.Data.Valid, nil
}
