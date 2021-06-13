package models

import (
	"bbs/util"
	"time"
)

// Class 分类
// relation
// 		Class `Has Many` Tribe
// 		Class `Belongs To` User
type Class struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	CreatorId int `gorm:"commit:创建者;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Name        string `gorm:"size:20;index"`
	Disc        string `gorm:"commit:简介"`
	LevelLimit  int    `gorm:"commit:等级限制"`
	UserGroupId int    `gorm:"commit:用户组限制"`
}

func (c *Class) Insert() int {
	util.DB.Table("class").Create(c)
	return c.Id
}

func NewClass(creatorId int, name string, disc string, levelLimit int, userGroupId int) *Class {
	return &Class{CreatorId: creatorId, Name: name, Disc: disc, LevelLimit: levelLimit, UserGroupId: userGroupId}
}

// Tribe 部落
// relation:
// 		Class `Has Many` Tribe
// 		Tribe `Has Many` Post
type Tribe struct {
	Id             int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt      time.Time `gorm:"comment:创建时间"`
	UpdatedAt      time.Time `orm:"comment:修改时间"`
	DeletedAt      time.Time `orm:"comment:删除时间"`
	ActivityDegree int       `gorm:"comment:活跃度"`
	Protect        bool      `orm:"commit:受保护"`

	Name  string `gorm:"size:10"`
	Color string `gorm:"size:7;type:char(7);commit:标题栏颜色"`
	Disc  string

	//
	Creator   User
	CreatorId int `gorm:"comment:创建人;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// Tribe `Belongs To` Class
	Class   Class
	ClassId int

	// Tribe `Has Many` Class
	Posts []Post `gorm:"foreignKey:TribeId"`
}
type TribeQry struct {
	Id             int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt      time.Time `gorm:"comment:创建时间"`
	ActivityDegree int       `gorm:"comment:活跃度"`
	Protect        bool      `orm:"commit:受保护"`
	Name           string    `gorm:"size:10"`
	Color          string    `gorm:"size:7;type:char(7);commit:标题栏颜色"`
	Disc           string
	CreatorId      int `gorm:"comment:创建人;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ClassId        int
	PostsCount     int `gorm:"-"`
}

func NewTribe(protect bool, name string, color string, disc string, creatorId int, classId int) *Tribe {
	return &Tribe{ActivityDegree: 0, Protect: protect, Name: name, Color: color, Disc: disc, CreatorId: creatorId, ClassId: classId}
}

func (t *Tribe) Insert() int {
	util.DB.Table("tribe").Create(t)
	return t.Id
}

type Tag struct {
	Id    int    `gorm:"primaryKey;autoIncrement;type:serial"`
	Name  string `gorm:"not null;unique;index"`
	Color string `gorm:"not null;default:"000000"`
}

func (t *Tag) Insert() int {
	util.DB.Table("tag").Create(t)
	return t.Id
}

func NewTag(name string, color string) *Tag {
	return &Tag{Name: name, Color: color}
}

// Post 帖子
// relation:
// 		Post `Has Many` Reply
// 		Post `Has Many` Favorites
// 		Tribe `Has Many` Post
// 		Post `Many To Many` Tags
// MarkDown : https://ld246.com/article/1549638745630
type Post struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	TribeId   int    `gorm:"comment:所属主题"`
	CreatorId int    `gorm:"comment:创建人;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MainImage string `gorm:"comment:${NetIo.baseUrl}/images/${data.mainImage}"`
	Title     string `gorm:"size:50;not null;index"`
	Content   string `gorm:"type:text;not null;comment:内容(MD);"`

	Tags []Tag `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Replys    []Reply     `gorm:"foreignKey:PostId"` // 回复 reverse
	Favorites []Favorites `gorm:"foreignKey:PostId"` // 收藏 reverse
	Likes     []Likes     `gorm:"foreignKey:PostId"` // 点赞 reverse

	FavoritesCount int `gorm:"-"`
	LikesCount     int `gorm:"-"`
	ReplyCount     int `gorm:"-"`

	CreatorName     string `gorm:"-"`
	CreatorFace     string `gorm:"-"`
	AlreadyFavorite bool   `gorm:"-"`
	AlreadyLike     bool   `gorm:"-"`
	TribeName       string `gorm:"-"`
	ClassName       string `gorm:"-"`
}

type PostQuery struct {
	Id        int
	CreatorId int
	TribeId   int
	Title     string
	MainImage string
	CreatedAt time.Time
}

func (p *Post) Insert() int {
	util.DB.Table("post").Create(p)
	return p.Id
}

func NewPost(tribeId int, creatorId int, title string, content string, MainImage string, tags []Tag) *Post {
	return &Post{TribeId: tribeId, CreatorId: creatorId, Title: title, Content: content, MainImage: MainImage, Tags: tags}
}

// Reply 回复
// Reply `Has Many` SubReply
type Reply struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	CreatorId int    `gorm:"comment:创建人;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;index"`
	PostId    int    `gorm:"comment:帖子;index"` // fk:Post.Id
	Content   string `gorm:"size:50;not null;comment:内容(MD)"`

	ParentId int `gorm:"commit:父节点;null"` // fk:Reply.Id // 0:root
}

func NewReply(creatorId int, postId int, parentId int, content string) *Reply {
	return &Reply{CreatorId: creatorId, PostId: postId, Content: content, ParentId: parentId}
}

func (r *Reply) Insert() int {
	if r.ParentId != 0 {
		util.DB.Table("reply").Create(r)
		return r.Id
	} else {
		util.DB.Table("reply").Select("CreatedAt", "CreatedAt", "CreatorId", "PostId", "Content").Create(r)
		return r.Id
	}
}
