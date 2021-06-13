package models

import (
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

func init() {
	//inits()
}

func inits() {
	logs.Info("[models] init(): util.DB.AutoMigrate()")
	logs.Info("[models] init(): SetupJoinTable: CharacterUsers{}")
	erruc := util.DB.SetupJoinTable(&Character{}, "Users", &CharacterUsers{})
	if erruc != nil {
		logs.Error("[models] CharacterUser:", erruc)
	}
	logs.Info("[models] init(): SetupJoinTable: CharacterAuths{}")
	errac := util.DB.SetupJoinTable(&Character{}, "Auths", &CharacterAuths{})
	if errac != nil {
		logs.Error("[models] CharacterAuth:", errac)
	}
	logs.Info("[models] init(): AutoMigrate: Character{} UserGroup{} User{} Profile{} Auth{}")
	err1 := util.DB.AutoMigrate(&Character{}, &UserGroup{}, &User{}, &Profile{}, &Auth{}, &CharacterAuths{}, &CharacterUsers{})
	if err1 != nil {
		logs.Error("[models] AutoMigrate User.go auth.go:", err1)
		return
	}
	logs.Info("[models] init(): AutoMigrate: Class{} Tribe{} Tag{} Post{} Reply{}")
	err2 := util.DB.AutoMigrate(&Class{}, &Tribe{}, &Tag{}, &Post{}, &Reply{})
	if err2 != nil {
		logs.Error("[models] AutoMigrate: Panel.go", err2)
		return
	}
	logs.Info("[models] init(): AutoMigrate: Class{} Tribe{} Tag{} Post{} Reply{}")
	err3 := util.DB.AutoMigrate(&Favorites{}, &Likes{}, &Follows{}, &Images{})
	if err3 != nil {
		logs.Error("[models] AutoMigrate: user_misc.go", err3)
		return
	}

	var userCount int64
	if util.DB.Model(&User{}).Count(&userCount); userCount == 0 {
		logs.Info("[models] init(): DataPrepare()")
		DataPrepare()
	}
}

type templateModule struct {
	// max 2147483647
	Id        int       `gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `orm:"comment:修改时间"`
}

// 名字								描述
// character varying(n), varchar(n)	变长，有长度限制
// character(n), char(n)			定长，不足补空白
// text								变长，无长度限制
// 这三种类型之间没有性能差别，除了当使用填充空白类型时的增加存储空间，
// 和当存储长度约束的列时一些检查存入时长度的额外的CPU周期。虽然在某些
// 其它的数据库系统里，character(n) 有一定的性能优势，但在PostgreSQL
// 里没有。事实上，character(n)通常是这三个中最慢的，因为额外存储成本。
// 在大多数情况下，应该使用text 或character varying。

// 警告：当使用 struct 更新时，GORM只会更新那些非零值的字段
// 对于下面的操作，不会发生任何更新，"", 0, false 都是其类型的零值

// 操作符  返回类型	描述	                    示例	                                                                结果
// @@	  boolean	tsvector是否匹配tsquery	select to_tsvector('fat cats ate rats') @@ to_tsquery('cat & rat');	t

// 1、字段的基本约束，通过gorm tag基本都可以设置。
// 2、gorm支持实体完整性约束。
// 3、域完整性约束中，外键约束需要通过 sql tag或调用api实现，check约束可以直接在字段上定义。
// 4、参照完整性gorm不能默认实现，必须通过sql tag或者调用api实现。
