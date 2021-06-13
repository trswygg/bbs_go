package controllers

import (
	"bbs/models"
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"time"
)

type TagController struct {
	web.Controller
}
type TagControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Count     int
	Data      interface{}
}

func (c *TagController) URLMapping() {
	logs.Trace("[URLMapping] mapping TagController")
	c.Mapping("tag", c.GetTagInfo)
	c.Mapping("tag/query", c.QueryTag)
	c.Mapping("tag/create", c.CreateTag)
	c.Mapping("tag/delete", c.DeleteTag)
}

func (c *TagController) sendReply(code int, result string, count int, data interface{}) {
	c.Data["json"] = TagControllerReply{
		Time:      time.Now().Unix(),
		ReplyCode: code,
		Result:    result,
		Count:     count,
		Data:      data,
	}
	err := c.ServeJSON()
	if err != nil {
		logs.Error("sendReply() err:", err.Error())
	}
	return
}

//oO0OoO0OoO0Oo Tag oO0OoO0OoO0Oo

// QueryTag 获取tag列表
// method get
// router /tag/:text:limit:offset [get]
func (c *TagController) QueryTag() {
	tag := new(models.Tag)
	var resultArr []models.Tag
	requireLimit, errl := c.GetInt("limit", 10)
	requireOffset, erri := c.GetInt("offset", 0)
	if errl != nil {
		c.sendReply(-1, "wrong param: limit", 0, nil)
		return
	} else if erri != nil {
		c.sendReply(-1, "wrong param index", 0, nil)
		return
	}
	// get user
	var total int64
	util.DB.Model(tag).Count(&total)
	util.DB.Model(tag).Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
	c.sendReply(0, "success", int(total), resultArr)
	return
}

// GetTagInfo 获取tag详细信息
// method get
// router /tag_info/:id [get]
func (c *TagController) GetTagInfo() {
	tag := new(models.Tag)
	requireId := c.GetString("id")
	// get user
	rows := util.DB.Find(&tag, requireId).RowsAffected
	if rows == 1 {
		c.sendReply(0, "success", 1, tag)
		return
	} else {
		c.sendReply(-1, "no such tag id:"+requireId, 0, nil)
		return
	}
}

// CreateTag 创建tag
// method post
func (c *TagController) CreateTag() {
	requireName := c.GetString("name")
	requireColor := c.GetString("color", "#000000")
	// get user
	_, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	if requireName == "" {
		c.sendReply(-1, "tag name must be NOT NULL", 0, nil)
		return
	}
	tag := models.NewTag(requireName, requireColor)
	result := util.DB.Create(tag)
	if result.Error != nil {
		c.sendReply(-1, "err:"+result.Error.Error(), 0, nil)
		return
	} else {
		c.sendReply(0, "success", int(result.RowsAffected), tag)
		return
	}
}

// DeleteTag 删除tag
// method post
// router /tag [post]
func (c *TagController) DeleteTag() {
	tag := new(models.Tag)
	id, err := c.GetInt("id")
	if err != nil {
		c.sendReply(-1, "wrong param id:"+err.Error(), 0, nil)
		return
	}
	group, b := models.GetUserGroupByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	if group == "root" || group == "admin" {
		//err := util.DB.Debug().Model(&models.Post{}).Association("Tags").Delete(&models.Tag{Id: id})
		resr := util.DB.Debug().Raw(`DELETE FROM post_tags WHERE "post_tags"."tag_id" = ?`, id)
		if resr.Error != nil {
			c.sendReply(-1, "DB error:"+resr.Error.Error(), 0, nil)
			return
		}
		logs.Trace("DeleteTag() 删除了 post_tags 关联", resr.RowsAffected, "条")
		res := util.DB.Debug().Select("post_tags").Delete(tag, id)
		if res.Error != nil {
			c.sendReply(-1, "DB error:"+res.Error.Error(), 0, nil)
			return
		}
		c.sendReply(0, "success", int(res.RowsAffected), nil)
		return
	}
	c.sendReply(-1, "no permission", 0, nil)
	return
}
