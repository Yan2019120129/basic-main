package sql

import (
	"basic/tools/utils"
)

const FinancialStatisticsTableName = "financial_statistics"
const FinancialStatisticsTableComment = "商城财务统计表"
const CreateFinancialStatistics = `CREATE TABLE ` + FinancialStatisticsTableName + ` (
  id bigint NOT NULL COMMENT '财务统计ID',
  admin_id bigint NOT NULL COMMENT '管理员ID',
  product_id bigint NOT NULL COMMENT '商品ID',
  profit decimal(10,2) DEFAULT NULL COMMENT '利润',
  unit_price decimal(10,2) DEFAULT NULL COMMENT '单价',
  wholesale_price decimal(10,2) DEFAULT NULL COMMENT '批发价',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + FinancialStatisticsTableComment + `';`

var BasicFinancialStatistics = &utils.InitTable{
	Name:        FinancialStatisticsTableName,
	Comment:     FinancialStatisticsTableComment,
	CreateTable: CreateFinancialStatistics,
}
