package models

import (
	"time"
)

type message struct {
	id        int
	CreatedAt time.Time `gorm:"comment:创建时间"`
	From      int
	FromName  string
	To        int
	ToName    string
	Msg       string
}
