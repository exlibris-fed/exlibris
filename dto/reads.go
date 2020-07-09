package dto

import "time"

type Read struct {
	Book
	Timestamp time.Time
}
