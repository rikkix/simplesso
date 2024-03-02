package loginreq

import "time"

type LoginReq struct {
	ID string
	Username string
	Expiry time.Time
	Dur int // secs
	Confirmed bool
	Code string
}