package models

import (
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
)

// Access Redis中实际储存的键值对
type Access struct {
	Key string
	UID int
}

var Redis *util.Redis

func init() {
	Redis = util.NewRedis()
	_, e := redis.String(Redis.Exec("PING", ""))
	if e != nil {
		logs.Error(e)
	}
}

func DoLateInit() {
	if Redis == nil {
		LateInit()
	}
}

func LateInit() {
	Redis = util.NewRedis()
	_, e := redis.String(Redis.Exec("PING", ""))
	if e != nil {
		logs.Error(e)
	}
}

// NewAccess 生成UID对应的Access
func NewAccess(UID int) (Access, error) {
	DoLateInit()

	Key := util.GeneKey()
	Redis.IntSet(Key, UID)

	// 更新存活时间，不应出错
	_ = Redis.Expire(Key, int64(3600*12))

	logs.Trace("NewAccess() ", UID, " | ", Key)
	return Access{Key: Key, UID: UID}, nil
}
func EmptyAccess() *Access {
	return &Access{}
}

// Drop
// 删除Access
// 严格模式，Key与UID必须符合记录
func (a *Access) Drop() bool {
	DoLateInit()

	if a.Verify() == true {
		n, _ := Redis.Delete(a.Key)
		return n == 1
	}
	// 不匹配
	return false
}

// Verify
// 验证Key与UID是否符合记录
// 要求：key非空
func (a *Access) Verify() bool {
	DoLateInit()

	if a.Key == "" {
		return false
	}
	b, _ := Redis.Exist(a.Key)
	if b == false {
		logs.Error("token %s not exist", a.Key)
		return false
	}

	value, err := Redis.IntGet(a.Key)
	// 更新存活时间
	_ = Redis.Expire(a.Key, int64(3600*12))

	if err != nil {
		logs.Error("redis exec 'GET %s' error: %s", a.Key, err.Error())
		return false
	}
	// done
	return value == a.UID
}

// Query
// 根据key查询UID
// success: 自动填充 Access return:true
// fail: return:false
func (a *Access) Query() bool {
	DoLateInit()

	if a.Key == "" {
		return false
	}
	value, err := Redis.IntGet(a.Key)
	if err != nil {
		return false
	}
	// 更新存活时间
	_ = Redis.Expire(a.Key, int64(3600))
	a.UID = value

	return true
}

func GetUserIdByKey(key string) (int, bool) {
	DoLateInit()
	if key == "" {
		return 0, false
	}
	value, err := Redis.IntGet(key)
	if err != nil {
		return 0, false
	}
	// 更新存活时间
	_ = Redis.Expire(key, int64(3600))
	return value, true
}

func GetUserGroupByKey(key string) (string, bool) {
	DoLateInit()
	if key == "" {
		return "", false
	}
	uid, err := Redis.IntGet(key)
	if err != nil {
		return "", false
	}
	// 更新存活时间
	_ = Redis.Expire(key, int64(3600))
	user := new(User)
	UserGroup := new(UserGroup)
	resU := util.DB.First(user, uid)
	if resU.Error == nil {
		resG := util.DB.First(UserGroup, user.UserGroupId)
		if resG.Error == nil {
			return UserGroup.GroupName, true
		}
	}
	return "", false
}
