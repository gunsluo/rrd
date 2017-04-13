package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type ChartData struct {
	DSName    string  `json:"dsName"`
	Mark      string  `json:"mark"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	var c float64
	for {

		c++
		d := &ChartData{
			DSName:    "demo",
			Mark:      "demo",
			Value:     c,
			Timestamp: time.Now().Unix(),
		}
		mJson, _ := json.Marshal(d)

		request := gorequest.New()
		_, body, _ := request.Post("http://127.0.0.1:8800/api/v1/push").Send(string(mJson)).End()
		fmt.Println(body)

		time.Sleep(time.Minute * 1)
	}
}
