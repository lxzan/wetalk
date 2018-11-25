package wetalk

import "time"

func SetInterval(t time.Duration, cb func()) {
	for {
		cb()
		time.Sleep(t)
	}
}
