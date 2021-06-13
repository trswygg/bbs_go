package models

import (
	"bbs/util"
	"time"
)

// UserGroup 用户组，描述基础权限
// relations:
// 		UserGroup `Has Many` User ;
type UserGroup struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	GroupName      string `gorm:"unique;size:20;comment:组名，简称"`
	DisplayName    string `gorm:"size:10;comment:全称"`
	Disc           string `gorm:"comment:简介"`
	Level          int8   `gorm:"unique;comment:用户组等级"`
	PermissionCode int    `gorm:"comment:保留:"` //TODO 保留备用

	// UserGroup `Has Many` User
	Users []User
}

func NewUserGroup(groupName string, displayName string, disc string, level int8, permissionCode int) *UserGroup {
	return &UserGroup{GroupName: groupName, DisplayName: displayName, Disc: disc, Level: level, PermissionCode: permissionCode}
}

// Insert 向`UserGroup`表中插入数据
func (UserGroup *UserGroup) Insert() int {
	util.DB.Table("user_group").Create(UserGroup)
	return UserGroup.Id
}

// Character 角色，描述用户被赋予的特殊权限
// relations:
// 		User `Many To Many` Character ;
// 		Auth `Many To Many` Character ;
type Character struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	Name string `gorm:"size:10"`

	Users []User `gorm:"many2many:character_users;"`
	Auths []Auth `gorm:"many2many:character_auths;"`
}

// Auth 权限，描述某项特殊的权限，使普通组成员可以管理某论坛
// relations:
// 		Auth `Many To Many` Character ;
type Auth struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	Name  string `gorm:"size:10;not null"`
	OP    string `gorm:"size:20;not null;comment:操作"`
	Scope int    `gorm:"size:20;not null;comment:最小用户组等级"`

	Characters []Character `gorm:"many2many:character_auths;"`
}

// 连接表
// autoIncrement:false

// CharacterUsers rel(m2m)
type CharacterUsers struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	ReviewerId int `gorm:"comment:创建人"`
	Reason     string

	UserId      int `gorm:"primaryKey"`
	CharacterId int `gorm:"primaryKey"`
}

// CharacterAuths rel(m2m)
type CharacterAuths struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	ReviewerId  int
	Reason      string
	AuthId      int `gorm:"primaryKey"`
	CharacterId int `gorm:"primaryKey"`
}
