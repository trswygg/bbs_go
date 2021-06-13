package controllers

import (
	"bbs/models"
	"bbs/util"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	//init rand seed
	rand.Seed(time.Now().UnixNano())
}

// UserController Operations about Users
type UserController struct {
	web.Controller
}

func (c *UserController) URLMapping() {
	logs.Trace("[URLMapping] mapping UserController")
	c.Mapping("login", c.Login)
	c.Mapping("logout", c.Logout)
	c.Mapping("info", c.GetUser)
	c.Mapping("info", c.UpdateUser)
	c.Mapping("reg", c.RegUser)
	c.Mapping("profile", c.GetProfile)
	c.Mapping("profile", c.UpdateProfile)

}

type UserControllerReply struct {
	Time      int64       `json:"timestamp"`
	ReplyCode int         `json:"code"`
	Result    string      `json:"result"`
	Data      interface{} `json:"data"`
}
type LastLogin struct {
	Time string `json:"time"`
	IP   string `json:"ip"`
}

// Verify 查询AccessKey所代表的UID
func Verify(Key string) int {
	acc := models.EmptyAccess()
	acc.Key = Key
	// 命中
	if acc.Query() {
		return acc.UID
	}
	return -1
}

func (c *UserController) sendReply(code int, result string, data interface{}) {
	c.Data["json"] = UserControllerReply{
		Time:      time.Now().Unix(),
		ReplyCode: code,
		Result:    result,
		Data:      data,
	}
	err := c.ServeJSON()
	if err != nil {
		logs.Error("sendReply() err:", err.Error())
	}
	return
}

// Login
// router /login [post]
func (c *UserController) Login() {
	Key := c.GetString("AccessKey")
	if Key != "" {
		// 有缓存
		if id, b := models.GetUserIdByKey(Key); b != false {
			// 命中
			// 已经登录
			c.sendReply(-1, "已经登录,id:"+strconv.Itoa(id), nil)
			return
		}
	}
	user := new(models.User)
	email := c.GetString("email")
	password := c.GetString("password")
	// SELECT
	res := util.DB.Where("email = ?", email).Where("password = ?", password).First(&user)
	logs.Trace("login email:", email, "password:", password)
	if res.Error != nil {
		c.sendReply(-1, "err:"+res.Error.Error(), nil)
		return
	}
	if res.RowsAffected == 1 {
		// 密码正确
		// update lastLogin info
		buf, _ := json.Marshal(&LastLogin{
			Time: time.Now().Format(time.UnixDate),
			IP:   c.Ctx.Input.IP(),
		})
		util.DB.First(&models.User{}, user.Id).Update("last_login", string(buf))
		access, _ := models.NewAccess(user.Id)
		c.Ctx.SetCookie("AccessKey", access.Key, 3600)
		c.sendReply(0, "login success", access)
		return
	} else {
		// 密码错误
		c.sendReply(-3, "wrong username or password,please retry", nil)
		return
	}
}

// Logout 退出登录
// method post
// router /logout [post]
func (c *UserController) Logout() {
	Key := c.GetString("AccessKey")
	id, _ := c.GetInt("id", 0)
	if Key != "" {
		// 有缓存
		acc := models.EmptyAccess()
		acc.UID = id
		acc.Key = Key
		if acc.Verify() {
			acc.Drop()
			c.Ctx.SetCookie("AccessKey", "")
			c.sendReply(0, "logout success", nil)
			return
		} else {
			c.Ctx.Output.Status = 403
			c.sendReply(-1, "wrong access", nil)
			return
		}
	}
	c.Ctx.Output.Status = 403
	c.sendReply(-1, "wrong access", nil)
	return
}

// GetUser 获取user信息
// method get
// router /info:id [get]
func (c *UserController) GetUser() {
	user := new(models.User)
	requireId := c.GetString("id")
	// get user
	util.DB.Joins("Profile").First(&user, requireId)
	if user.Name != "" {
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Favorites{}).Count(&user.FavoritesCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Likes{}).Count(&user.LikesCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Images{}).Count(&user.ImagesCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Class{}).Count(&user.ClassesCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Tribe{}).Count(&user.TribesCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Post{}).Count(&user.PostsCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Reply{}).Count(&user.ReplysCount)
		util.DB.Where("creator_id = ?", user.Id).Model(&models.Follows{}).Count(&user.FollowsCount)
		util.DB.Where("user_id = ?", user.Id).Model(&models.Follows{}).Count(&user.FansCount)
		c.sendReply(0, "success", user)
		return
	} else {
		c.sendReply(-1, "no such user id:"+requireId, nil)
		return
	}
}

// UpdateUser 更新user信息
// method post
// Param  form-data face name email password
// router /info [post]
func (c *UserController) UpdateUser() {
	user := new(models.User)
	password := c.GetString("password")
	email := c.GetString("email")
	name := c.GetString("name")
	face := c.GetString("face")
	// get user
	Key := c.GetString("AccessKey")
	if Key != "" {
		// 有缓存
		acc := models.EmptyAccess()
		acc.Key = Key
		if acc.Query() == true {
			util.DB.Find(user, acc.UID)
			// update info
			// give exp 3
			user.Exp += 3
			//  "" 0 false 不会被更新
			util.DB.Model(&user).Updates(models.User{Name: name, Email: email, Password: password, Face: face})
			c.sendReply(0, "success", user)
			return
		} else {
			logs.Trace("acc:", acc)
			c.sendReply(-2, "wrong access", nil)
			return
		}
	} else {
		c.sendReply(-1, "no AccessKey found", nil)
		return
	}
}

// RegUser 注册user
// method post
// router /reg [post]
func (c *UserController) RegUser() {
	buf := new(models.User)

	userGroup := new(models.UserGroup)
	util.DB.Where("group_name = ?", "user").First(userGroup)
	name := c.GetString("name")
	email := c.GetString("email")
	password := c.GetString("password", "123456")
	user := &models.User{
		Name:        name,
		Email:       email,
		Password:    password,
		Prestige:    0,
		UserGroupId: userGroup.Id,
	}
	// get user
	if util.DB.Where("email = ?", email).Find(&buf).RowsAffected < 1 {

		res := util.DB.Create(user)
		if res.Error != nil {
			c.sendReply(-1, res.Error.Error(), nil)
			return
		}
		c.sendReply(0, "success", user)
		return
	} else {
		c.sendReply(-2, "duplicate email", nil)
		return
	}

}

// GetProfile 获取user profile
// method get
// router /profile :id [get]
func (c *UserController) GetProfile() {
	user := new(models.User)
	profile := new(models.Profile)
	requireId, errn := strconv.Atoi(c.GetString("id"))
	hasUser := util.DB.Find(&user, requireId).RowsAffected
	if errn == nil && requireId >= 0 && hasUser == 1 {
		// get user
		if util.DB.Where("user_id = ?", requireId).Find(profile).RowsAffected == 0 {
			// new profile
			util.DB.Create(models.NewProfile(requireId))
		} // do nothing
		c.sendReply(0, "success", profile)
		return
	} else {
		if errn != nil {
			logs.Error("GetProfile() :", errn.Error())
		}
		c.sendReply(-1, "no such user id:"+strconv.Itoa(requireId), nil)
		return
	}
}

// UpdateProfile 更新user profile
// method get
// router /profile [post]
func (c *UserController) UpdateProfile() {
	profile := new(models.Profile)
	// get user
	Key := c.GetString("AccessKey")
	if Key != "" {
		// 有缓存
		acc := models.EmptyAccess()
		acc.Key = Key
		if acc.Query() == true {
			util.DB.Where("user_id = ?", acc.UID).First(&profile)
			// update info
			sign := c.GetString("sign")
			phoneVisibility, _ := c.GetInt8("phone_visibility", 0)
			phone := c.GetString("phone")
			emailVisibility, _ := c.GetInt8("email_visibility", 0)
			email := c.GetString("email")
			ageVisibility, _ := c.GetInt8("age_visibility", 0)
			age, _ := c.GetInt8("age", 0)
			sexVisibility, _ := c.GetInt8("sex_visibility", 0)
			sex := c.GetString("sex")
			birthdayVisibility, _ := c.GetInt8("birthday_visibility", 0)
			birthday, _ := c.GetInt64("birthday", 0)

			buf := &models.Profile{
				Sign:               sign,
				Phone:              phone,
				PhoneVisibility:    phoneVisibility,
				Email:              email,
				EmailVisibility:    emailVisibility,
				Age:                age,
				AgeVisibility:      ageVisibility,
				Sex:                sex,
				SexVisibility:      sexVisibility,
				Birthday:           time.Unix(birthday, 0),
				BirthdayVisibility: birthdayVisibility,
			}
			util.DB.Model(&profile).Where("user_id = ?", acc.UID).Updates(buf)
			c.sendReply(0, "success", buf)
			return
		} else {
			c.sendReply(-1, "no AccessKey found", nil)
			return
		}
	} else {
		c.sendReply(-1, "no AccessKey found", nil)
		return
	}
}

// DoFavorites 添加收藏项
// method post
// router /favorites [post]
func (c *UserController) DoFavorites() {

}

// GetFavorites 获取收藏
// method get
// router /favorites [get]
func (c *UserController) GetFavorites() {

}

// GetLevel
// max Exp 2147483647
// https://www.zhihu.com/question/63595057
func GetLevel(Exp int32) int {
	// max Level:6
	if Exp >= 28800 {
		return 6
	} else if Exp >= 10800 {
		return 5
	} else if Exp >= 4500 {
		return 5
	} else if Exp >= 1500 {
		return 3
	} else if Exp >= 200 {
		return 2
	} else if Exp >= 10 {
		return 1
	}
	return 0
}

// GiveExp 增加`Id`对应user的`Exp`
func GiveExp(Exp int32, Id int) {
	MaxExp := 28800
	// data prepare
	user := new(models.User)
	util.DB.Find(user, Id)
	// sql return 0 rows
	if user.Name == "" {
		logs.Error("GiveLevel():", "Wrong user `Id` :", Id)
		return
	} else {
		// begin
		oldExp := user.Exp
		buf := oldExp + Exp
		newExp := int32(math.Min(float64(buf), float64(MaxExp)))
		logs.Trace("[Trace] GiveLevel()", user.Id, user.Name, oldExp, "->", newExp)
		user.Exp = newExp

		// save
		util.DB.Save(user)
	}
}
