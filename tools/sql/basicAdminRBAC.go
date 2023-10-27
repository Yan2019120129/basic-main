package sql

import (
	"basic/tools/utils"
)

const AdminRBACItemTableName = "admin_auth_item"
const AdminRBACItemTableComment = "管理权限item"
const CreateAdminRBACItem = `CREATE TABLE ` + AdminRBACItemTableName + ` (
	name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '名称',
	type TINYINT NOT NULL DEFAULT 0 COMMENT '类型',
	description VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
	rule_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '规则名称',
	data VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	PRIMARY KEY (name),
	KEY auth_item_type (type)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AdminRBACItemTableComment + `';`
const InsertAdminRBACItem = ``

var BasicAdminRBACItem = &utils.InitTable{
	Name:        AdminRBACItemTableName,
	Comment:     AdminRBACItemTableComment,
	CreateTable: CreateAdminRBACItem,
	InsertTable: InsertAdminRBACItem,
}

const AdminRBACChildTableName = "admin_auth_child"
const AdminRBACChildTableComment = "管理权限child"
const CreateAdminRBACChild = `CREATE TABLE ` + AdminRBACChildTableName + ` (
	parent VARCHAR(64) NOT NULL COMMENT '父级',
	child VARCHAR(64) NOT NULL COMMENT '子级',
	type TINYINT NOT NULL DEFAULT 0 COMMENT '类型',
	PRIMARY KEY (parent, child),
	KEY auth_child_child (child),
	CONSTRAINT auth_item_child_1 FOREIGN KEY (parent) REFERENCES admin_auth_item (name) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT auth_item_child_2 FOREIGN KEY (child) REFERENCES admin_auth_item (name) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AdminRBACChildTableComment + `';`
const InsertAdminRBACChild = ``

var BasicAdminRBACChild = &utils.InitTable{
	Name:        AdminRBACChildTableName,
	Comment:     AdminRBACChildTableComment,
	CreateTable: CreateAdminRBACChild,
	InsertTable: InsertAdminRBACChild,
}

const AdminRBACAssignmentTableName = "admin_auth_assignment"
const AdminRBACAssignmentTableComment = "管理权限assignment"
const CreateAdminRBACAssignment = `CREATE TABLE ` + AdminRBACAssignmentTableName + ` (
	item_name VARCHAR(64) NOT NULL COMMENT '名称',
	user_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
	PRIMARY KEY (item_name, user_id),
	CONSTRAINT auth_assignment_1 FOREIGN KEY (item_name) REFERENCES admin_auth_item (name) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AdminRBACAssignmentTableComment + `';`
const InsertAdminRBACAssignment = ``

var BasicAdminRBACAssignment = &utils.InitTable{
	Name:        AdminRBACAssignmentTableName,
	Comment:     AdminRBACAssignmentTableComment,
	CreateTable: CreateAdminRBACAssignment,
	InsertTable: InsertAdminRBACAssignment,
}
