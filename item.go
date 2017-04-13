package rrd

// Item rrd item data
type Item struct {
	Value     JsonFloat `json:"value"`
	Timestamp int64     `json:"timestamp"`
}
