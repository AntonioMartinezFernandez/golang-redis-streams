package utils

import (
	"math/rand"
	"time"

	ulid "github.com/oklog/ulid/v2"
)

func NewUlid() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	timestamp := ulid.Timestamp(time.Now())
	newUlid, _ := ulid.New(timestamp, entropy)
	return newUlid.String()
}
