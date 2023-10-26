package controllers

import (
	"basic/home/controllers/index"
	"basic/home/controllers/product"
	"basic/home/controllers/user"
	"basic/home/controllers/wallet"

	"github.com/so68/zfeng/router"
)

func Router() []*router.Handle {
	routerList := make([]*router.Handle, 0)

	// 添加默认路由
	routerList = append(routerList, []*router.Handle{
		//	基础方法
		router.NewTokenHandle("上传图片", "POST", "/upload", index.Upload),
		router.NewHandle("首页信息", "POST", "/index", index.Index),
		router.NewHandle("生成验证码", "GET", "/captcha/generate", index.GenerateCaptcha),
		router.NewHandle("显示验证码", "GET", "/captcha/image", index.ImageCaptcha),
		router.NewHandle("预处理数据", "POST", "/prefetch", index.PreFetch),
		router.NewHandle("语言包数据", "POST", "/locales", index.Locales),
		router.NewHandle("用户注册", "POST", "/register", index.Register),
		router.NewHandle("用户登陆", "POST", "/login", index.Login),
		router.NewHandle("配置内容", "POST", "/article", index.Article),
		router.NewHandle("下载文件", "POST", "/down", index.Download),

		// 用户方法
		router.NewTokenHandle("帮助中心", "POST", "/user/helpers", user.Helpers),
		router.NewTokenHandle("用户信息", "POST", "/user/info", user.Info),
		router.NewTokenHandle("用户账单", "POST", "/user/bill", user.Bill),
		router.NewTokenHandle("更新用户信息", "POST", "/user/update", user.Update),
		router.NewTokenHandle("更新用户登陆密码", "POST", "/user/update/password", user.UpdatePassword),
		router.NewTokenHandle("更新用户安全密钥", "POST", "/user/update/security", user.UpdateSecurity),
		router.NewTokenHandle("用户验证信息", "POST", "/user/verify/info", user.VerifyInfo),
		router.NewTokenHandle("用户验证", "POST", "/user/verify", user.Verify),
		router.NewTokenHandle("用户团队", "POST", "/user/team", user.Team),
		router.NewTokenHandle("用户等级列表", "POST", "/user/level/index", user.Level),
		router.NewTokenHandle("用户购买等级", "POST", "/user/level", user.LevelOrder),
		router.NewTokenHandle("用户邀请信息", "POST", "/user/invite", user.Invite),

		// 钱包操作
		router.NewTokenHandle("钱包订单列表", "POST", "/wallet/index", wallet.Index),
		router.NewTokenHandle("钱包充值信息", "POST", "/wallet/deposit/info", wallet.DepositInfo),
		router.NewTokenHandle("钱包充值创建", "POST", "/wallet/deposit", wallet.Deposit),
		router.NewTokenHandle("钱包提现创建", "POST", "/wallet/withdraw", wallet.Withdraw),
		router.NewTokenHandle("钱包提现信息", "POST", "/wallet/withdraw/info", wallet.WithdrawInfo),
		router.NewTokenHandle("钱包账户绑定", "POST", "/wallet/account", wallet.Account),
		router.NewTokenHandle("钱包账户列表", "POST", "/wallet/account/index", wallet.AccountIndex),
		router.NewTokenHandle("钱包账户更新", "POST", "/wallet/account/update", wallet.AccountUpdate),
		router.NewTokenHandle("钱包账户删除", "POST", "/wallet/account/delete", wallet.AccountDelete),

		// 产品方法
		router.NewTokenHandle("产品商品列表", "POST", "/product/index", product.Index),
		router.NewTokenHandle("产品商品购买", "POST", "/product/create", product.Create),
		router.NewTokenHandle("产品订单列表", "POST", "/product/order", product.Order),
		router.NewTokenHandle("产品商品详情", "POST", "/product/details", product.Details),
	}...)

	return routerList
}
