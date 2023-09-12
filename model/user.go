package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/hd2yao/ginblog/utils/err_msg"
)

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(20);not null" json:"user_name" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

// CheckUser 查询用户是否存在
func CheckUser(name string) int {
	var user User
	db.Select("id").Where("username = ?", name).First(&user)
	if user.ID > 0 {
		return err_msg.ErrorUsernameUsed
	}
	return err_msg.SUCCESS
}

// CheckUpUser 更新查询
func CheckUpUser(id int, name string) int {
	var user User
	db.Select("id, username").Where("username = ?", name).First(&user)
	if user.ID == uint(id) {
		return err_msg.SUCCESS
	}
	if user.ID > 0 {
		return err_msg.ErrorUsernameUsed
	}
	return err_msg.SUCCESS
}

// CreateUser 新增用户
func CreateUser(data *User) int {
	err := db.Create(&data).Error
	if err != nil {
		return err_msg.ERROR
	}
	return err_msg.SUCCESS
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Where("id = ?", id).Limit(1).Find(&user).Error
	if err != nil {
		return user, err_msg.ERROR
	}
	return user, err_msg.SUCCESS
}

// GetUsers 查询用户列表
func GetUsers(userName string, pageSize, pageNum int) ([]User, int64) {
	var users []User
	var counts int64

	if userName != "" {
		db.Select("id,user_name,role,created_at").
			Where("user_name LIKE ?", "%"+userName+"%").
			Offset((pageNum - 1) * pageSize).Limit(pageSize).
			Find(&users)
		db.Model(&users).
			Where("user_name LIKE ?", "%"+userName+"%").
			Count(&counts)
		return users, counts
	}
	db.Select("id,user_name,role,created_at").
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&users)
	db.Model(&users).Count(&counts)

	return users, counts
}

// EditUser 编辑用户信息
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["user_name"] = data.UserName
	maps["role"] = data.Role
	err := db.Model(&user).Where("id = ?", id).Updates(&maps).Error
	if err != nil {
		return err_msg.ERROR
	}
	return err_msg.SUCCESS
}

// ChangePassword 修改密码
func ChangePassword(id int, data *User) int {
	//var user User
	//var maps = make(map[string]interface{})
	//maps["password"] = data.Password

	err := db.Select("password").Where("id = ?", id).Updates(&data).Error
	if err != nil {
		return err_msg.ERROR
	}
	return err_msg.SUCCESS
}

// DeleteUser 删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err_msg.ERROR
	}
	return err_msg.SUCCESS
}

// BeforeCreate 密码加密&权限控制
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	u.Role = 2
	return nil
}

// BeforeUpdate 密码加密
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

// ScryptPw 生成密码
func ScryptPw(password string) string {
	const cost = 10

	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal(err)
	}
	return string(HashPw)
}

// CheckLogin 后台登录验证
func CheckLogin(userName string, password string) (User, int) {
	var user User
	var passwordErr error

	db.Where("user_name = ?", userName).First(&user)

	passwordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user.ID == 0 {
		return user, err_msg.ErrorUserNotExist
	}
	if passwordErr != nil {
		return user, err_msg.ErrorPasswordWrong
	}
	if user.Role != 1 {
		return user, err_msg.ErrorUserNoRight
	}
	return user, err_msg.SUCCESS
}

// CheckLoginFront 前台登录
func CheckLoginFront(userName string, password string) (User, int) {
	var user User
	var passwordErr error

	db.Where("user_name = ?", userName).First(&user)

	passwordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user.ID == 0 {
		return user, err_msg.ErrorUserNotExist
	}
	if passwordErr != nil {
		return user, err_msg.ErrorPasswordWrong
	}
	return user, err_msg.SUCCESS
}
