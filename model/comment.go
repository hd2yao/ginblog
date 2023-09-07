package model

import (
    "github.com/hd2yao/ginblog/utils/err_msg"
    "gorm.io/gorm"
)

type Comment struct {
    gorm.Model
    UserId    uint   `json:"user_id"`
    ArticleId uint   `json:"article_id"`
    Title     string `json:"article_title"`
    UserName  string `json:"user_name"`
    Content   string `gorm:"type:varchar(500);not null;" json:"content"`
    Status    int8   `gorm:"type:tinyint;default:2" json:"status"`
}

// AddComment 新增评论
func AddComment(data *Comment) int {
    err := db.Create(&data).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}

// GetComment 查询单个评论
func GetComment(id int) (Comment, int) {
    var comment Comment
    err := db.Where("id = ?", id).First(&comment).Error
    if err != nil {
        return comment, err_msg.ERROR
    }
    return comment, err_msg.SUCCESS
}

// GetCommentList 后台所有获取评论列表
func GetCommentList(pageSize, pageNum int) ([]Comment, int, int64) {
    var commentList []Comment
    var counts int64
    db.Find(&commentList).Count(&counts)
    err := db.Model(&commentList).
        Joins("LEFT JOIN article ON comment.article_id = article.id").
        Joins("LEFT JOIN user ON comment.user_id = user.id").
        Select("comment.id, article.title,user_id,article_id, user.username, " +
            "comment.content, comment.status,comment.created_at,comment.deleted_at").
        Order("Created_At DESC").
        Limit(pageSize).Offset((pageNum - 1) * pageSize).
        Scan(&commentList).Error
    if err != nil {
        return commentList, err_msg.ERROR, 0
    }
    return commentList, err_msg.SUCCESS, counts
}

// GetCommentCount 获取评论数量
func GetCommentCount(id int) int64 {
    var comment Comment
    var counts int64
    db.Find(&comment).Where("article_id = ?", id).Where("status = ?", 1).Count(&counts)
    return counts
}

// GetCommentListFront 展示页面获取评论列表
func GetCommentListFront(id int, pageSize, pageNum int) ([]Comment, int, int64) {
    var commentList []Comment
    counts := GetCommentCount(id)
    err := db.Model(&Comment{}).
        Joins("LEFT JOIN article ON comment.article_id = article.id").
        Joins("LEFT JOIN user ON comment.user_id = user.id").
        Select("comment.id, article.title, user_id, article_id, user.username, "+
            "comment.content, comment.status,comment.created_at,comment.deleted_at").
        Where("article_id = ?", id).
        Where("status = ?", 1).
        Order("Created_At DESC").
        Limit(pageSize).Offset((pageNum - 1) * pageSize).
        Scan(&commentList).Error
    if err != nil {
        return commentList, err_msg.ERROR, 0
    }
    return commentList, err_msg.SUCCESS, counts
}

// DeleteComment 删除评论
func DeleteComment(id uint) int {
    var comment Comment
    err := db.Where("id = ?", id).Delete(&comment).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}

// CheckComment 通过评论
func CheckComment(id int, data *Comment) int {
    var comment Comment
    var res Comment
    var article Article
    var maps = make(map[string]interface{})
    maps["status"] = data.Status
    err := db.Model(&comment).Where("id = ?", id).Updates(maps).First(&res).Error
    db.Model(&article).Where("id = ?", res.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}

// UncheckComment 撤下评论
func UncheckComment(id int, data *Comment) int {
    var comment Comment
    var res Comment
    var article Article
    var maps = make(map[string]interface{})
    maps["status"] = data.Status

    err = db.Model(&comment).Where("id = ?", id).Updates(maps).First(&res).Error
    db.Model(&article).Where("id = ?", res.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}
