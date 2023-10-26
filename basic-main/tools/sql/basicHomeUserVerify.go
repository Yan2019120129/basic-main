package sql

import (
	"basic/tools/utils"
)

const HomeUserVerifyTableName = "user_verify"
const HomeUserVerifyTableComment = "用户认证"
const CreateHomeUserVerify = `CREATE TABLE ` + HomeUserVerifyTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	user_id  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '类型 1身份证 2护照',
	real_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '真实姓名',
	id_number VARCHAR(50) NOT NULL DEFAULT '' COMMENT '证件号码',
	id_photo1 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '证件照1',
	id_photo2 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '证件照2',
	id_photo3 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '证件照3',
	address VARCHAR(255) NOT NULL DEFAULT '' COMMENT '地址',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -1拒绝｜10审核｜20通过',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeUserVerifyTableComment + `';`

const InsertHomeUserVerify = ``

var BasicHomeUserVerify = &utils.InitTable{
	Name:        HomeUserVerifyTableName,
	Comment:     HomeUserVerifyTableComment,
	CreateTable: CreateHomeUserVerify,
	InsertTable: InsertHomeUserVerify,
}
