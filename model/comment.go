package model

import "gorm.io/gorm"

type Comment struct {
    gorm.Model
    UserId    uint   `json:"user_id"`
    ArticleId uint   `json:"article_id"`
    Title     string `json:"article_title"`
    UserName  string `json:"user_name"`
    Content   string `gorm:"type:varchar(500);not null;" json:"content"`
    Status    int8   `gorm:"type:tinyint;default:2" json:"status"`
}
