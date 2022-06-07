package ios_session

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imroc/req/v3"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto"
	"dingdong/internal/app/dto/address"
	"dingdong/internal/app/dto/session_dto"
	"dingdong/internal/app/pkg/errs"
	"dingdong/pkg/crypto"
	"dingdong/pkg/js"
	"dingdong/pkg/json"
	"dingdong/pkg/textual"
	"dingdong/pkg/uri"
)

var (
	once sync.Once
	s    *session
)

type session struct {
	ImSecret string
	Headers  map[string]string
	Params   map[string]string
	Client   *req.Client
	Address  *address.Item
}

func Initialize(dir string) {
	once.Do(func() {
		// client := req.DevMode().EnableForceHTTP1()
		// client := req.C().EnableForceHTTP1()
		client := req.C().EnableForceHTTP1().
			SetCommonRetryCondition(retryCondition).
			SetCommonRetryInterval(retryInterval).
			SetCommonRetryHook(retryHook).
			SetCommonRetryBackoffInterval(1*time.Millisecond, 10*time.Millisecond)

		s = &session{
			Client: client,
		}

		setImSecret()
		chooseSession(dir, false)
		chooseAddr("")
	})
}

func InitializeMock(dir string, addr string) {
	once.Do(func() {
		client := req.DevMode().EnableForceHTTP1()
		// client := req.C().EnableForceHTTP1()

		s = &session{
			Client: client,
		}

		setImSecret()
		chooseSession(dir, true)
		chooseAddr(addr)
	})
}

func retryCondition(resp *req.Response, err error) bool {
	if err != nil || resp.StatusCode != http.StatusOK {
		return true
	}
	body, err := resp.ToBytes()
	if err != nil {
		return true
	}
	success := json.Get(body, "success").ToBool()
	return !success
}

func retryInterval(resp *req.Response, attempt int) time.Duration {
	duration := 150 + rand.Intn(50)
	return time.Duration(duration) * time.Millisecond
}

func retryHook(resp *req.Response, err error) {
	if err != nil {
		log.Printf("Request error => %v", err)
	}
	// r := resp.Request.RawRequest
	// log.Println("Retry request =>", r.Method, r.URL.Path)
}

func Client() *req.Client {
	return s.Client
}

func Address() *address.Item {
	return s.Address
}

func setImSecret() {
	conf := config.GetDingDong()
	s.ImSecret = conf.ImSecret
}

func chooseSession(dir string, isTest bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	options := make([]string, 0, len(files))
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".chlsj") {
			options = append(options, name)
		}
	}

	if len(options) == 0 {
		panic(errs.SelectSessionFailed)
	}

	var file string
	if isTest || len(options) == 1 {
		file = options[0]
	} else {
		sv := &survey.Select{
			Message: "请选择session文件",
			Options: options,
		}
		err = survey.AskOne(sv, &file)
		if err != nil {
			panic(err)
		}
	}

	bs, err := ioutil.ReadFile(dir + "/" + file)
	if err != nil {
		panic(err)
	}

	list := make([]session_dto.Data, 0, 1<<1)
	json.MustDecode(bs, &list)
	if len(list) == 0 {
		panic("session list is empty")
	}

	SetHeader(list[0])
	SetParams(list[0])
}

func chooseAddr(match string) {
	addrList, err := GetAddress()
	if err != nil {
		panic(err)
	}
	if len(addrList) == 1 {
		s.Address = addrList[0]
		log.Println(json.MustEncodePrettyString(s.Address))
		log.Printf("默认收货地址: %s %s %s", s.Address.Location.Address, s.Address.Location.Name, s.Address.AddrDetail)
		return
	}

	var index int

	options := make([]string, 0, len(addrList))
	for i, v := range addrList {
		option := fmt.Sprintf("%s %s %s %s", v.Location.Address, v.Location.Name, v.AddrDetail, v.StationName)
		options = append(options, option)
		if match != "" && strings.Contains(option, match) {
			index = i
		}
	}

	if match == "" {
		var addr string
		sv := &survey.Select{
			Message: "请选择收货地址",
			Options: options,
		}
		if err := survey.AskOne(sv, &addr); err != nil {
			panic(err)
		}
		index = textual.IndexOf(addr, options)
	}

	s.Address = addrList[index]
	s.Headers["ddmc-city-number"] = s.Address.CityNumber
	s.Headers["ddmc-station-id"] = s.Address.StationId
	s.Headers["ddmc-longitude"] = strconv.FormatFloat(s.Address.Location.Location[0], 'f', -1, 64)
	s.Headers["ddmc-latitude"] = strconv.FormatFloat(s.Address.Location.Location[1], 'f', -1, 64)
	s.Params["city_number"] = s.Headers["ddmc-city-number"]
	s.Params["station_id"] = s.Headers["ddmc-station-id"]
	s.Params["longitude"] = s.Headers["ddmc-longitude"]
	s.Params["latitude"] = s.Headers["ddmc-latitude"]
	// log.Printf("Address => %#v", s.Address)
	log.Printf("已选择收货地址: %s %s %s %s", s.Address.Location.Address, s.Address.Location.Name, s.Address.AddrDetail, s.Address.StationName)
	return
}

func SetHeader(data session_dto.Data) {
	headers := make(map[string]string)
	for _, v := range data.Request.Header.Headers {
		key := strings.ToLower(v.Name)
		if strings.HasPrefix(key, ":") {
			continue
		}
		if key == "nars" || key == "sesi" || key == "sign" || key == "accept-encoding" {
			continue
		}
		headers[key] = v.Value
	}
	headers["im_secret"] = s.ImSecret
	s.Headers = headers
}

func SetParams(data session_dto.Data) {
	query := data.Query
	if strings.Contains(data.Query, "?") {
		query = strings.Split(data.Query, "?")[1]
	}
	queryStr, err := url.QueryUnescape(query)
	if err != nil {
		panic(err)
	}
	// log.Println("queryStr =>", queryStr)
	queryList := strings.Split(queryStr, "&")
	res := make(map[string]string)
	for _, v := range queryList {
		key, value, found := strings.Cut(v, "=")
		if !found {
			continue
		}
		// 别改这里, 测试签名时需要存储的旧 seqid
		if key == "sign" {
			continue
		}
		res[key] = value
	}
	s.Params = res
}

func TakeHeaders(timestamp int64, fields []string) map[string]string {
	res := make(map[string]string)
	if len(fields) > 0 {
		for _, v := range fields {
			key := strings.ToLower(v)
			if v, ok := s.Headers[key]; ok {
				res[key] = v
			}
		}
	} else {
		for k, v := range s.Headers {
			res[k] = v
		}
	}
	encrypt := crypto.MD5(fmt.Sprintf("private_key=%s&time=%d", s.ImSecret, timestamp))
	res["time"] = fmt.Sprintf("%d,%s", timestamp, encrypt)
	res["content-type"] = "application/x-www-form-urlencoded"
	return res
}

func TakeParams(timestamp int64, fields []string) map[string]string {
	res := make(map[string]string)
	if len(fields) > 0 {
		for _, v := range fields {
			res[v] = s.Params[v]
		}
	} else {
		for k, v := range s.Params {
			res[k] = v
		}
	}
	res["time"] = strconv.Itoa(int(timestamp))
	return res
}

func Sign(params map[string]string) (map[string]string, error) {
	conf := config.GetDingDong()
	str := json.MustEncodeFast(params)
	res, err := js.Call("js/ios_sign.js", "ios_sign", conf.ImSecret, str)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	signed := make(map[string]string)
	json.MustDecodeFromString(res.String(), &signed)
	return signed, nil
}

func EncodeFormData(params map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range params {
		res[url.QueryEscape(k)] = url.QueryEscape(v)
	}
	return res
}

// EncodeFormDataToString 排序QueryEscape
func EncodeFormDataToString(params map[string]string) string {
	if params == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := params[k]
		keyEscaped := uri.QueryEscape(k)
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		buf.WriteString(uri.QueryEscape(v))
	}
	return buf.String()
}

func GetAddress() ([]*address.Item, error) {
	api := "https://sunquan.api.ddxq.mobi/api/v1/user/address/"

	now := time.Now().Unix()
	headers := TakeHeaders(now, []string{"accept", "accept-encoding", "accept-language", "content-type", "cookie", "ddmc-api-version", "ddmc-app-client-id", "ddmc-build-version", "ddmc-channel", "ddmc-city-number", "ddmc-country-code", "ddmc-device-id", "ddmc-device-model", "ddmc-device-name", "ddmc-device-token", "ddmc-idfa", "ddmc-ip", "ddmc-language-code", "ddmc-latitude", "ddmc-locale-identifier", "ddmc-longitude", "ddmc-os-version", "ddmc-station-id", "ddmc-uid", "time", "user-agent"})
	params := TakeParams(now, nil)
	params["source_type"] = "5"

	result := address.Result{}
	errMsg := new(dto.ErrorMessage)
	resp, err := Client().R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&result).
		SetError(errMsg).
		SetRetryCount(5).
		Send(http.MethodGet, api)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	if !resp.IsSuccess() {
		return nil, errs.Wrap(errs.GetAddressFailed, resp.String())
	}
	if !result.Success {
		return nil, errs.Wrap(errs.GetAddressFailed, resp.String())
	}
	if len(result.Data.Valid) == 0 {
		return nil, errs.WithStack(errs.NoValidAddress)
	}
	return result.Data.Valid, nil
}
