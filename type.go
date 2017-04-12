package rrd

var DataSourceTypes = struct {
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
