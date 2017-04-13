package rrd

// DSTypes data source type
var DSTypes = struct {
	Gauge    string
	Counter  string
	Absolute string
	Derive   string
}{
	"GAUGE",
	"COUNTER",
	"ABSOLUTE",
	"DERIVE",
}

// RRATypes rra type
var RRATypes = struct {
	Average string
	Min     string
	Max     string
	Last    string
}{
	"AVERAGE",
	"MIN",
	"MAX",
	"LAST",
}
