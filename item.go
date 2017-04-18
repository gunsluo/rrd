package rrd

// Item rrd item data
type Item struct {
	Value     JSONFloat `json:"value"`
	Timestamp int64     `json:"timestamp"`
}
