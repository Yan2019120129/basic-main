package sql

import (
	"basic/tools/utils"
)

const HomeUserLevelTableName = "user_level"
const HomeUserLevelTableComment = "用户等级"
const CreateHomeUserLevel = `CREATE TABLE ` + HomeUserLevelTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
    name VARCHAR(50) NOT NULL DEFAULT '' COMMENT '名称',
    icon VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
    level TINYINT NOT NULL DEFAULT 1 COMMENT '等级',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '购买金额',
	days SMALLINT NOT NULL DEFAULT -1 COMMENT '购买天数 -1无限时间',
    status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
    data VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '数据',
    created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
    updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeUserLevelTableComment + `';`

const InsertHomeUserLevel = `INSERT INTO ` + HomeUserLevelTableName + `(admin_id, name, icon, level, money, days, data) VALUES
(1, 'Level1', '/assets/images/icon/diamond.png', 1, 100, -1, 'level1Content'),
(1, 'Level2', '/assets/images/icon/diamond.png', 2, 500, -1, 'level2Content'),
(1, 'Level3', '/assets/images/icon/diamond.png', 3, 1500, -1, 'level3Content'),
(1, 'Level4', '/assets/images/icon/diamond.png', 4, 3000, -1, 'level4Content'),
(1, 'Level5', '/assets/images/icon/diamond.png', 5, 8000, -1, 'level5Content');`

var BasicHomeUserLevel = &utils.InitTable{
	Name:        HomeUserLevelTableName,
	Comment:     HomeUserLevelTableComment,
	CreateTable: CreateHomeUserLevel,
	InsertTable: InsertHomeUserLevel,
}
