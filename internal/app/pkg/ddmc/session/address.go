package session

import (
	"net/http"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto/address"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
)

func GetAddress() ([]address.Item, error) {
	url := "https://sunquan.api.ddxq.mobi/api/v1/user/address/"

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
	query, err := Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := address.Result{}
	_, err = Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		SetRetryCount(5).
		Send(http.MethodGet, url)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	// log.Println(resp.String())
	if !result.Success {
		return nil, errs.WithMessage(code.InvalidResponse, result.Message)
	}
	if len(result.Data.Valid) == 0 {
		return nil, errs.New(code.NoValidAddress)
	}

	// res := make(map[string]address.Item)
	// for _, v := range result.Data.Valid {
	// 	str := fmt.Sprintf("%s %s %s", v.UserName, v.Location.Address, v.AddrDetail)
	// 	res[str] = v
	// }
	return result.Data.Valid, nil
}
