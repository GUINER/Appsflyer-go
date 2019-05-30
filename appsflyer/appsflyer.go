package appsflyer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	AppID               = "com.your_appid"
	AppsFlyerUrlFormat  = "https://api2.appsflyer.com/inappevent/%s"
	AppsFlyerDevKey     = "your_appsflyer_devkey"
	AppsFlyerTimeFormat = "2006-01-02 15:04:05.000" //utc时间格式
)

type Event struct {
	AppID  string         `json:"app_id"`  // AppsFyler的ID
	DevKey string         `json:"dev_key"` // 校验key
	Body   AppsFlyerEvent `json:"body"`    // 传输数据
}

type AppsFlyerEvent struct {
	AppsflyerID    string `json:"appsflyer_id"`               // 手机设备的AppsflyerID
	AdvertisingID  string `json:"advertising_id,omitempty"`   // 安卓设备ID（可选)
	IDFA           string `json:"idfa,omitempty"`             // IOS（可选)
	CustomerUserId string `json:"customer_user_id,omitempty"` // IOS（可选)
	BundleID       string `json:"bundle_id,omitempty"`        // IOS（可选)
	EventName      string `json:"eventName"`                  // 事件名称
	EventValue     string `json:"eventValue"`                 // 事件值（可选），json string or ''
	EventTime      string `json:"eventTime"`                  // 时间, utc时间
	AfEventsApi    string `json:"af_events_api"`              // 默认为true
}

//
type AppsFlyer struct {
	Client *http.Client
	event  *Event
}

// json to string
// eventValue: map or struct
func (s *AppsFlyer) Json2String(eventValue interface{}) (jsonStr string) {
	jsonByte, err := json.Marshal(eventValue)
	if err != nil {
		fmt.Printf("EncodeEventValue error: %v", err)
		return
	}
	jsonStr = string(jsonByte)
	return
}

// AppsFlyer打点
func (s *AppsFlyer) SendEvent(event *Event) (err error) {
	if event == nil {
		return fmt.Errorf("SendEvent# param event is nil")
	}
	s.event = event
	if err = s.checkValid(); err != nil {
		return
	}
	s.checkClient()
	if err = s.post2AppsFlyer(); err != nil {
		return
	}
	return
}

func (s *AppsFlyer) checkValid() (err error) {
	if s.event.AppID == "" {
		return fmt.Errorf("AppID is null")
	}
	if s.event.DevKey == "" {
		return fmt.Errorf("DevKey is null")
	}
	if s.event.Body.AppsflyerID == "" {
		return fmt.Errorf("AppsflyerID is null")
	}
	if s.event.Body.EventName == "" {
		return fmt.Errorf("EventName is null")
	}
	return
}

//
func (s *AppsFlyer) checkClient() {
	if s.Client == nil {
		s.NewClient()
	}
}

// http client
func (s *AppsFlyer) NewClient() {
	s.Client = &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}}
}

var codeList = []string{
	200: "OK when the request has been processed by the AppsFlyer system.",
	401: "unauthorized when the key provided in the authentication header is not the dev-key for this app.",
	400: "Bad request when the request failed at least one of the validation criteria",
	500: "Internal server error this indicates a server error",
}

// post to AppsFlyer
func (s *AppsFlyer) post2AppsFlyer() (err error) {
	params, err := json.Marshal(s.event.Body)
	if err != nil {
		return
	}

	appsFlyerUrl := fmt.Sprintf(AppsFlyerUrlFormat, s.event.AppID)
	req, err := http.NewRequest(http.MethodPost, appsFlyerUrl, bytes.NewReader(params))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authentication", s.event.DevKey)

	resp, err := s.Client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	fmt.Println("Post to AppsFlyer ----->", "statusCode =", resp.StatusCode, "eventName =", s.event.Body.EventName)
	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode:%d, ErrMsg:%s", resp.StatusCode, codeList[resp.StatusCode])
	}
	return
}
