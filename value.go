package pine

import "time"

type TPQ struct {
	Timestamp time.Time
	Px        float64
	Qty       float64
}
