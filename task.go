package main

import (
	"bbs/controllers"
	"bbs/models"
	"bbs/util"
	"fmt"
	"github.com/beego/beego/v2/adapter/toolbox"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

func init() {
	initTaskShowDB()
	//initTaskAutoClean()
}

func initTaskShowDB() {
	tk := toolbox.NewTask("gorm.ShowDB", "0 0 0/12 * * *", func() error {
		sqlDB, errDB := util.DB.DB()
		if errDB != nil {
			logs.Error("[gorm] ShowDB() :", errDB)
		}
		stats := sqlDB.Stats()
		fmt.Printf("[gorm] ShowDB() %+v \n", stats)
		return errDB
	})
	toolbox.AddTask("gorm.ShowDB", tk)
}

// initTaskAutoClean 自动清理 per7day
func initTaskAutoClean() {
	tk := toolbox.NewTask("bbs.autoClean", "0 0 4 * * 1", func() error {
		logs.Info("[task] start bbs.AutoClean")
		var tribes []models.Tribe
		util.DB.Unscoped().Where(" activity_degree = 0 AND protect = false").Find(&tribes)
		now := time.Now()
		for i := range tribes {
			// Delete is null (0001-01-01 00:00:00+00)
			if tribes[i].DeletedAt.IsZero() {
				controller := &controllers.MailBoxController{}
				controller.SendMailFromSys(geneMsgHide(tribes[i]), 2, tribes[i].CreatorId)
				logs.Info("[task.AutoClean] set DeleteAt on %#v", tribes[i])
				db := util.DB.Delete(&tribes[i])
				if db.Error != nil {
					return db.Error
				}
			} else if now.Sub(tribes[i].DeletedAt).Hours() > 7*24 {
				controller := &controllers.MailBoxController{}
				controller.SendMailFromSys(geneMsgDelete(tribes[i]), 2, tribes[i].CreatorId)
				db := util.DB.Unscoped().Delete(&tribes[i]) // Drop
				logs.Info("[task.AutoClean] DELETE %#v", tribes[i])
				if db.Error != nil {
					return db.Error
				}
			}
		}
		return nil
	})
	toolbox.AddTask("gorm.ShowDB", tk)
}

func geneMsgHide(tribe models.Tribe) string {
	return ""
}
func geneMsgDelete(tribe models.Tribe) string {
	return ""
}
