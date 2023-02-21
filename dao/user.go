package dao

import "log"

// User 表结构体
type User struct {
	Id       int64  `gorm:"column:user_id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

// 表名映射
func (User) TableName() string {
	return "user"
}

// 获取全部用户
func QueryAllUser() ([]User, error) {
	users := []User{}
	err := db.Find(&users).Error
	if err != nil {
		log.Println(err.Error())
		return users, err
	}
	return users, nil
}

// 根据id查找指定用户
func QueryUserById(id int64) (User, error) {
	user := User{}
	err := db.Where("user_id=?", id).First(&user).Error
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// 根据username查找指定用户
func QueryUserByUsername(name string) (User, error) {
	user := User{}
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// 新增用户
func InsertUser(user *User) bool {
	err := db.Create(&user).Error
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
