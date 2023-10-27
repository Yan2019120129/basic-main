package models

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserCommentAttrs 数据库模型属性
type UserCommentAttrs struct {
	Id         int64  `json:"id"`          //评论ID
	UserId     int64  `json:"user_id"`     //用户ID
	AdminId    int64  `json:"admin_id"`    //
	ProductId  int64  `json:"product_id"`  //商品ID
	StarRating string `json:"star_rating"` //星级
	Username   string `json:"username"`    //用户名
	UserAvatar string `json:"user_avatar"` //用户头像
	Comment    string `json:"comment"`     //评论
	Time       int64  `json:"time"`        //时间
	Status     int64  `json:"status"`      //状态：1新增，2追加，3回复
}

// UserComment 数据库模型
type UserComment struct {
	define.Db
}

// NewUserComment 创建数据库模型
func NewUserComment(tx *sql.Tx) *UserComment {
	return &UserComment{
		database.DbPool.NewDb(tx).Table("user_comment"),
	}
}

// AndWhere where条件
func (c *UserComment) AndWhere(str string, arg ...any) *UserComment {
	c.Db.AndWhere(str, arg...)
	return c
}

// FindOne 查询单挑
func (c *UserComment) FindOne() *UserCommentAttrs {
	attrs := new(UserCommentAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.UserId, &attrs.ProductId, &attrs.StarRating, &attrs.Username, &attrs.UserAvatar, &attrs.Comment, &attrs.Time, &attrs.Status, &attrs.AdminId)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserComment) FindMany() []*UserCommentAttrs {
	data := make([]*UserCommentAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserCommentAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.ProductId, &tmp.StarRating, &tmp.Username, &tmp.UserAvatar, &tmp.Comment, &tmp.Time, &tmp.Status, &tmp.AdminId)
		data = append(data, tmp)
	})
	return data
}
