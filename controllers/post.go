package controllers

import (
	"bbs/models"
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
	"time"
)

type PostController struct {
	web.Controller
}
type PostControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Count     int
	Data      interface{}
}

const queryRAW = `SELECT
	"id",
	title,
	created_at,
	creator_id,
	tribe_id
FROM
	post 
WHERE
	to_tsvector( 'jiebaqry', "content" ) @@ websearch_to_tsquery( 'jiebacfg', ? )
ORDER BY id DESC
LIMIT ? OFFSET ?
`

func (c *PostController) URLMapping() {
	c.Mapping("post", c.GetPostInfo)
	c.Mapping("post/query", c.QueryPosts)
	c.Mapping("post/timeline", c.PostTimeLine)
	c.Mapping("post/create", c.NewPost)
	c.Mapping("post/delete", c.DeletePost)

}

func (c *PostController) sendReply(code int, result string, count int, data interface{}) {
	c.Data["json"] = PostControllerReply{
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

//oO0OoO0OoO0Oo Post oO0OoO0OoO0Oo

// QueryPosts 获取QueryPost列表
// method get
// param :text:text_type:limit:offset [post]
func (c *PostController) QueryPosts() {
	post := new(models.Post)
	resultArr := make([]models.PostQuery, 10)
	requireText := c.GetString("text")
	requireType := c.GetString("text_type")
	requireLimit, errl := c.GetInt("limit", 3)
	requireOffset, erri := c.GetInt("offset", 0)
	if errl != nil {
		c.sendReply(-1, "wrong param: limit", 0, nil)
		return
	} else if erri != nil {
		c.sendReply(-1, "wrong param index", 0, nil)
		return
	}
	switch requireType {
	case "title":
		var total int64
		db := util.DB.Model(post).Order("id").
			Where("title LIKE ?", "%"+requireText+"%").
			Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
		if db.Error != nil {
			c.sendReply(0, "sql fail:"+db.Error.Error(), int(0), resultArr)
			return
		}
		c.sendReply(0, "success", int(total), resultArr)
		return
	case "content":
		var total int64
		db := util.DB.Raw(queryRAW, requireText, requireLimit, requireOffset).Find(&resultArr)
		total = db.RowsAffected
		if db.Error != nil {
			c.sendReply(0, "sql fail:"+db.Error.Error(), int(-1), resultArr)
			return
		}
		c.sendReply(0, "success", int(total), resultArr)
	default:
		c.sendReply(-2, "no text_type support", 0, nil)
		return

	}

}

// PostTimeLine 时间线-综合
// @param type [all,user,class,tribe,rand,tribes,  follow]
//              全部,某人,分类id,部落id,随机, 部落-user,关注(auth)
func (c *PostController) PostTimeLine() {
	post := new(models.Post)
	resultArr := make([]models.PostQuery, 0)
	requireType := c.GetString("type", "all")
	requireId := c.GetString("id", "1")
	requireLimit, errl := c.GetInt("limit", 10)
	requireOffset, erri := c.GetInt("offset", 0)
	logs.Debug("(c *PostController) PostTimeLine() type=%v,id=%v,limit=%v,offset=%v cookie:AccessKey=%v",
		requireType, requireId, requireLimit, requireOffset, c.GetString("AccessKey"))
	if errl != nil {
		c.sendReply(-1, "wrong param: limit", 0, nil)
		return
	} else if erri != nil {
		c.sendReply(-1, "wrong param index", 0, nil)
		return
	}
	var total int64
	var db *gorm.DB
	// requireType [all,user,class,tribe,rand,tribes,   follow]
	switch requireType {
	case "all":
		db = util.DB.Model(post).Order("id DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	case "user": //  CreatorId
		db = util.DB.Model(post).Where("creator_id = ?", requireId).Order("id DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	case "class": // TribeId in ?
		subSql := util.DB.Model(models.Tribe{}).Select("id").Where("class_id = ?", requireId)
		db = util.DB.Model(post).Where("tribe_id IN (?)", subSql).Order("created_at DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	case "tribe":
		db = util.DB.Model(post).Where("tribe_id = ?", requireId).Order("id DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	case "rand": // order by random()
		db = util.DB.Model(post).Order("random()").Limit(requireLimit).Find(&resultArr)
		total = db.RowsAffected
	case "tribes": //  自己的部落
		userId, b := models.GetUserIdByKey(c.GetString("AccessKey"))
		if !b {
			c.sendReply(-1, "you should login first", 0, nil)
			return
		}
		subSql := util.DB.Model(models.Tribe{}).Select(" id").Where("creator_id = ?", userId)
		db = util.DB.Model(post).Where("tribe_id IN (?)", subSql).Order("created_at DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	case "follow": //关注的人
		userId, b := models.GetUserIdByKey(c.GetString("AccessKey"))
		if !b {
			c.sendReply(-1, "you should login first", 0, nil)
			return
		}
		subSql := util.DB.Model(models.Follows{}).Select("user_id").Where("creator_id = ?", userId)
		db = util.DB.Model(post).Where("creator_id IN (?)", subSql).Order("created_at DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	default:
		db = util.DB.Model(post).Order("id DESC").Limit(requireLimit).Offset(requireOffset).Find(&resultArr)
		total = db.RowsAffected
	}
	if db.Error != nil {
		c.sendReply(-1, "sql fail:"+db.Error.Error(), 0, resultArr)
		return
	}
	c.sendReply(0, "success", int(total), resultArr)
	return
}

// GetPostInfo 获取post详细信息 登录查看正文
// method get
// @param bool content
// TODO 访问加分
func (c *PostController) GetPostInfo() {
	post := new(models.Post)
	user := new(models.User)
	tribe := new(models.Tribe)
	class := new(models.Class)
	requireId := c.GetString("id")
	bContent, _ := c.GetBool("content", false)
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		bContent = false
	}
	// get user
	rows := util.DB.Preload("Tags").Preload("Replys").Find(post, requireId).RowsAffected
	if rows == 1 {
		var fCount, lCount, rCount int64
		var countBLike, countBFavorite int64
		util.DB.Model(&models.Favorites{}).Where("post_id = ?", requireId).Count(&fCount)
		util.DB.Model(&models.Likes{}).Where("post_id = ?", requireId).Count(&lCount)
		util.DB.Model(&models.Reply{}).Where("post_id = ?", requireId).Count(&rCount)
		util.DB.Select("name", "face").Find(user, post.CreatorId)
		util.DB.Select("name", "class_id").Find(tribe, post.TribeId)
		util.DB.Select("name").Find(class, tribe.ClassId)
		util.DB.Model(&models.Likes{}).Where("post_id = ? AND creator_id = ?", requireId, uid).
			Count(&countBLike)
		util.DB.Model(&models.Favorites{}).Where("post_id = ? AND creator_id = ?", requireId, uid).
			Count(&countBFavorite)
		post.FavoritesCount = int(fCount)
		post.LikesCount = int(lCount)
		post.ReplyCount = int(rCount)
		post.CreatorName = user.Name
		post.CreatorFace = user.Face
		post.TribeName = tribe.Name
		post.ClassName = class.Name
		if countBLike > 0 {
			post.AlreadyLike = true
		}
		if countBFavorite > 0 {
			post.AlreadyFavorite = true
		}
		if bContent == false {
			post.Content = "登录才能查看"
		}
		c.sendReply(0, "success", 1, post)
		return
	} else {
		c.sendReply(-1, "no such tag id:"+requireId, 0, nil)
		return
	}
}

// NewPost 创建post
// method post
// TODO 创建加分
func (c *PostController) NewPost() {
	tribeId, _ := c.GetInt("tribe_id")
	creatorId, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	title := c.GetString("title")
	content := c.GetString("content")
	image := c.GetString("main_image", "bg.jpg")
	tagIds := c.GetString("tag_ids")
	var tags []models.Tag
	post := models.NewPost(tribeId, creatorId, title, content, image, nil)
	ids, _ := util.StrToIntArray(tagIds)
	for _, n := range ids {
		tags = append(tags, models.Tag{Id: n})
	}
	post.Tags = tags
	res := util.DB.Create(&post)
	if res.Error != nil {
		c.sendReply(-1, "could not create post:"+res.Error.Error(), 0, nil)
		return
	}
	//errsql := util.DB.Model(&post).Association("Tags").Replace(&tags)
	//if errsql != nil {
	//	util.DB.Delete(&post) // rollback
	//	c.sendReply(-1, "could not insert Tags:"+errsql.Error(), 0, nil)
	//	return
	//}
	post.Content = ""
	c.sendReply(0, "success", 1, post)
	return
}

// DeletePost 删除帖子
// TODO 扣分
func (c *PostController) DeletePost() {
	post := new(models.Post)
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	id, err := c.GetInt("id")
	if err != nil {
		c.sendReply(-1, "wrong param id:"+err.Error(), 0, nil)
		return
	}
	util.DB.First(post, id)
	if post.CreatorId != uid {
		c.sendReply(-1, "wrong access", 0, nil)
		return
	}
	util.DB.Debug().Select("Tags", "Replys", "Favorites", "Likes").Delete(post, id)
}

type MailBoxController struct {
	web.Controller
}

func (c *MailBoxController) SendMailFromSys(msg string, from, to int) {

}
