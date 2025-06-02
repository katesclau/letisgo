package models

import (
	"github.com/oklog/ulid/v2"
)

type Session struct {
	ID  ulid.ULID `json:"id"`
	JWT string    `json:"jwt"`
	TTL int64     `json:"ttl"`
}
