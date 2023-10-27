package sql

import (
	"basic/tools/utils"
)

const AdminMenuTableName = "admin_menu"
const AdminMenuTableComment = "管理菜单"
const CreateAdminMenu = `CREATE TABLE ` + AdminMenuTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '菜单名称',
	parent INT NOT NULL DEFAULT 0 COMMENT '父级ID',
	route VARCHAR(50) NOT NULL DEFAULT '' COMMENT '路由',
	sort TINYINT NOT NULL DEFAULT 0 COMMENT '排序',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -1禁用 10启用',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AdminMenuTableComment + `';`

const InsertAdminMenu = `INSERT INTO ` + AdminMenuTableName + ` VALUES
(1, '管理设置', 0, '', 100, 10, '{"icon": "sym_o_manage_accounts"}'),
(2, '管理列表', 1, '/manage/index', 0, 10, '{"icon": "sym_o_manage_accounts", "temp": "/manage/index"}'),
(3, '角色管理', 1, '/role/index', 1, 10, '{"icon": "sym_o_6_ft_apart", "temp": "/manage/role"}'),
(4, '菜单管理', 1, '/menu/index', 2, 10, '{"icon": "sym_o_list", "temp": "/manage/menu"}'),
(5, '操作日志', 1, '/logs/index', 3, 10, '{"icon": "sym_o_description", "temp": "/manage/logs"}'),

(10, '配置管理', 0, '', 101, 10, '{"icon": "sym_o_settings"}'),
(11, '支付设置', 10, '/payment/index', 0, 10, '{"icon": "sym_o_credit_card", "temp": "/setting/payment"}'),
(12, '资产设置', 10, '/assets/index', 1, 10, '{"icon": "sym_o_account_balance", "temp": "/setting/assets"}'),
(13, '等级设置', 10, '/level/index', 1, 10, '{"icon": "sym_o_brightness_auto", "temp": "/setting/level"}'),
(14, '国家设置', 10, '/country/index', 1, 10, '{"icon": "sym_o_flag_circle", "temp": "/setting/country"}'),
(15, '语言设置', 10, '/lang/index', 1, 10, '{"icon": "sym_o_language", "temp": "/setting/language"}'),
(16, '语言字典', 10, '/dictionary/index', 1, 10, '{"icon": "sym_o_g_translate", "temp": "/setting/dictionary"}'),
(17, '配置设置', 10, '/setting/index', 2, 10, '{"icon": "sym_o_history_edu", "temp": "/setting/index"}'),

(20, '用户管理', 0, '', 0, 10, '{"icon": "sym_o_engineering"}'),
(21, '用户列表', 20, '/user/index', 0, 10, '{"icon": "sym_o_people", "temp": "/user/index"}'),
(22, '用户资产', 20, '/user/assets/index', 1, 10, '{"icon": "sym_o_account_balance", "temp": "/user/assets"}'),
(23, '用户关系', 20, '/user/relation', 2, 10, '{"icon": "sym_o_diversity_1", "temp": "/user/relation"}'),
(24, '用户账户', 20, '/wallet/account/index', 3, 10, '{"icon": "sym_o_credit_card", "temp": "/user/wallet/account"}'),
(25, '用户认证', 20, '/user/verify/index', 4, 10, '{"icon": "sym_o_badge", "temp": "/user/verify"}'),
(30, '财务中心', 0, '', 1, 10, '{"icon": "sym_o_account_balance"}'),
(31, '用户账单', 30, '/user/bill/index', 0, 10, '{"icon": "sym_o_request_quote", "temp": "/user/bill"}'),
(32, '充值订单', 30, '/wallet/deposit/index', 1, 10, '{"icon": "sym_o_format_textdirection_r_to_l", "temp": "/user/wallet/deposit"}'),
(33, '提现订单', 30, '/wallet/withdraw/index', 2, 10, '{"icon": "sym_o_format_textdirection_l_to_r", "temp": "/user/wallet/withdraw"}'),
(40, '产品中心', 0, '', 2, 10, '{"icon": "sym_o_storefront"}'),
(41, '产品分类', 40, '/product/category/index', 0, 10, '{"icon": "sym_o_category", "temp": "/product/category"}'),
(42, '产品列表', 40, '/product/index/index', 0, 10, '{"icon": "sym_o_store", "temp": "/product/index"}'),
(43, '产品订单', 40, '/product/order/index', 0, 10, '{"icon": "sym_o_local_shipping", "temp": "/product/order"}');`

var BasicAdminMenu = &utils.InitTable{
	Name:        AdminMenuTableName,
	Comment:     AdminMenuTableComment,
	CreateTable: CreateAdminMenu,
	InsertTable: InsertAdminMenu,
}
