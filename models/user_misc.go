package models

import "time"

// Favorites 收藏夹
// relations
//  	User `Has Many` Favorites
//  	Post `Has Many` Favorites
type Favorites struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	CreatorId int `gorm:""`
	PostId    int `gorm:"comment:帖子"`
	PostTitle string
}

// Likes 点赞记录
// relations
//  	User `Has Many` Likes
//  	Post `Has Many` Likes
type Likes struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	CreatorId int `gorm:""`
	PostId    int `gorm:"comment:帖子"`
}

// Follows 关注列表
// relations
//  	User `Has Many` Follows
type Follows struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	CreatorId int
	UserId    int `gorm:"comment:被关注者;index;"`
}

// Images 关注列表
// relations
//  	User `Has Many` Images
type Images struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:修改时间"`

	CreatorId   int
	CreatorName string `gorm:"-"`
	CreatorFace string `gorm:"-"`
	ImageName   string `gorm:"index;unique;not null"` // 写死的子路径 /static/images/XXX.jpg
}
