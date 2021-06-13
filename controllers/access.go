package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"time"
)

type AccessController struct {
	web.Controller
}

type AccessControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Resources string
	Count     int
	Data      interface{}
}

// id	group_name		display_name	disc								level	permission_code
// 1	root			版主				不服来干								5		0
// 2	admin			管理员			系统管理员，拥有管理所有板块与用户的权限	4		0
// 3	authorizeduser	授权用户			授权用户，拥有管理所授权板块的权限		3		0
// 4	user			普通用户			普通用户，拥有基础的交互功能				2		0
// 5	guest			访客				访客，允许浏览基础的网页				1		0
// 6	noaccess		禁止访问			禁止访问，账户禁封中					-1		0
// 7	default			default			default								0		0
// auth || character
// CharacterUsers || CharacterAuths

func (c *AccessController) URLMapping() {
	logs.Trace("[URLMapping] mapping AccessController")

}
func (c *AccessController) sendReply(code int, result string, resources string, count int, data interface{}) {
	c.Data["json"] = AccessControllerReply{
		Time:      time.Now().Unix(),
		ReplyCode: code,
		Result:    result,
		Resources: resources,
		Count:     count,
		Data:      data,
	}
	err := c.ServeJSON()
	if err != nil {
		logs.Error("sendReply() err:", err.Error())
	}
	return
}

func (c AccessController) GetUserGroup() {

}

func (c *AccessController) ChangeUserGroupForUser() {

}

// -=-=-=-=-=-=-=-=-=-=-= Auth =-=-=-=-=-=-=-=-=-=-=-

func (c *AccessController) GetAuth() {

}

func (c *AccessController) FindAuth() {

}

func (c *AccessController) UpdateAuth() {

}

func (c *AccessController) DropAuth() {

}

// -=-=-=-=-=-=-=-=-=-=-= Character =-=-=-=-=-=-=-=-=-=-=-

func (c *AccessController) GetCharacter() {

}

func (c *AccessController) FindCharacter() {

}

func (c *AccessController) UpdateCharacter() {

}

func (c *AccessController) DropCharacter() {

}

// AbleTo 遍历审核表
func AbleTo(uid int, resources string) bool {
	return false
}
