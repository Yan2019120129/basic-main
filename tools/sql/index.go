package sql

import (
	"basic/tools/utils"
)

// InitializeTables 初始化数据表
func InitializeTables() []*utils.InitTable {
	tables := []*utils.InitTable{
		BasicAdminRBACItem,           //	权限表
		BasicAdminRBACChild,          //	权限表
		BasicAdminRBACAssignment,     //	权限表
		BasicAccessLogs,              //	访问日志
		BasicAdminMenu,               //	管理菜单
		BasicAdminSetting,            //	管理设置
		BasicAdminUser,               //	后台管理
		BasicHomeUser,                //	前端用户
		BasicUserSetting,             //	用户配置
		BasicHomeLang,                //	用户语言
		BasicHomeLangDictionary,      //	语言字典
		BasicHomeUserInvite,          //	用户邀请码
		BasicHomeCountry,             //	国家列表
		BasicHomeUserLevel,           //	等级列表
		BasicHomeUserLevelOrder,      //	等级订单
		BasicHomeUserBill,            //	用户账单
		BasicHomeWalletPayment,       //	钱包支付
		BasicHomeUserWalletAccount,   //	钱包账户
		BasicHomeUserWalletOrder,     //	钱包订单
		BasicHomeUserVerify,          //	用户认证
		BasicHomeAssets,              //	资产列表
		BasicHomeUserAssets,          //	用户资产
		BasicHomeProductCategory,     //	产品分类
		BasicHomeProduct,             //	产品
		BasicHomeProductOrder,        //	产品订单
		BasicChatConversation,        //	聊天会话
		BasicChatConversationMessage, //	聊天会话消息

		BasicFinancialStatistics,   //	商城财务统计表
		BasicUserComment,           //	用户评论表
		BasicCommodity,             //	商城商品表
		BasicStore,                 //	商城商店表
		BasicCommodityOrder,        //	商城商品订单信息表
		BasicShopMyWallet,          //	商城用户钱包记录
		BasicShopStatisticalRecord, //	商城店铺统计记录表
		BasicShopUser,              //	商城用户表
		BasicUserAddress,           //	收货地址表
		BasicUserAttention,         //	用户关注表
		BasicUserFavorites,         //	用户收藏表
		BasicCategory,              //	商城类目表
	}

	// 项目表
	tables = append(tables, []*utils.InitTable{}...)

	return tables
}

// InitializeAuth 初始化RBAC
func InitializeAuth() *utils.Permission {
	return &utils.Permission{
		Roles: []string{"默认后台管理员", "默认组长管理员"},
		RoleOnlyRouter: map[string][]string{
			"默认组长管理员": {"上传文件", "首页信息", "管理信息", "更新密码", "更新信息", "管理Websocket连接", "会话列表", "会话信息", "发送消息", "清除未读", "会话消息", "操作日志", "用户管理", "用户新增", "用户更新", "用户删除", "用户关系", "用户资产", "资产更新", "资产删除", "资产新增", "用户账单列表", "用户账单类型", "用户认证列表", "用户认证审核", "用户提现账户列表", "用户提现账户更新", "用户充值列表", "用户充值审核", "用户提现列表", "用户提现审核", "产品订单列表", "产品订单更新", "国家Options", "平台资产数组", "产品分类目录", "支付提现类型", "产品列表选项"},
		},
		RoleFilterRouter: map[string][]string{
			"默认后台管理员": {"数据库表", "数据表信息", "权限数组", "角色列表", "角色更新", "角色新增", "角色删除", "菜单列表", "菜单更新"},
		},
	}
}
