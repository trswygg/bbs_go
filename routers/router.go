package routers

import (
	"bbs/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// api v2
	ns := web.NewNamespace("/api",

		// 	/api/misc
		web.NSNamespace("/misc",
			web.NSRouter("/captcha", &controllers.CaptchaController{}, "get:Get"),
			web.NSRouter("/captcha", &controllers.CaptchaController{}, "post:Post"),
			//web.NSInclude(&controllers.CaptchaController{}),
		),
		web.NSNamespace("/image",
			web.NSRouter("/upload", &controllers.ImageController{}, "post:Upload"),
			web.NSRouter("/list", &controllers.ImageController{}, "get:GetLists"),
			web.NSRouter("/delete", &controllers.ImageController{}, "delete:Delete"),
			//web.NSInclude(&controllers.CaptchaController{}),
		),
		// /api/user
		web.NSNamespace("/user",
			web.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			web.NSRouter("/logout", &controllers.UserController{}, "post:Logout"),
			web.NSRouter("/info", &controllers.UserController{}, "get:GetUser"),
			web.NSRouter("/info", &controllers.UserController{}, "post:UpdateUser"),
			web.NSRouter("/reg", &controllers.UserController{}, "post:RegUser"),
			web.NSRouter("/profile", &controllers.UserController{}, "get:GetProfile"),
			web.NSRouter("/profile", &controllers.UserController{}, "post:UpdateProfile"),
		),
		//api/bbs
		web.NSNamespace("/bbs",
			web.NSRouter("/class", &controllers.ClassController{}, "get:GetClassInfo"),
			web.NSRouter("/class/list", &controllers.ClassController{}, "post:QueryClass"),
			web.NSRouter("/class/update", &controllers.ClassController{}, "post:UpdateClass"),
			web.NSRouter("/class/create", &controllers.ClassController{}, "post:CreateClass"),
			web.NSRouter("/class/by_user", &controllers.ClassController{}, "get:QueryClassByUser"),
			web.NSRouter("/class/subscribe", &controllers.ClassController{}, "get:SubscribeClass"),

			web.NSRouter("/tag", &controllers.TagController{}, "get:GetTagInfo"),
			web.NSRouter("/tag/query", &controllers.TagController{}, "get:QueryTag"),
			web.NSRouter("/tag/create", &controllers.TagController{}, "post:CreateTag"),
			web.NSRouter("/tag/delete", &controllers.TagController{}, "delete:DeleteTag"),

			web.NSRouter("/tribe", &controllers.TribeController{}, "get:GetTribeInfo"),
			web.NSRouter("/tribe/list", &controllers.TribeController{}, "get:QueryTribes"),
			web.NSRouter("/tribe/by_user", &controllers.TribeController{}, "get:GetTribesByUser"),
			web.NSRouter("/tribe/subscribe", &controllers.TribeController{}, "post:SubscribeTribe"),
			web.NSRouter("/tribe/update", &controllers.TribeController{}, "post:UpdateTribe"),
			web.NSRouter("/tribe/create", &controllers.TribeController{}, "post:CreateTribe"),
			web.NSRouter("/tribe/delete", &controllers.TribeController{}, "post:DeleteTribe"),

			web.NSRouter("/post", &controllers.PostController{}, "get:GetPostInfo"),
			web.NSRouter("/post/query", &controllers.PostController{}, "get:QueryPosts"),
			web.NSRouter("/post/timeline", &controllers.PostController{}, "get:PostTimeLine"),
			web.NSRouter("/post/create", &controllers.PostController{}, "post:NewPost"),
			web.NSRouter("/post/delete", &controllers.PostController{}, "delete:DeletePost"),

			web.NSRouter("/reply", &controllers.ReplyController{}, "get:GetReplyInfo"),
			web.NSRouter("/reply/update", &controllers.ReplyController{}, "post:ChReply"),
			web.NSRouter("/reply/create", &controllers.ReplyController{}, "post:NewReply"),
			web.NSRouter("/reply/delete", &controllers.ReplyController{}, "delete:DropReply"),

			web.NSRouter("/favorite", &controllers.ReplyController{}, "post:DoFavorite"),
			web.NSRouter("/like", &controllers.ReplyController{}, "post:DoLike"),
			web.NSRouter("/unfavorite", &controllers.ReplyController{}, "post:UnFavorite"),
			web.NSRouter("/unlike", &controllers.ReplyController{}, "post:UnLike"),

			//web.NSInclude(&controllers.ClassController{}),
		),
	)
	web.AddNamespace(ns)
}

//func init() {
//	// 登陆路由
//	beego.Router("/user/login/", &controllers.UserLoginController{})
//	// 注销路由
//	beego.Router("/user/logout/", &controllers.UserLogoutController{})
//	// 注册路由
//	beego.Router("/user/register/", &controllers.UserRegisterController{})
//}
