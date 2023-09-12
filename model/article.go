package model

import (
    "gorm.io/gorm"

    "github.com/hd2yao/ginblog/utils/err_msg"
)

type Article struct {
    Category Category `gorm:"foreign_key:Cid"`
    gorm.Model
    Title        string `gorm:"type:varchar(100);not null" json:"title"`
    Cid          int    `gorm:"type:int;not null" json:"cid"`
    Desc         string `gorm:"type:varchar(200)" json:"desc"`
    Content      string `gorm:"type:longtext" json:"content"`
    Img          string `gorm:"type:varchar(100)" json:"img"`
    CommentCount int    `gorm:"type:int;not null;default:0" json:"comment_count"`
    ReadCount    int    `gorm:"type:int;not null;default:0" json:"read_count"`
}

// CreateArt 新增文章
func CreateArt(data *Article) int {
    err := db.Create(&data).Error
    if err != nil {
        return err_msg.ERROR // 500
    }
    return err_msg.SUCCESS
}

// GetCateArtList 查询分类下的所有文章
func GetCateArtList(cid int, pageSize int, pageNum int) ([]Article, int, int64) {
    var cateArtList []Article
    var counts int64

    err := db.Preload("Category").
        Where("cid = ?", cid).Find(&cateArtList).
        Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
    db.Model(&cateArtList).Where("cid = ?", cid).Count(&counts)
    if err != nil {
        return nil, err_msg.ErrorCateNotExist, 0
    }
    return cateArtList, err_msg.SUCCESS, counts
}

// GetArtInfo 查询指定文章
func GetArtInfo(id int) (Article, int) {
    var art Article
    err := db.Where("id = ?", id).Preload("Category").First(&art).Error
    db.Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
    if err != nil {
        return art, err_msg.ErrorArticleNotExist
    }
    return art, err_msg.SUCCESS
}

// GetArtList 查询文章列表
func GetArtList(title string, pageSize int, pageNum int) ([]Article, int, int64) {
    var artList []Article
    var counts int64
    err := db.Select("article.id, title, img, created_at, updated_at, `desc`, comment_count, read_count, category.name").
        Order("Created_At DESC").Joins("Category").
        Where("title LIKE ?", "%"+title+"%").
        Limit(pageSize).Offset((pageSize - 1) * pageNum).Find(&artList).Error
    // 单独计数
    db.Model(&artList).Where("title LIKE ?", "%"+title+"%").Count(&counts)
    if err != nil {
        return nil, err_msg.ERROR, 0
    }
    return artList, err_msg.SUCCESS, counts
}

// EditArt 编辑文章
func EditArt(id int, data *Article) int {
    var art Article
    var maps = make(map[string]interface{})
    maps["title"] = data.Title
    maps["cid"] = data.Cid
    maps["desc"] = data.Desc
    maps["content"] = data.Content
    maps["img"] = data.Img

    err := db.Model(&art).Where("id = ?", id).Updates(&maps).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}

// DeleteArt 删除文章
func DeleteArt(id int) int {
    var art Article
    err := db.Where("id = ?", id).Delete(&art).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}
