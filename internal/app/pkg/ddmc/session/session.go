package session

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imroc/req/v3"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto/address"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/js"
	"dingdong/pkg/json"
	"dingdong/pkg/textual"
)

var (
	once sync.Once
	s    *session
)

type session struct {
	UserID  string
	JsFile  string // js文件路径
	Client  *req.Client
	Address address.Item
}

func Initialize(jsFile string) {
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
			JsFile: jsFile,
		}

		setUserID()
		chooseAddr()
	})
}

func InitializeMock(jsFile string) {
	once.Do(func() {
		// client := req.DevMode().EnableForceHTTP1()
		client := req.C().EnableForceHTTP1()

		s = &session{
			Client: client,
			JsFile: jsFile,
		}

		conf := config.Get()
		s.UserID = conf.Mock["ddmc-uid"]
		s.Address = address.Item{
			Id:         conf.Mock["address_id"],
			CityNumber: conf.Mock["ddmc-city-number"],
			StationId:  conf.Mock["ddmc-station-id"],
		}
		longitude, _ := strconv.ParseFloat(conf.Mock["ddmc-longitude"], 64)
		latitude, _ := strconv.ParseFloat(conf.Mock["ddmc-latitude"], 64)
		s.Address.Location.Location = []float64{longitude, latitude}
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
		log.Println("Request error =>", err.Error())
	}
	r := resp.Request.RawRequest
	log.Println("Retry request =>", r.Method, r.URL)
}

func Client() *req.Client {
	return s.Client
}

func JsFile() string {
	return s.JsFile
}

func Address() address.Item {
	return s.Address
}

func setUserID() {
	user, err := GetUser()
	if err != nil {
		panic(err)
	}
	s.UserID = user.ID
}

func chooseAddr() {
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

	options := make([]string, 0, len(addrList))
	for _, v := range addrList {
		options = append(options, fmt.Sprintf("%s %s %s", v.Location.Address, v.Location.Name, v.AddrDetail))
	}

	var addr string
	sv := &survey.Select{
		Message: "请选择收货地址",
		Options: options,
	}
	if err := survey.AskOne(sv, &addr); err != nil {
		panic(errs.Wrap(code.SelectAddressFailed, err))
	}

	index := textual.IndexOf(addr, options)
	s.Address = addrList[index]
	log.Printf("已选择收货地址: %s %s %s", s.Address.Location.Address, s.Address.Location.Name, s.Address.AddrDetail)
	return
}

// GetHeaders 抓包后参考项目中的image/headers.jpeg 把信息一行一行copy到下面 没有的key不需要复制
// AddressId string // 收货地址id
func GetHeaders() map[string]string {
	headers := make(map[string]string)
	// headers["accept-encoding"] = "gzip, deflate, br" // 压缩有乱码
	headers["ddmc-api-version"] = "9.50.0"
	headers["ddmc-app-client-id"] = "4"
	headers["ddmc-build-version"] = "2.83.1"
	headers["ddmc-channel"] = "applet"
	headers["ddmc-ip"] = ""
	headers["ddmc-os-version"] = "[object Undefined]"
	headers["ddmc-time"] = strconv.Itoa(int(time.Now().Unix()))
	headers["referer"] = "https://servicewechat.com/wx1e113254eda17715/430/page-frame.html"
	headers["ddmc-city-number"] = s.Address.CityNumber // 城市id
	headers["ddmc-longitude"] = strconv.FormatFloat(s.Address.Location.Location[0], 'f', -1, 64)
	headers["ddmc-latitude"] = strconv.FormatFloat(s.Address.Location.Location[1], 'f', -1, 64)
	headers["ddmc-station-id"] = s.Address.StationId // 发货站点id
	headers["ddmc-uid"] = s.UserID                   // 用户id

	h := config.Get().Headers
	headers["cookie"] = h["cookie"]
	headers["user-agent"] = h["user-agent"]
	headers["ddmc-device-id"] = h["ddmc-device-id"] // 设备id
	return headers
}

// GetParams 抓包后参考项目中的image/params.jpeg 把信息一行一行copy到下面 没有的key不需要复制
func GetParams(headers map[string]string) map[string]string {
	params := make(map[string]string)
	params["api_version"] = headers["ddmc-api-version"]
	params["app_version"] = headers["ddmc-build-version"]
	params["app_client_id"] = "4"
	params["applet_source"] = ""
	params["channel"] = "applet"
	params["city_number"] = headers["ddmc-city-number"]
	params["h5_source"] = ""
	params["longitude"] = headers["ddmc-longitude"]
	params["latitude"] = headers["ddmc-latitude"]
	params["openid"] = headers["ddmc-device-id"]
	params["s_id"] = strings.TrimLeft(headers["cookie"], "DDXQSESSID=")
	params["sharer_uid"] = ""
	params["station_id"] = headers["ddmc-station-id"]
	params["time"] = headers["ddmc-time"]
	params["uid"] = headers["ddmc-uid"]

	p := config.Get().Params
	params["device_token"] = p["device_token"] // 设备token
	return params
}

func Sign(params map[string]string) (map[string]string, error) {
	res, err := js.Call(JsFile(), "sign", json.MustEncodeToString(params))
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	json.MustDecodeFromString(res.String(), &m)
	params["nars"] = m["nars"]
	params["sesi"] = m["sesi"]
	return params, nil
}
