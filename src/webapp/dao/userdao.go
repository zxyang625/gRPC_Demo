package dao

import (
	"webapp/model"
	"webapp/utils"
)

//根据用户名和密码从数据库中查询一条记录
func CheckUserNameAndPassword(username string, password string) (*model.User, error) {
	sqlStr := "select id, username, password, email from users where username = ? and password = ?"
	row := utils.Db.QueryRow(sqlStr, username, password)
	user := &model.User{}
	row.Scan(user.ID, user.Username, user.Password, user.Email)
	return user, nil
}

//根据用户名从数据库中查询一条记录
func CheckUserName(username string) (*model.User, error) {
	sqlStr := "select id, username, password, email from users where username = ?"
	row := utils.Db.QueryRow(sqlStr, username)
	user := &model.User{}
	row.Scan(user.ID, user.Username, user.Password, user.Email)
	return user, nil
}

//向数据库中插入用户信息
func SaveUser(username string, password string, email string) error {
	sqlStr := "insert into users(username, password, email) values(?,?,?)"
	_, err := utils.Db.Exec(sqlStr, username, password, email)
	if err != nil {
		return err
	}
	return nil
}
