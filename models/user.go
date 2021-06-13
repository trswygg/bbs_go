package models

import (
	"bbs/util"
	_ "gorm.io/gorm"
	"time"
)

// User 用户通用信息
type User struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`
	Name      string    `gorm:"size:15;index;not null"` // find by name
	Email     string    `gorm:"size:30;index;unique;not null"`
	Password  string    `gorm:"comment:md5;not null" json:"-"`
	Face      string    `gorm:"type:text;comment:64*64"` // base64

	Exp      int32 `gorm:"comment:经验值;default:0"`
	Prestige int   `gorm:"comment:威望;default:0"`

	BanTime   time.Time `gorm:"comment:封禁时间"`
	MuteTime  time.Time `gorm:"comment:禁言时间"`
	LastLogin string    `gorm:"size:100;comment:上次登录信息;default:first time"`
	Profile   Profile
	// UserGroup `Has Many` User
	UserGroupId int
	// User `Many To Many` Character ;
	Characters []Character `gorm:"many2many:character_users;" json:"-"`

	//  	User `Has Many` Favorites ； 收藏
	//  	User `Has Many` Likes 	  ； 点赞
	//  	User `Has Many` Images 	  ； 相册
	//  	User `Has Many` Class	  ； 分区
	//  	User `Has Many` Tribe	  ； 部落
	//  	User `Has Many` Post	  ； 帖子
	//  	User `Has Many` Reply	  ； 回复
	//  	User `Has Many` Follows   ； 关注
	Favorites []Favorites `gorm:"foreignKey:CreatorId" json:"-"`
	Likes     []Likes     `gorm:"foreignKey:CreatorId" json:"-"`
	Images    []Images    `gorm:"foreignKey:CreatorId" json:"-"`
	Classes   []Class     `gorm:"foreignKey:CreatorId" json:"-"`
	Tribes    []Tribe     `gorm:"foreignKey:CreatorId" json:"-"`
	Posts     []Post      `gorm:"foreignKey:CreatorId" json:"-"`
	Replys    []Reply     `gorm:"foreignKey:CreatorId" json:"-"`
	Follows   []Follows   `gorm:"foreignKey:CreatorId" json:"-"`
	// view
	FavoritesCount int64 `gorm:"-"`
	LikesCount     int64 `gorm:"-"`
	ImagesCount    int64 `gorm:"-"`
	ClassesCount   int64 `gorm:"-"`
	TribesCount    int64 `gorm:"-"`
	PostsCount     int64 `gorm:"-"`
	ReplysCount    int64 `gorm:"-"`
	FollowsCount   int64 `gorm:"-"`
	FansCount      int64 `gorm:"-"`
}

func NewUser(name string, email string, password string, exp int32, groupId int) *User {
	return &User{Name: name, Email: email, Password: password, Exp: exp, UserGroupId: groupId}
}

// Insert 向`user`表中插入数据(包括profile)
func (user *User) Insert() int {
	util.DB.Table("user").Create(user)
	util.DB.Table("profile").Create(&Profile{UserId: user.Id})
	return user.Id
}

// Profile 用户主页及相关设置
// 1tom 关联：帖子，收藏夹，点赞记录，相册，关注表，粉丝表
// see user_misc.go
// *Visibility: 可见性[0:自己，1:粉丝,2:访客,3:注册用户,all:0]
// relations:
// 		Profile `Belongs To` User ; UserId 是其外键
type Profile struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`

	UserId    int    `gorm:"index"` //严格模式，唯一对应
	MainImage string //主页图片
	// misc
	Sign               string    `gorm:"size:127;comment:个性签名;default:这个人很懒，什么都没有写"`
	Phone              string    `gorm:"comment:手机号;size:15"`
	PhoneVisibility    int8      `gorm:"comment:手机号可见性;not null;default:0" json:"-"`
	Email              string    `gorm:"comment:邮箱;size:32"`
	EmailVisibility    int8      `gorm:"comment:邮箱可见性;not null;default:0" json:"-"`
	Age                int8      `gorm:"comment:年龄"`
	AgeVisibility      int8      `gorm:"comment:年龄可见性;not null;default:0" json:"-"`
	Sex                string    `gorm:"comment:性别;size:8"`
	SexVisibility      int8      `gorm:"comment:性别可见性;not null;default:0" json:"-"`
	Birthday           time.Time `gorm:"comment:生日"`
	BirthdayVisibility int8      `gorm:"comment:生日可见性;not null;default:0" json:"-"`
	//Image           time.Time `gorm:"comment:相册"`
	ImageVisibility int8 `gorm:"comment:相册可见性;not null;default:0"`
}

func NewProfile(userId int) *Profile {
	return &Profile{UserId: userId, Sex: "保密", Birthday: time.Unix(0, 0)}
}

// Insert 向`user`表中插入数据
func (profile *Profile) Insert() int {
	util.DB.Table("profile").Create(profile)
	return profile.Id
}

// RewardRecord 奖惩记录
type RewardRecord struct {
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`

	User   User
	UserId int

	Reviewer   User   `gorm:"foreignKey:ReviewerId"`
	ReviewerId int    `gorm:"comment:创建人"`
	Reason     string `gorm:"not null"`

	Exp      int32 `gorm:"comment:经验值;default:0"`
	Prestige int   `gorm:"comment:威望;default:0"`
}
