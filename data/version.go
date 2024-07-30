package data

import "time"

func getVersion() int64 {
	return time.Now().Unix()
}
