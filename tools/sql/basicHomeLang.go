package sql

import (
	"basic/tools/utils"
)

const BasicHomeLangTableName = "lang"
const BasicHomeLangTableComment = "用户语言"
const CreateBasicHomeLang = `CREATE TABLE ` + BasicHomeLangTableName + ` (
	id         	 INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    admin_id   	 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
    name       	 VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '名称',
	alias		 VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '别名',
    icon       	 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
    sort       	 TINYINT      NOT NULL DEFAULT 99 COMMENT '排序',
    status     	 TINYINT      NOT NULL DEFAULT 10 COMMENT '状态 -1禁用｜10启用',
    data       	 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
    created_at 	 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicHomeLangTableComment + `';`

const InsertBasicHomeLang = `INSERT INTO ` + BasicHomeLangTableName + `(admin_id, name, alias, icon, status, data) VALUES
(1, '简体中文', 'zh-CN', '/assets/images/country/china.png', 10, '简体中文'),
(1, '繁體中文','zh-TW', '/assets/images/country/taiwan.png', -1, '繁体中文'),
(1, 'English','en-US', '/assets/images/country/usa.png', -1, '英格兰语'),
(1, 'عربي', 'ar-AE', '/assets/images/country/united_arab_emirates.png', -1, '阿拉伯语'),
(1, 'беларускі', 'be-BY', '/assets/images/country/belarus.png', -1, '白俄罗斯语'),
(1, 'български', 'bg-BG', '/assets/images/country/bulgaria.png', -1, '保加利亚语'),
(1, 'čeština', 'cs-CZ', '/assets/images/country/czech.png', -1, '捷克语'),
(1, 'dansk', 'da-DK', '/assets/images/country/denmark.png', -1, '丹麦语'),
(1, 'Deutsch', 'de-DE', '/assets/images/country/germany.png', -1, '德语'),
(1, 'Ελληνικά', 'el-GR', '/assets/images/country/greece.png', -1, '希腊语'),
(1, 'español', 'es-ES', '/assets/images/country/spain.png', -1, '西班牙语'),
(1, 'eesti keel', 'et-EE', '/assets/images/country/estonia.png', -1, '爱沙尼亚语'),
(1, 'Suomalainen', 'fi-FI', '/assets/images/country/finland.png', -1, '芬兰语'),
(1, 'Français', 'fr-FR', '/assets/images/country/france.png', -1, '法语'),
(1, 'Hrvatski', 'hr-HR', '/assets/images/country/croatia.png', -1, '克罗地亚语'),
(1, 'Magyar', 'hu-HU', '/assets/images/country/hungary.png', -1, '匈牙利语'),
(1, 'íslenskur', 'is-IS', '/assets/images/country/iceland.png', -1, '冰岛语'),
(1, 'italiano', 'it-IT', '/assets/images/country/italy.png', -1, '意大利语'),
(1, '日本', 'ja-JP', '/assets/images/country/japan.png', -1, '日语'),
(1, 'Melayu', 'ms-MY', '/assets/images/country/malaysia.png', -1, '马来语'),
(1, 'Tiếng Việt', 'vi-VN', '/assets/images/country/vietnam.png', -1, '越南语'),
(1, '한국인', 'ko-KR', '/assets/images/country/north_korea.png', -1, '朝鲜语(韩语)'),
(1, 'lietuvių', 'lt-LT', '/assets/images/country/lithuania.png', -1, '立陶宛语'),
(1, 'македонски', 'mk-MK', '/assets/images/country/macedonia.png', -1, '马其顿语'),
(1, 'Nederlands', 'nl-NL', '/assets/images/country/netherlands.png', -1, '荷兰语'),
(1, 'norsk', 'no-NO', '/assets/images/country/norway.png', -1, '挪威语'),
(1, 'Polski', 'pl-PL', '/assets/images/country/poland.png', -1, '波兰语'),
(1, 'Português', 'pt-PT', '/assets/images/country/portugal.png', -1, '葡萄牙语'),
(1, 'Română', 'ro-RO', '/assets/images/country/romania.png', -1, '罗马尼亚语'),
(1, 'Русский', 'ru-RU', '/assets/images/country/russia.png', -1, '俄语'),
(1, 'Hrvatski', 'sh-YU', '/assets/images/country/croatia.png', -1, '克罗地亚语'),
(1, 'slovenský', 'sk-SK', '/assets/images/country/slovakia.png', -1, '斯洛伐克语'),
(1, 'Slovenščina', 'sl-SI', '/assets/images/country/slovenia.png', -1, '斯洛文尼亚语'),
(1, 'shqiptare', 'sq-AL', '/assets/images/country/albania.png', -1, '阿尔巴尼亚语'),
(1, 'svenska', 'sv-SE', '/assets/images/country/sweden.png', -1, '瑞典语'),
(1, 'แบบไทย', 'th-TH', '/assets/images/country/thailand.png', -1, '泰语'),
(1, 'Türkçe', 'tr-TR', '/assets/images/country/turkey.png', -1, '土耳其语'),
(1, 'українська', 'uk-UA', '/assets/images/country/ukraine.png', -1, '乌克兰语'),
(1, 'Српски', 'sr-YU', '/assets/images/country/serbia.png', -1, '塞尔维亚语'),
(1, 'עִברִית', 'iw-IL', '/assets/images/country/israel.png', -1, '希伯来语'),
(1, 'हिंदी', 'hi-IN', '/assets/images/country/india.png', -1, '印地语'),
(1, 'Indonesia', 'id-ID', '/assets/images/country/indonesia.png', -1, '印尼语');`

var BasicHomeLang = &utils.InitTable{
	Name:        BasicHomeLangTableName,
	Comment:     BasicHomeLangTableComment,
	CreateTable: CreateBasicHomeLang,
	InsertTable: InsertBasicHomeLang,
}
