package models

import (
	"github.com/beego/beego/v2/core/logs"
)

const Rules = `
# 版规

欢迎，A岛是一个前台匿名后台实名的，旨在为大家提供一个理性客观中立的讨论环境，A岛核心精神：开放包容，理性客观，有事说事，就事论事
请注意以下岛规和各版版规，**继续访问本站则视为理解并愿意遵守所有版规及[免责声明](https://adnmb3.com/t/11689471)**，对版规有疑问请参考和询问[版务集中串](https://adnmb3.com/t/17703121)：

- 以下违规内容将被采取碎槽清退等严重处理：1、涉政（敏感人物、标志、事件、民族、宗教、地区）、涉黄（露点、18X作品）、及其它违反中华人民共和国法律法规（毒赌嫖、犯罪、欺诈、造谣、邪教、魔法）的内容，在任何国内网站都必须被优先处理；2、跳脸群嘲岛民、恶意挑衅岛管理、屡次违规仍再犯的内容，恶意破坏本网站秩序的行为将被拒绝访问。
- 以下违规内容将被采取碎饼碎槽等一般处理：1、与A岛精神和社会公德相悖的内容，包括但不限于带妈及以上的嘴臭、屡次发表/明知某内容不适合在A岛公开讨论（见下条）但仍要发表的恶意引战钓鱼内容、频繁骚扰岛民正常讨论；2、与第一条有关的擦边球/红线内容，包括a) 上下任一部位比基尼级别的暴露、表面正常但有色情暗示的诱惑图文、明确有18X作品名称的任何正常讨论，b) 有政治元素的笑话、无官媒来源的社会时政消息、发布官媒消息时带个人片面解读、转进或者引发涉政相关的讨论（无论观点立场为何），c) 指路涉黄涉政内容至境外网站、侵犯正常岛民隐私、未明确虚构写作的自残自杀和残害他人的描述，d) 歪曲版规、带红名正常管理的节奏，排挤未违规的岛民；3、商业广告、推广链接、宣传QQ群等站外高风险内容。
- 以下非违规但不适合在A岛公开讨论的内容将被采取[SAGE](https://adnmb3.com/t/12366940)/删串碎饼等警告：1、沙发、队形、留名、刷版、新人签到、回复可见、无脑跟风、纯顶串水串等无意义的内容，为保证A岛讨论质量而不鼓励发表；2、晒妹、秀优越、纯伸手党、倒垃圾、地图炮等容易引起串中普遍敌视的内容，为保证A岛讨论氛围而禁止发表；3、特定话题未发往专版的串版内容，包括但不限于时事热点只发速报2，女性视角只发姐妹1，无公众讨论意义只发日记，圈子事（AB站等）和版务只发圈内。
- **本站为独立自主运营的私人论坛，使用本站是一种由本站授予您的特权而不是您固有的权利，管理员（红名）随时有权主观片面决定删除或保留内容、允许或禁止用户在本站发表内容。如果对此感到不适，可以选择不使用本站。如果你认为管理员可以做得更好，欢迎随时自愿[加入红名管理](https://adnmb3.com/t/14051281)团队。**

以下是一些系统设定：

- 后台注册账号并验证手机后便可以发表内容并解除99页浏览限制，您的个人信息不会以任何形式出现在前台。
- “饼干”是一个7位随机字符串，作为用户在前台发言时的唯一标识信息，每位注册用户在后台实名验证后即拥有2个饼干槽位，注册后如无违规行为（未被碎过槽）则每月增加1槽位，最多5槽位，每个饼干槽位均可单独领取、应用、导出或删除（销毁）饼干。客户端APP可扫描二维码导入饼干。
- 发文可以用>>或者>字符进行引用，点击No.编号也可以。发文有时间间隔，具体间隔请查看版块说明。
- 管理员（红名）认为不适合讨论但不直接严重违反版规的内容会被sage，被sage的串24小时后无法再回复。
- 可附加图片类型：GIF, JPG, PNG。
- 附加图片最大上限 2 M ( 2048K )。宽 250像素、高 250像素以上时自动缩小尺寸显示。
- 默认开启水印，可选择关闭。

来源： [首页 - A岛匿名版](https://adnmb3.com/Forum#ixzz6sadqhAy7) https://adnmb3.com/Forum#ixzz6sadqhAy7
`

func DataPrepare() {

	UserGroupRoot := NewUserGroup("root", "版主", "不服来干", 5, 0).Insert()
	UserGroupAdmin := NewUserGroup("admin", "管理员", "系统管理员，拥有管理所有板块与用户的权限", 4, 0).Insert()
	UserGroupAuthUser := NewUserGroup("authorizeduser", "授权用户", "授权用户，拥有管理所授权板块的权限", 3, 0).Insert()
	UserGroupUser := NewUserGroup("user", "普通用户", "普通用户，拥有基础的交互功能", 2, 0).Insert()
	UserGroupGuest := NewUserGroup("guest", "访客", "访客，允许浏览基础的网页", 1, 0).Insert()
	UserGroupNoAccess := NewUserGroup("noaccess", "禁止访问", "禁止访问，账户禁封中", -1, 0).Insert()
	UserGroupDefault := NewUserGroup("default", "default", "default", 0, 0).Insert()
	logs.Trace("[DBPrepare] NewUserGroup()", UserGroupRoot, UserGroupAdmin, UserGroupAuthUser, UserGroupUser, UserGroupGuest, UserGroupNoAccess, UserGroupDefault)

	UserRoot := NewUser("root", "root@user.com", "root", 32767, UserGroupRoot).Insert()
	UserAdmin := NewUser("admin", "admin@user.com", "admin", 1000, UserGroupAdmin).Insert()
	UserAuthUser := NewUser("auth", "auth@user.com", "auth", 100, UserGroupAuthUser).Insert()
	UserUser0 := NewUser("user0", "user0@user.com", "user", 10, UserGroupUser).Insert()
	UserUser1 := NewUser("user1", "user1@user.com", "user", 10, UserGroupUser).Insert()
	UserUser2 := NewUser("user2", "user2@user.com", "user", 10, UserGroupUser).Insert()
	UserUser3 := NewUser("user3", "user3@user.com", "user", 10, UserGroupUser).Insert()
	UserGuest := NewUser("guest", "guest@user.com", "guest", 10, UserGroupGuest).Insert()
	UserNoAccess := NewUser("noaccess", "nacc@user.com", "noaccess", -10, UserGroupNoAccess).Insert()
	UserDefault := NewUser("d", "default@user.com", "d", 0, UserGroupDefault).Insert()
	logs.Trace("[DBPrepare] NewUser()", UserRoot, UserAdmin, UserAuthUser, UserUser0, UserUser1, UserUser2, UserUser3, UserGuest, UserNoAccess, UserDefault)

	classMisc := NewClass(UserAdmin, "杂谈", "杂谈版，找不到分类的可以先放在这里", 0, UserDefault).Insert()
	classTimeline := NewClass(UserAdmin, "时间线", "综合版，集齐一手消息", 0, UserDefault).Insert()
	classYWH := NewClass(UserAdmin, "亚文化", "沵竾湜②佽沅嬤", 0, UserDefault).Insert()
	classTech := NewClass(UserAdmin, "科技", "科学技术是第一生产力 /n/t ————邓小平", 0, UserDefault).Insert()
	classMedia := NewClass(UserAdmin, "影视", "影视杂谈", 0, UserDefault).Insert()
	classAdmin := NewClass(UserAdmin, "值班室", "门卫大爷们的值班室", 0, UserDefault).Insert()
	classDark := NewClass(UserRoot, "deep♂dark♂fantasy", "影视杂谈", 3, UserGroupAuthUser).Insert()
	logs.Trace("[DBPrepare] NewClass()", classMisc, classTimeline, classYWH, classTech, classMedia, classAdmin, classDark)

	TribeMisc := NewTribe(false, "杂谈版", "#000000", "杂谈版1", UserAdmin, classMisc).Insert()
	TribeOld := NewTribe(false, "怀旧", "#F5DEB3", "杂谈版1", UserAdmin, classMisc).Insert()
	TribeNews := NewTribe(false, "新闻", "#000000", "新闻版", UserAdmin, classMisc).Insert()
	TribeQA := NewTribe(false, "问答", "#696969", "Q&A", UserAdmin, classMisc).Insert()
	TribeTimeLine := NewTribe(false, "时间线1", "#708090", "时间线1", UserAdmin, classTimeline).Insert()
	TribeComputer := NewTribe(false, "电脑", "#000000", "这就不用介绍了吧", UserAdmin, classTech).Insert()
	TribeACGN := NewTribe(false, "二次元", "#66CCFF", "莪湜②佽沅！", UserAdmin, classTech).Insert()
	TribeNet := NewTribe(false, "网络", "#000000", "Across the Great Wall we can reach every corner in the world.", UserAdmin, classTech).Insert()
	TribeReport := NewTribe(false, "举报", "#FFD700", "Q&A", UserAdmin, classMisc).Insert()
	TribeBoard := NewTribe(false, "公告", "#DC143C", "Q&A", UserAdmin, classMisc).Insert()
	logs.Trace("[DBPrepare] NewTribe()", TribeMisc, TribeOld, TribeNews, TribeQA, TribeTimeLine, TribeComputer, TribeACGN, TribeNet, TribeReport, TribeBoard)

	Red := NewTag("赤色", "#FF5722").Insert()
	Orange := NewTag("橙色", "#FFB800").Insert()
	Green := NewTag("墨绿", "#009688").Insert()
	Cyan := NewTag("藏青", "#2F4056").Insert()
	Blue := NewTag("蓝色", "#1E9FFF").Insert()
	Black := NewTag("雅黑", "#393D49").Insert()
	Gray := NewTag("银灰", "#F6F6F6").Insert()
	logs.Trace("[DBPrepare] CreateTag()", Red, Orange, Green, Green, Cyan, Blue, Black, Gray)

	PostRoles := NewPost(TribeBoard, UserRoot, "版规", Rules, "", []Tag{*NewTag("红色", "FF0000")}).Insert()
	logs.Trace("[DBPrepare] NewPost()", PostRoles)

	Reply0 := NewReply(UserAdmin, PostRoles, 0, "这是第一条评论").Insert()
	Reply0_0 := NewReply(UserAdmin, PostRoles, Reply0, "这是第一条子评论").Insert()
	Reply0_0_0 := NewReply(UserAuthUser, PostRoles, Reply0, "这是第一条子评论的子评论").Insert()
	Reply0_0_1 := NewReply(UserUser0, PostRoles, Reply0, "这是第一条子评论的子评论").Insert()
	logs.Trace("[DBPrepare] NewPost()", Reply0, Reply0_0, Reply0_0_0, Reply0_0_1)
}
