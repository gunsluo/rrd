package rrd

import "time"

type DataSource struct {
	Name  string
	Type  string
	Start time.Time
	Step  uint
	Max   int
	Min   int
}
