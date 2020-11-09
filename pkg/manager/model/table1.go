package model

import "time"

type Table1 struct {
	Id				int64     `xorm:"not null pk autoincr BIGINT(20)"`
	Gid				string    `xorm:"not null unique VARCHAR(36)"`
	Uuid			string    `xorm:"not null unique VARCHAR(36)"`
	Name			string    `xorm:"not null VARCHAR(36)"`
	Deleted			int       `xorm:"index TINYINT(1)"`
	DeletedAt		time.Time `xorm:"DATETIME"`
	TenantId		string    `xorm:"VARCHAR(36)"`
	Created			time.Time `xorm:"DATETIME"`
	Updated			time.Time `xorm:"DATETIME updated"`
}