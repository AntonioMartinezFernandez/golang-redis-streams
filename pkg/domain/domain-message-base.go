package domain

import (
	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

type DomainMessageBase struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func NewDomainMessageBase(messageType string) *DomainMessageBase {
	return &DomainMessageBase{
		Id:   pkg_utils.NewUuid(),
		Type: messageType,
	}
}

func (bm *DomainMessageBase) GetId() string {
	return bm.Id
}

func (bm *DomainMessageBase) GetType() string {
	return bm.Type
}
