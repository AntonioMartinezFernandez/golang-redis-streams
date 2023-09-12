package utils

import (
	"math/rand"
	"time"

	ulid "github.com/oklog/ulid/v2"
)

func NewUlid() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	newUlid, _ := ulid.New(ms, entropy)
	return newUlid.String()
}
