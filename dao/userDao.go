package dao

import "log"

// TableUser: 对应数据库 User 表结构
type TableUser struct {
	Id       int64				// 编号
	Name     string				// 姓名
	Password string				// 密码
}

// TableName: 修改表名映射
func (tableUser TableUser) TableName() string {
	return "users"
}

// GetTableUserList: 获取全部 TableUser 对象
func GetTableUserList() ([]TableUser, error) {
	// 创建一个空的 TableUser 对象切片，用于存储查询结果
	tableUsers := []TableUser{}
	// 通过数据库 Db 查询，结果保存于 tableUser 中
	if err := Db.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	// 查询成功返回 TableUser 对象的切片
	return tableUsers, nil
}

// GetTableUserByUsername: 根据 username 获得 TableUser 对象
func GetTableUserByUsername(name string) (TableUser, error) {
	// 创建一个空的 TableUser 对象，用于存储查询结果
	tableUser := TableUser{}
	// 通过数据库 Db 查询，结果保存于 tableUser 中
	if err := Db.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	// 查询成功返回 TableUser 对象
	return tableUser, nil
}

// GetTableUserById: 根据 user_id 获得 TableUser 对象
func GetTableUserById(id int64) (TableUser, error) {
	// 创建一个空的 TableUser 对象，用于存储查询结果
	tableUser := TableUser{}
	if err := Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	// 查询成功返回 TableUser 对象
	return tableUser, nil
}

// InsertTableUser: 将 tableUser 插入数据表中
func InsertTableUser(tableUser *TableUser) bool {
	// 传入一个 TableUser 对象，将其插入数据表中
	if err := Db.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	// 插入成功
	return true
}

// Db.Find() 			代表查询
// Db.Where().First()	筛选
// Db.Create()			插入