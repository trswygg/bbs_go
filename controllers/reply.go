package controllers

import (
	"bbs/models"
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
	"time"
)

type ReplyController struct {
	web.Controller
}
type ReplyControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Count     int
	Data      interface{}
}

func (c *ReplyController) URLMapping() {
	logs.Trace("[URLMapping] mapping ReplyController")
	c.Mapping("reply", c.GetReplyInfo)
	c.Mapping("reply/update", c.ChReply)
	c.Mapping("reply/create", c.NewReply)
	c.Mapping("reply/delete", c.DropReply)

	c.Mapping("favorite", c.DoFavorite)
	c.Mapping("like", c.DoLike)
	c.Mapping("unfavorite", c.UnFavorite)
	c.Mapping("unlike", c.UnLike)
}

func (c *ReplyController) sendReply(code int, result string, count int, data interface{}) {
	c.Data["json"] = ReplyControllerReply{
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

// NewReply 创建post的reply
// method post
// param post_id parent_id content
func (c *ReplyController) NewReply() {
	//    CreatorId int       `gorm:"comment:创建人;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;index"`
	//    PostId    int       `gorm:"comment:帖子;index"`
	//    Content   string    `gorm:"size:50;not null;comment:内容(MD)"`
	//    ParentId  int       `gorm:"commit:父节点;null"`
	creatorId, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	postId, _ := c.GetInt("post_id", 0)
	parentId, _ := c.GetInt("parent_id", 0)
	content := c.GetString("content", "NaN")
	r := models.NewReply(creatorId, postId, parentId, content)
	var res *gorm.DB
	// create
	res = util.DB.Table("reply").Create(&r)
	logs.Info("create reply %#v", r)
	if res.Error != nil {
		c.sendReply(-1, "error on create reply :"+res.Error.Error(), 0, nil)
		return
	} else {
		c.sendReply(0, "success", int(res.RowsAffected), r)
		return
	}
}

// ChReply 修改reply
// method post
// router /reply_info [post]
func (c *ReplyController) ChReply() {}

// DropReply 删除reply
// method post
// router /reply_info [post]
func (c *ReplyController) DropReply() {}

// GetReplyInfo 获取reply
// method get
// router /reply:id [get]
func (c *ReplyController) GetReplyInfo() {
	reply := new(models.Reply)
	requireId := c.GetString("id")
	// get user
	rows := util.DB.Model(&models.Reply{}).Find(reply, requireId).RowsAffected
	if rows == 1 {
		c.sendReply(0, "success", 1, reply)
		return
	} else {
		c.sendReply(-1, "no such tag id:"+requireId, 0, nil)
		return
	}
}

// DoFavorite 收藏
// method post
// router /favorite:post_id [post]
func (c *ReplyController) DoFavorite() {
	favorites := new(models.Favorites)
	post := new(models.Post)
	requirePostId, _ := c.GetInt("post_id")
	// get user
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	sub := util.DB.First(post, requirePostId) // postTitle
	var countHave int64
	mainQ := util.DB.Model(favorites).Where("post_id = ? AND creator_id = ?", requirePostId, uid).Count(&countHave)
	if mainQ.Error != nil {
		c.sendReply(-1, "sql error"+mainQ.Error.Error(), 0, nil)
		return
	}
	if countHave != 0 {
		c.sendReply(-1, "already", 0, nil)
		return
	}
	if sub.Error != nil {
		c.sendReply(-1, "err:"+sub.Error.Error(), 0, nil)
		return
	}
	favorites.CreatorId = uid
	favorites.PostId = requirePostId
	favorites.PostTitle = post.Title
	result := util.DB.Create(favorites)
	if result.Error != nil {
		c.sendReply(-1, "err:"+result.Error.Error(), 0, nil)
		return
	} else {
		c.sendReply(0, "success", int(result.RowsAffected), favorites)
		return
	}
}

// DoLike 喜欢
// method post
// router /favorite:post_id [post]
func (c *ReplyController) DoLike() {
	likes := new(models.Likes)
	requirePostId, _ := c.GetInt("post_id")
	// get user
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	var countHave int64
	mainQ := util.DB.Model(likes).Where("post_id = ? AND creator_id = ?", requirePostId, uid).Count(&countHave)
	if mainQ.Error != nil {
		c.sendReply(-1, "sql error"+mainQ.Error.Error(), 0, nil)
		return
	}
	if countHave != 0 {
		c.sendReply(-1, "already", 0, nil)
		return
	}
	likes.CreatorId = uid
	likes.PostId = requirePostId
	result := util.DB.Create(likes)
	if result.Error != nil {
		c.sendReply(-1, "err:"+result.Error.Error(), 0, nil)
		return
	} else {
		c.sendReply(0, "success", int(result.RowsAffected), likes)
		return
	}
}

// UnFavorite 取消收藏
// method post
// router /unfavorite:post_id [post]
func (c *ReplyController) UnFavorite() {
	favorite := new(models.Favorites)
	requirePostId, _ := c.GetInt("post_id")
	// get user
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	res := util.DB.Where("post_id = ?", requirePostId).Where("creator_id = ?", uid).Delete(favorite)
	if res.Error != nil {
		c.sendReply(-1, "err:"+res.Error.Error(), 0, nil)
		return
	} else {
		c.sendReply(0, "success", int(res.RowsAffected), favorite)
		return
	}
}

// UnLike 取消喜欢
// method post
// router /unlike:post_id [post]
func (c *ReplyController) UnLike() {
	likes := new(models.Likes)
	requirePostId, _ := c.GetInt("post_id")
	// get user
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	res := util.DB.Where("post_id = ?", requirePostId).Where("creator_id = ?", uid).Delete(likes)
	if res.Error != nil {
		c.sendReply(-1, "err:"+res.Error.Error(), 0, nil)
		return
	} else {
		c.sendReply(0, "success", int(res.RowsAffected), likes)
		return
	}
}
