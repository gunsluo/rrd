package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gunsluo/rrd"
	"github.com/yubo/rrdlite"
)

func main() {

	//now := time.Now()
	//start := now.Add(time.Duration(-24) * time.Hour)
	start := time.Unix(1491912600, 0)
	filename := "test.rrd"

	ds := &rrd.DataSource{
		Name:  "test",
		Type:  rrd.DataSourceTypes.Gauge,
		Start: start,
		Step:  60,
		Max:   1000,
		Min:   0,
	}
	err := create(filename, ds)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var items []*Item
	for i := 1; i <= 720; i++ {
		item := new(Item)
		item.Value = float64(i)
		item.Timestamp = start.Add(time.Duration(i*60) * time.Second).Unix()
		item.DSType = rrd.DataSourceTypes.Gauge
		items = append(items, item)
	}

	err = update(filename, items)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	end := start.Add(time.Duration(720*60) * time.Second)
	items1, err := fetch(filename, rrd.RRATypes.Average, start.Unix(), end.Unix(), 60)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for idx, v := range items1 {
		fmt.Println(idx, "->", v)
	}

	items2, err := fetch(filename, rrd.RRATypes.Average, start.Unix(), end.Unix(), 60*5)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for idx, v := range items2 {
		fmt.Println(idx, "->", v)
	}

	fmt.Println("ok")
}

func create(filename string, ds *rrd.DataSource) error {

	c := rrdlite.NewCreator(filename, ds.Start, ds.Step)
	c.DS(ds.Name, ds.Type, ds.Step*2, ds.Min, ds.Max)

	// 设置各种归档策略
	c.RRA(rrd.RRATypes.Average, 0.5, 1, 720)
	c.RRA(rrd.RRATypes.Average, 0.5, 5, 576)

	return c.Create(true)
}

type Item struct {
	Value     float64
	Timestamp int64
	DSType    string
}

func update(filename string, items []*Item) error {
	u := rrdlite.NewUpdater(filename)
	//u.SetTemplate("g")

	for _, item := range items {
		v := math.Abs(item.Value)
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if item.DSType == rrd.DataSourceTypes.Derive || item.DSType == rrd.DataSourceTypes.Counter {
			u.Cache(item.Timestamp, int(item.Value))
		} else {
			u.Cache(item.Timestamp, item.Value)
		}
	}

	return u.Update()
}

func fetch(filename string, cf string, start, end int64, step int) ([]*Item, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

	fetchRes, err := rrdlite.Fetch(filename, cf, start_t, end_t, step_t)
	if err != nil {
		return []*Item{}, err
	}

	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	ret := make([]*Item, size)

	start_ts := fetchRes.Start.Unix()
	step_s := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := start_ts + int64(i+1)*int64(step_s)
		d := &Item{
			Timestamp: ts,
			Value:     val,
		}
		ret[i] = d
	}

	return ret, nil
}
