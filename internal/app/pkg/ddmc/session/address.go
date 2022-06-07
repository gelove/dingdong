package session

import (
	"net/http"

	"dingdong/internal/app/dto/address"
	"dingdong/internal/app/pkg/errs"
)

func GetAddress() ([]*address.Item, error) {
	api := "https://sunquan.api.ddxq.mobi/api/v1/user/address/"

	headers := GetUserHeader()
	params := GetUserParams(headers)
	params["source_type"] = "5"

	result := address.Result{}
	resp, err := Client().R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&result).
		SetRetryCount(5).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	if !result.Success {
		return nil, errs.Wrap(errs.GetAddressFailed, resp.String())
	}
	if len(result.Data.Valid) == 0 {
		return nil, errs.WithStack(errs.NoValidAddress)
	}

	return result.Data.Valid, nil
}
