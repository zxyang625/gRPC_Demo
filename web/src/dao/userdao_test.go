package dao

import (
	"fmt"
	"model"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// fmt.Println("测试bookdao中的方法")
	m.Run()
}

func TestUser(t *testing.T) {
	// fmt.Println("测试userdao中的函数")
	// t.Run("验证用户名或密码：", testLogin)
	// t.Run("验证用户名：", testRegist)
	// t.Run("保存用户：", testSave)
}

func testLogin(t *testing.T) {
	user, _ := CheckUserNameAndPassword("admin", "123456")
	fmt.Println("获取用户信息是：", user)
}
func testRegist(t *testing.T) {
	user, _ := CheckUserName("admin")
	fmt.Println("获取用户信息是：", user)
}
func testSave(t *testing.T) {
	SaveUser("admin3", "123456", "admin@atguigu.com")
}

func TestBook(t *testing.T) {
	// fmt.Println("测试bookdao中的相关函数")
	// t.Run("测试获取所有图书", testGetBooks)
	// t.Run("测试添加图书", testAddBook)
	// t.Run("测试删除图书", testDeleteBook)
	// t.Run("测试获取一本图书", testGetBook)
	// t.Run("测试更新图书", testUpdateBook)
	// t.Run("测试获取带分页的图书", testGetPageBooks)
	// t.Run("测试获取带分页和价格范围的图书", testGetPageBooksByPrice)
}

func testGetBooks(t *testing.T) {
	books, _ := GetBooks()
	//遍历得到每一本图书
	for k, v := range books {
		fmt.Printf("第%v本图书的信息是：%v\n", k+1, v)
	}
}
func testAddBook(t *testing.T) {
	book := &model.Book{
		Title:   "三国演义",
		Author:  "罗贯中",
		Price:   88.88,
		Sales:   100,
		Stock:   100,
		ImgPath: "/static/img/default.jpg",
	}
	//调用添加图书的函数
	AddBook(book)
}
func testDeleteBook(t *testing.T) {
	//调用删除图书的函数
	DeleteBook("34")
}
func testGetBook(t *testing.T) {
	//调用获取图书的函数
	book, _ := GetBookByID("32")
	fmt.Println("获取的图书信息是：", book)
}
func testUpdateBook(t *testing.T) {
	book := &model.Book{
		ID:      32,
		Title:   "3个女人与105个男人的故事",
		Author:  "罗贯中",
		Price:   66.66,
		Sales:   10000,
		Stock:   1,
		ImgPath: "/static/img/default.jpg",
	}
	//调用更新图书的函数
	UpdateBook(book)
}

func testGetPageBooks(t *testing.T) {
	page, _ := GetPageBooks("9")
	fmt.Println("当前页是：", page.PageNo)
	fmt.Println("总页数是：", page.TotalPageNo)
	fmt.Println("总记录数是：", page.TotalRecord)
	fmt.Println("当前页中的图书有：")
	for _, v := range page.Books {
		fmt.Println("图书的信息是：", v)
	}
}
func testGetPageBooksByPrice(t *testing.T) {
	page, _ := GetPageBooksByPrice("3", "10", "30")
	fmt.Println("当前页是：", page.PageNo)
	fmt.Println("总页数是：", page.TotalPageNo)
	fmt.Println("总记录数是：", page.TotalRecord)
	fmt.Println("当前页中的图书有：")
	for _, v := range page.Books {
		fmt.Println("图书的信息是：", v)
	}
}

func TestSession(t *testing.T) {
	// fmt.Println("测试Session相关函数")
	// t.Run("测试添加Session", testAddSession)
	// t.Run("测试删除Session", testDeleteSession)
	// t.Run("测试获取Session", testGetSession)
}

func testAddSession(t *testing.T) {
	sess := &model.Session{
		SessionID: "13838381438",
		UserName:  "马蓉",
		UserID:    5,
	}
	AddSession(sess)
}

func testDeleteSession(t *testing.T) {
	DeleteSession("13838381438")
}
func testGetSession(t *testing.T) {
	sess, _ := GetSession("c65d2a76-9447-44cc-5fe8-c183e1414076")
	fmt.Println("Session的信息是：", sess)
}

func TestCart(t *testing.T) {
	// fmt.Println("测试购物车的相关函数")
	// t.Run("测试添加购物车", testAddCart)
	// t.Run("测试根据图书的id获取对应的购物项", testGetCartItemByBookID)
	// t.Run("测试根据购物车的id获取所有的购物项", testGetCartItemsByCartID)
	// t.Run("测试根据用户的id获取对应的购物车", testGetCartByUserID)
	// t.Run("测试根据图书的id和购物车的id以及输入的图书的数量更新购物项", testUpdateBookCount)
	// t.Run("测试购物车的id删除购物项和购物车", testDeleteCartByCartID)
	// t.Run("测试删除购物项", testDeleteCartItemByID)
}

func testAddCart(t *testing.T) {
	//设置要买的第一本书
	book := &model.Book{
		ID:    1,
		Price: 27.20,
	}
	//设置要买的第二本书
	book2 := &model.Book{
		ID:    2,
		Price: 23.00,
	}
	//创建一个购物项切片
	var cartItems []*model.CartItem
	//创建两个购物项
	cartItem := &model.CartItem{
		Book:   book,
		Count:  10,
		CartID: "66668888",
	}
	cartItems = append(cartItems, cartItem)
	cartItem2 := &model.CartItem{
		Book:   book2,
		Count:  10,
		CartID: "66668888",
	}
	cartItems = append(cartItems, cartItem2)
	//创建购物车
	cart := &model.Cart{
		CartID:    "66668888",
		CartItems: cartItems,
		UserID:    1,
	}
	//将购物车插入到数据库中
	AddCart(cart)
}

func testGetCartItemByBookID(t *testing.T) {
	cartItem, _ := GetCartItemByBookIDAndCartID("1", "66668888")
	fmt.Println("图书id=1的购物项的信息是：", cartItem)
}
func testGetCartItemsByCartID(t *testing.T) {
	cartItems, _ := GetCartItemsByCartID("66668888")
	for k, v := range cartItems {
		fmt.Printf("第%v个购物项是：%v\n", k+1, v)
	}
}

func testGetCartByUserID(t *testing.T) {
	cart, _ := GetCartByUserID(3)
	fmt.Println("id为2的用户的购物车信息是：", cart)
}

func testUpdateBookCount(t *testing.T) {
	// UpdateBookCount(100, 1, "66668888")
}
func testDeleteCartByCartID(t *testing.T) {
	DeleteCartByCartID("80bb8008-8383-47d0-4694-5eae94f39ffd")
}

func testDeleteCartItemByID(t *testing.T) {
	DeleteCartItemByID("21")
}

func TestOrder(t *testing.T) {
	fmt.Println("测试订单相关函数")
	// t.Run("测试添加订单和订单项", testAddOrder)
	// t.Run("测试获取所有的订单", testGetOrders)
	// t.Run("测试获取所有的订单项", testGetOrderItems)
	// t.Run("测试获取我的订单", testGetMyOrders)
	t.Run("测试发货和收货", testUpdateOrderState)

}

func testAddOrder(t *testing.T) {
	//生成订单号
	orderID := "88888888"
	//创建订单
	order := &model.Order{
		OrderID:     orderID,
		CreateTime:  time.Now().String(),
		TotalCount:  2,
		TotalAmount: 400,
		State:       0,
		UserID:      1,
	}
	//创建订单项
	orderItem := &model.OrderItem{
		Count:   1,
		Amount:  300,
		Title:   "三国演义",
		Author:  "罗贯中",
		Price:   300,
		ImgPath: "/static/img/default.jpg",
		OrderID: orderID,
	}
	orderItem2 := &model.OrderItem{
		Count:   1,
		Amount:  100,
		Title:   "西游记",
		Author:  "吴承恩",
		Price:   100,
		ImgPath: "/static/img/default.jpg",
		OrderID: orderID,
	}
	//保存订单
	AddOrder(order)
	//保存订单项
	AddOrderItem(orderItem)
	AddOrderItem(orderItem2)
}
func testGetOrders(t *testing.T) {
	orders, _ := GetOrders()
	for _, v := range orders {
		fmt.Println("订单信息是：", v)
	}
}
func testGetOrderItems(t *testing.T) {
	orderItems, _ := GetOrderItemsByOrderID("9a738546-d240-4c1a-7b3d-a3837100977f")
	for _, v := range orderItems {
		fmt.Println("订单项的信息是：", v)
	}
}

func testGetMyOrders(t *testing.T) {
	orders, _ := GetMyOrders(2)
	for _, v := range orders {
		fmt.Println("我的订单有：", v)
	}
}
func testUpdateOrderState(t *testing.T) {
	UpdateOrderState("5823f37e-f4e6-4a39-7567-2a7c0fd8e638", 1)
}
