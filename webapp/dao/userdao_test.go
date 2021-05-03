package dao

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	fmt.Println("测试userdao中的函数")
	t.Run("验证用户名和密码:", testLogin)
	t.Run("验证用户名:", testRegist)
	t.Run("验证保存用户:", testSave)
}

func testLogin(t *testing.T) {
	user, _ := CheckUserNameAndPassword("admin", "123456")
	fmt.Println(user)
}

func testRegist(t *testing.T) {
	user, _ := CheckUserName("admin")
	fmt.Println(user)
}
func testSave(t *testing.T) {
	SaveUser("admin", "123456", "123456@qq.com")
}
