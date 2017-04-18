package rrd

import (
	"fmt"
	"math"
)

// JSONFloat float64 type json by rrd
type JSONFloat float64

// MarshalJSON json marshalJson
func (v JSONFloat) MarshalJSON() ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("%f", f)), nil
}
