package main

import (
	"Appsflyer-go/appsflyer"
	"fmt"
	"time"
)

func main() {
	af := appsflyer.AppsFlyer{}
	eventValue := map[string]interface{}{
		"name": "gg",
		"age":  18,
	}
	event := appsflyer.Event{
		DevKey: appsflyer.AppsFlyerDevKey,
		AppID:  appsflyer.AppID,
		Body: appsflyer.AppsFlyerEvent{
			AppsflyerID: "your appsFlyer id",
			EventName:   "appsflyer_event",
			EventValue:  af.Json2String(eventValue), // json 字符串
			EventTime:   time.Now().In(time.UTC).Format(appsflyer.AppsFlyerTimeFormat),
			AfEventsApi: "true", // 必须为true
		},
	}
	if err := af.SendEvent(&event); err != nil {
		fmt.Printf("AppsFlyer SendEvent error, %v\n", err)
		return
	}
}
