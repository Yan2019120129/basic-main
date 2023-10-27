package my_test

import (
	"database/sql"
	"github.com/so68/zfeng/database"
	"github.com/so68/zfeng/database/define"
)

// UserCommentAttrs 表示用户评论信息的结构体
type UserCommentAttrs struct {
	ID         int64  `json:"id"`          // 评论ID
	ProductID  int64  `json:"product_id"`  // 商品ID
	StarRating string `json:"star_rating"` // 星级
	UserID     int64  `json:"user_id"`     // 用户ID
	Username   string `json:"username"`    // 用户名
	UserAvatar string `json:"user_avatar"` // 用户头像
	Comment    string `json:"comment"`     // 评论
	Time       int    `json:"time"`        // 时间
	Status     int    `json:"status"`      // 状态：1新增，2追加，3回复
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

// AndWhere 添加 WHERE 子句
func (u *UserComment) AndWhere(str string, args ...interface{}) *UserComment {
	u.Db.AndWhere(str, args...)
	return u
}

// FindOne 查询单个评论信息
func (u *UserComment) FindOne() *UserCommentAttrs {
	comment := new(UserCommentAttrs)
	u.QueryRow(func(row *sql.Row) {
		err := row.Scan(&comment.ID, &comment.ProductID, &comment.StarRating, &comment.UserID, &comment.Username, &comment.UserAvatar, &comment.Comment, &comment.Time, &comment.Status)
		if err != nil {
			comment = nil
		}
	})
	return comment
}

// FindMany 查询多个评论信息
func (u *UserComment) FindMany() []*UserCommentAttrs {
	data := make([]*UserCommentAttrs, 0)
	u.Query(func(rows *sql.Rows) {
		tmp := new(UserCommentAttrs)
		_ = rows.Scan(&tmp.ID, &tmp.ProductID, &tmp.StarRating, &tmp.UserID, &tmp.Username, &tmp.UserAvatar, &tmp.Comment, &tmp.Time, &tmp.Status)
		data = append(data, tmp)
	})
	return data
}
