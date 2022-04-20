package service

import (
	"fmt"
	"net/http"

	"dingdong/internal/app/dto/address"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
)

func GetAddress() (map[string]address.Item, error) {
	url := "https://sunquan.api.ddxq.mobi/api/v1/user/address/"

	headers := session.GetHeaders()
	params := session.GetParams(headers)
	query, err := session.Sign(params)
	if err != nil {
		return nil, errs.Wrap(code.SignFailed, err)
	}

	result := address.Result{}
	_, err = session.Client().R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetResult(&result).
		// SetRetryCount(50).
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

	res := make(map[string]address.Item)
	for _, v := range result.Data.Valid {
		str := fmt.Sprintf("%s %s %s", v.UserName, v.Location.Address, v.AddrDetail)
		res[str] = v
	}
	return res, nil
}
