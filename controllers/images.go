package controllers

import (
	"bbs/models"
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"time"
)

type ImageController struct {
	beego.Controller
}

func (c *ImageController) URLMapping() {
	logs.Trace("[URLMapping] mapping UserController")
	c.Mapping("upload", c.Upload)
	c.Mapping("list", c.GetLists) // get
	c.Mapping("delete", c.Delete) // post
}

type ImageControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Count     int
	Data      interface{}
}

func (c *ImageController) sendReply(code int, result string, count int, data interface{}) {
	c.Data["json"] = ImageControllerReply{
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

// Upload 上传图片
// method post
func (c *ImageController) Upload() {
	image := new(models.Images)
	// get user
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	f, h, err := c.GetFile("image") //获取上传的文件
	if err != nil {
		c.sendReply(-1, "error on getting file <image> :"+err.Error(), 0, nil)
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.sendReply(-2, "wrong file type", 0, nil)
		return
	}
	//创建目录
	wd, errwd := os.Getwd()
	if errwd != nil {
		c.sendReply(-2, "os error: "+errwd.Error(), 0, nil)
		return
	}
	uploadDir := wd + "/static/images/" + time.Now().Format("20060102")
	err = os.MkdirAll(uploadDir, 666)
	if err != nil {
		c.sendReply(-2, "server error:"+err.Error(), 0, nil)
		return
	}
	//构造文件名称
	rand.Seed(time.Now().UnixNano())
	randNum := util.GeneId()
	fileName := randNum + ext
	//this.Ctx.WriteString(  fileName )
	fpath := uploadDir + fileName
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
		}
	}(f) //关闭上传的文件，不然的话会出现临时文件不能清除的情况

	// insert
	image.CreatorId = uid
	image.ImageName = fileName
	res := util.DB.Create(image)
	if res.Error != nil {
		c.sendReply(-4, "DB error:"+res.Error.Error(), 0, nil)
		return
	}
	err = c.SaveToFile("image", fpath)
	logs.Trace("c.SaveToFile('image' " + fpath + ")")
	if err != nil {
		c.sendReply(-3, "server error:"+err.Error(), 0, nil)
		return
	}
	c.sendReply(0, "success", 1, image)
}

// GetLists 获取图片
// method get
// param id
// router /image:page:limit 【get】
func (c ImageController) GetLists() {
	var imageList = make([]models.Images, 0)
	var user = new(models.User)
	id, _ := c.GetInt("id", 0)
	requireLimit, errl := c.GetInt("limit", 10)
	requireOffset, erri := c.GetInt("offset", 0)
	if errl != nil {
		c.sendReply(-1, "wrong param: limit", 0, nil)
		return
	} else if erri != nil {
		c.sendReply(-1, "wrong param index", 0, nil)
		return
	}
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	util.DB.Preload("Profile").First(user, id)
	if uid != id && user.Profile.ImageVisibility == 0 {
		logs.Debug("uid:%d,id:%d", uid, id)
		logs.Debug("%v", user.Profile)
		c.sendReply(-1, "nothing here", 0, nil)
	}
	res := util.DB.Where(" creator_id = ?", id).Limit(requireLimit).Offset(requireOffset).Find(&imageList)
	for i, _ := range imageList {
		imageList[i].CreatorName = user.Name
		imageList[i].CreatorFace = user.Face
	}
	if res.Error != nil {
		c.sendReply(-1, "DB error:"+res.Error.Error(), 0, nil)
		return
	}
	c.sendReply(0, "success:", len(imageList), &imageList)
	return
}

// Delete 删除图片
// method post
// router /image/drop
func (c ImageController) Delete() {
	var images models.Images
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	id, err := c.GetInt("id")
	if err != nil || id == 0 {
		c.sendReply(-1, "wrong param id", 0, nil)
		return
	}
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	res := util.DB.Debug().Where("creator_id = ?", uid).Delete(&images, id)
	if res.Error != nil {
		c.sendReply(-1, "DB error:"+res.Error.Error(), 0, nil)
		return
	}
	// Notice RowsAffected 为 0 依旧返回成功
	c.sendReply(-1, "success", int(res.RowsAffected), nil)
	return
}
