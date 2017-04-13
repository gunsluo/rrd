package main

import (
	"fmt"
	"time"

	"github.com/gunsluo/rrd"
)

func main() {

	tool := rrd.NewTool().Cfg(&rrd.Config{
		RRDDir: "./rrd",
		DSS: []rrd.DSConfig{
			{
				Name: "test",
				Type: rrd.DSTypes.Gauge,
				Step: 60,
				Min:  "0",
				Max:  "1000",
				RRAS: []rrd.RRAConfig{
					{
						Type:  rrd.RRATypes.Average,
						Steps: 1,
						Rows:  720,
					}, {
						Type:  rrd.RRATypes.Average,
						Xff:   0.5,
						Steps: 5,
						Rows:  576,
					},
				},
			},
		},
	})

	var items []*rrd.Item
	now := time.Now()
	for i := 0; i < 720; i++ {
		item := new(rrd.Item)
		item.Value = float64(i)
		item.Timestamp = now.Add(time.Duration(i*60) * time.Second).Unix()
		items = append(items, item)
	}
	err := tool.Write("test", "1", items)
	if err != nil {
		panic(err)
	}

	end := now.Add(time.Duration(720*60) * time.Second)
	startTs := now.Unix()
	endTs := end.Unix()
	startTs = startTs - startTs%int64(60) - int64(60)
	endTs = endTs - endTs%int64(60) - int64(60)

	itemsRet, err := tool.Fetch(rrd.RRATypes.Average, "test", "1", rrd.Unix(startTs), rrd.Unix(endTs), 60)
	if err != nil {
		panic(err)
	}

	for idx, v := range itemsRet {
		fmt.Println(idx, "->", v)
	}

	fmt.Println("ok")
}
