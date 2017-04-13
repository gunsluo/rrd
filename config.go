package rrd

// Config rrd config
type Config struct {
	RRDDir string
	DSS    []DSConfig
}

// DSConfig rrd ds config
type DSConfig struct {
	Name      string
	Type      string
	Step      uint
	Heartbeat uint
	Max       int
	Min       int
	RRAS      []RRAConfig
}

// RRAConfig rrd rra config
type RRAConfig struct {
	Type  string
	Xff   float64
	Steps int
	Rows  int
}
