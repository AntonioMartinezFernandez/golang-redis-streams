package utils

import uuid "github.com/satori/go.uuid"

func NewUuid() string {
	return uuid.NewV4().String()
}
