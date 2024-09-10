package interfaces

import "gorm.io/gorm"

type Injections struct {
	DB *gorm.DB
}

func NewInjection(db *gorm.DB) *Injections {
	return &Injections{
		DB: db,
	}
}

type Ops struct{}
