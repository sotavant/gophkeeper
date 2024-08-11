package data

import "time"

func getVersion() uint64 {
	return uint64(time.Now().Unix())
}
