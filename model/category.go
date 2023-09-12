package model

import (
    "gorm.io/gorm"

    "github.com/hd2yao/ginblog/utils/err_msg"
)

type Category struct {
    ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
    Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// CategoryIsExist 检查分类名是否存在
func CategoryIsExist(name string) int {
    var cate Category
    db.Select("id").Where("name LIKE ?", "%"+name+"%").First(&cate)
    if cate.ID > 0 {
        return err_msg.ErrorCateNameUsed
    }
    return err_msg.SUCCESS
}

// CreateCate 新增分类
func CreateCate(data *Category) int {
    err := db.Create(&data).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}

// GetCateInfo 查询指定的分类信息
func GetCateInfo(id int) (Category, int) {
    var cate Category
    db.Where("id = ?", id).First(&cate)
    return cate, err_msg.SUCCESS
}

// GetCateList 查询分类列表
func GetCateList(name string, pageSize, pageNum int) ([]Category, int, int64) {
    var cateList []Category
    var counts int64
    err := db.Where("name LIKE ?", "%"+name+"%").Limit(pageSize).Offset((pageSize - 1) * pageNum).
        Find(&cateList).Error
    db.Model(&cateList).Where("name LIKE ?", "%"+name+"%").Count(&counts)
    if err != nil && err != gorm.ErrRecordNotFound {
        return nil, err_msg.ERROR, 0
    }
    return cateList, err_msg.SUCCESS, counts
}

// EditCate 编辑分类信息
func EditCate(id int, data *Category) int {
    var cate Category
    var maps = make(map[string]interface{})
    maps["name"] = data.Name
    err := db.Model(&cate).Updates(&maps).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}

// DeleteCate 删除分类
func DeleteCate(id int) int {
    var cate Category
    err := db.Where("id = ?", id).Delete(&cate).Error
    if err != nil {
        return err_msg.ERROR
    }
    return err_msg.SUCCESS
}
