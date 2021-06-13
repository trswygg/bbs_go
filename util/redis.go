package util

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	db "github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

//                _._
//           _.-``__ ''-._
//      _.-``    `.  `_.  ''-._           Redis 6.2.1 (00000000/0) 64 bit
//  .-`` .-```.  ```\/    _.,_ ''-._
// (    '      ,       .-`  | `,    )     Running in standalone mode
// |`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379
// |    `-._   `._    /     _.-'    |     PID: 1
//  `-._    `-._  `-./  _.-'    _.-'
// |`-._`-._    `-.__.-'    _.-'_.-'|
// |    `-._`-._        _.-'_.-'    |           http://redis.io
//  `-._    `-._`-.__.-'_.-'    _.-'
// |`-._`-._    `-.__.-'    _.-'_.-'|
// |    `-._`-._        _.-'_.-'    |
//  `-._    `-._`-.__.-'_.-'    _.-'
//      `-._    `-.__.-'    _.-'
//          `-._        _.-'
//              `-.__.-'

type Redis struct {
	pool *db.Pool
}

// 启动顺序有问题，只能先写死
var (
	Method   string
	Url      string
	PASSWORD string
)

func init() {

	Method, _ = web.AppConfig.String("redis.Method")
	Url, _ = web.AppConfig.String("redis.Url")
	PASSWORD, _ = web.AppConfig.String("redis.PASSWORD")
	logs.Info("[Redis] init()", Method, Url, PASSWORD)
}

func NewRedis() *Redis {
	redis := new(Redis)
	logs.Trace("redis ", Method, Url, PASSWORD)
	redis.pool = &db.Pool{
		MaxIdle:     16,
		IdleTimeout: 20 * time.Second,
		Dial: func() (db.Conn, error) {
			c, err := db.Dial(
				Method,
				Url,
			)
			if err == nil {
				_, err := c.Do("AUTH", PASSWORD)
				if err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c db.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return redis
}

//MaxIdle:     3,
//IdleTimeout: 240 * time.Second,
//Dial: func() (redis.Conn, error) {
//	c, err := redis.Dial("tcp", server)
//	if err != nil {
//		return nil, err
//	}
//	if password != "" {
//		_, err := c.Do("AUTH", password)
//		if err != nil {
//			c.Close()
//			return nil, err
//		}
//	}
//	return c, err
//},
//TestOnBorrow: func(c redis.Conn, t time.Time) error {
//	_, err := c.Do("PING")
//	return err
//},

func (r *Redis) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}
	logs.Trace("REDIS | EXEC() | ", cmd, " | ", params)
	return con.Do(cmd, params...)
}

// Expire 设置过期时间(秒)
func (r *Redis) Expire(name string, ExpireTime int64) error {
	// 设置key 的过期时间
	_, err := r.Exec("EXPIRE", name, ExpireTime)
	return err
}

// Expireat 设置过期时间(Time)
func (r *Redis) Expireat(name string, ExpireatTime time.Time) error {
	// 设置key 的过期时间
	_, err := r.Exec(" EXPIREAT", name, int32(ExpireatTime.Unix()))
	return err
}

// StringGet
// 获取value
// key:string value:string
func (r *Redis) StringGet(Key string) (string, error) {
	str, err := db.String(r.Exec("GET", Key))
	return str, err
}

// key:string value:Int
func (r *Redis) IntGet(Key string) (int, error) {
	n, err := db.Int(r.Exec("GET", Key))
	return n, err
}

// 设置 value
// 总是返回OK，故不处理错误

// key:string value:string
func (r *Redis) StringSet(Key, Value string) {
	_, _ = r.Exec("SET", Key, Value)

}

//key:string value:int
func (r *Redis) IntSet(Key string, value int) {
	_, _ = r.Exec("SET", Key, strconv.Itoa(value))

}

// Exist 获取key是否存在
func (r *Redis) Exist(Key string) (bool, error) {
	b, err := db.Bool(r.Exec("EXISTS", Key))
	return b, err
}

// Delete
// 删除键值对
// return should be 1
func (r *Redis) Delete(Key string) (int, error) {
	n, err := db.Int(r.Exec("DEL", Key))
	return n, err
}
