package loginreq

import "time"

type LoginReq struct {
	ID string
	Username string
	Expiry time.Time
	Dur time.Duration
	Confirmed bool
	Code string
}