package session

import (
	"log"
	"net/http"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto/user"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
)

func GetUser() (*user.Info, error) {
	api := "https://sunquan.api.ddxq.mobi/api/v1/user/detail/"

	h := config.Get().Headers
	headers := map[string]string{
		"Host":   "sunquan.api.ddxq.mobi",
		"cookie": h["cookie"],
	}

	params := make(map[string]string)
	params["channel"] = "applet"
	params["api_version"] = "9.50.0"
	params["app_version"] = "2.83.1"
	params["app_client_id"] = "4"
	params["uid"] = ""
	params["applet_source"] = ""
	params["h5_source"] = ""
	params["sharer_uid"] = ""
	params["s_id"] = ""
	params["openid"] = ""
	params["device_token"] = ""
	params["source_type"] = "5"
	query, err := Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	var result user.Result
	_, err = Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		SetRetryCount(5).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return nil, errs.WithMessage(code.InvalidResponse, result.Message)
	}

	log.Printf("获取用户信息成功, id: %s, name: %s", result.Data.UserInfo.ID, result.Data.UserInfo.Name)
	return &result.Data.UserInfo, nil
}
