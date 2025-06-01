package server

import "time"

func Now() time.Time {
	return time.Now() // want `time.Now\(\) should not be used`
}
