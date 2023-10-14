package domain

type DomainMessage interface {
	GetType() string
	AsMap() map[string]interface{}
}
