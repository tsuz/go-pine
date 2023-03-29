package testutil

import "time"

func GenerateTime() time.Time {
	return time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
}
