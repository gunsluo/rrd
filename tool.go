package rrd

import (
	"fmt"
	"math"
	"time"

	"github.com/toolkits/file"
	"github.com/yubo/rrdlite"
)

// Tool rrd tool
type Tool struct {
	rrdDir string
	dss    map[string]*DSConfig
}

// NewTool return new rrd tool
func NewTool() *Tool {
	return &Tool{
		rrdDir: "/tmp/rrd/",
	}
}

// Cfg config rrd tool
func (t *Tool) Cfg(cfg *Config) *Tool {
	if cfg.RRDDir != "" {
		t.rrdDir = cfg.RRDDir
	}

	t.dss = make(map[string]*DSConfig, len(cfg.DSS))
	for _, ds := range cfg.DSS {
		d := new(DSConfig)
		d.Name = ds.Name
		d.Type = ds.Type
		if ds.Step == 0 {
			d.Step = 60
		} else {
			d.Step = ds.Step
		}
		if ds.Heartbeat == 0 {
			d.Heartbeat = d.Step * 2
		} else {
			d.Heartbeat = d.Heartbeat
		}
		if ds.Min == "" {
			d.Min = "U"
		} else {
			d.Min = ds.Min
		}
		if ds.Max == "" {
			d.Max = "U"
		} else {
			d.Max = ds.Max
		}

		for _, rra := range ds.RRAS {
			r := new(RRAConfig)
			r.Type = rra.Type
			if rra.Xff == 0 {
				r.Xff = 0.5
			} else {
				r.Xff = rra.Xff
			}
			r.Steps = rra.Steps
			r.Rows = rra.Rows
			d.RRAS = append(d.RRAS, *r)
		}

		t.dss[ds.Name] = d
	}

	return t
}

func (t *Tool) create(ds, dbfile string) error {

	d, ok := t.dss[ds]
	if !ok {
		return fmt.Errorf("no config data source[%s]", ds)
	}

	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	c := rrdlite.NewCreator(dbfile, start, d.Step)
	c.DS(d.Name, d.Type, d.Heartbeat, d.Min, d.Max)

	for _, rra := range d.RRAS {
		c.RRA(rra.Type, rra.Xff, rra.Steps, rra.Rows)
	}

	return c.Create(true)
}

// Write write data to rrd file
func (t *Tool) Write(ds, pk string, items []*Item) error {

	if len(items) == 0 {
		return nil
	}

	d, ok := t.dss[ds]
	if !ok {
		return fmt.Errorf("no config data source[%s]", ds)
	}

	dbfile := formatDBfile(t.rrdDir, ds, d.Type, pk)

	if !file.IsExist(dbfile) {
		baseDir := file.Dir(dbfile)

		err := file.InsureDir(baseDir)
		if err != nil {
			return err
		}

		err = t.create(ds, dbfile)
		if err != nil {
			return err
		}
	}

	u := rrdlite.NewUpdater(dbfile)
	for _, item := range items {
		v := math.Abs(item.Value)
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if d.Type == DSTypes.Derive || d.Type == DSTypes.Counter {
			u.Cache(item.Timestamp, int(item.Value))
		} else {
			u.Cache(item.Timestamp, item.Value)
		}
	}

	return u.Update()
}

// Fetch fetch data from rrd file
func (t *Tool) Fetch(rraType, ds, pk string, start, end time.Time, step int) (items []*Item, err error) {

	d, ok := t.dss[ds]
	if !ok {
		return nil, fmt.Errorf("no config data source[%s]", ds)
	}

	dbfile := formatDBfile(t.rrdDir, ds, d.Type, pk)
	stepDt := time.Duration(step) * time.Second

	fetchRes, err := rrdlite.Fetch(dbfile, rraType, start, end, stepDt)
	if err != nil {
		return []*Item{}, err
	}
	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	ret := make([]*Item, size)

	startTsFromRes := fetchRes.Start.Unix()
	stepDtFromRes := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := startTsFromRes + int64(i)*int64(stepDtFromRes)
		d := &Item{
			Timestamp: ts,
			Value:     val,
		}
		ret[i] = d
	}

	return ret, nil
}
