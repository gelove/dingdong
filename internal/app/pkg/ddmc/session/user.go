package session

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto/user_dto"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
)

func GetUserHeader() map[string]string {
	h := config.GetDingDong().Headers
	headers := map[string]string{
		"Host":            "sunquan.api.ddxq.mobi",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Origin":          "https://wx.m.ddxq.mobi",
		"Cookie":          h["cookie"],
		"Accept":          "application/json, text/plain, */*",
		"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x1800123f) NetType/WIFI Language/zh_CN",
		"Referer":         "https://wx.m.ddxq.mobi/",
		"Accept-Language": "zh-CN,zh-Hans;q=0.9",
	}
	return headers
}

func GetUserParams(headers map[string]string) map[string]string {
	params := map[string]string{
		// "uid" : "",
		// "longitude" : "",
		// "latitude" : "",
		// "station_id" : "",
		// "city_number" : "",
		// "device_token" : "",
		"api_version":   "9.50.2",
		"app_version":   "2.85.3",
		"applet_source": "",
		"app_client_id": "3",
		"h5_source":     "",
		"wx":            "1",
		"sharer_uid":    "",
		"s_id":          strings.TrimLeft(strings.Split(headers["Cookie"], ";")[0], "DDXQSESSID="),
		"openid":        "",
		"time":          strconv.Itoa(int(time.Now().Unix())),
	}
	return params
}

func GetUser() (*user_dto.Info, error) {
	api := "https://sunquan.api.ddxq.mobi/api/v1/user/detail/"

	headers := GetUserHeader()
	params := GetUserParams(headers)
	params["source_type"] = "mine_page"

	var result user_dto.Result
	_, err := Client().R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&result).
		SetRetryCount(5).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.Wrap(code.RequestFailed, err)
	}
	if !result.Success {
		return nil, errs.WithMessage(code.ResponseError, "获取用户信息失败 => "+json.MustEncodeToString(result))
	}

	log.Printf("获取用户信息成功, id: %s, name: %s", result.Data.UserInfo.ID, result.Data.UserInfo.Name)
	return &result.Data.UserInfo, nil
}
