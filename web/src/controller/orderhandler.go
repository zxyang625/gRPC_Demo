package controller

import (
	"dao"
	"html/template"
	"model"
	"net/http"
	"time"
	"utils"
)

//Checkout //去结账
func Checkout(w http.ResponseWriter, r *http.Request) {
	//获取session
	_, session := dao.IsLogin(r)
	//获取用户的id
	userID := session.UserID
	//获取购物车
	cart, _ := dao.GetCartByUserID(userID)
	//生成订单号
	orderID := utils.CreateUUID()
	//创建生成订单的时间
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	//创建Order
	order := &model.Order{
		OrderID:     orderID,
		CreateTime:  timeStr,
		TotalCount:  cart.TotalCount,
		TotalAmount: cart.TotalAmount,
		State:       0,
		UserID:      int64(userID),
	}
	//将订单保存到数据库中
	dao.AddOrder(order)
	//保存订单项
	//获取购物车中的购物项
	cartItems := cart.CartItems
	//遍历得到每一个购物项
	for _, v := range cartItems {
		//创建OrderItem
		orderItem := &model.OrderItem{
			Count:   v.Count,
			Amount:  v.Amount,
			Title:   v.Book.Title,
			Author:  v.Book.Author,
			Price:   v.Book.Price,
			ImgPath: v.Book.ImgPath,
			OrderID: orderID,
		}
		//将购物项保存到数据库中
		dao.AddOrderItem(orderItem)
		//更新当前购物项中图书的库存和销量
		book := v.Book
		book.Sales = book.Sales + int(v.Count)
		book.Stock = book.Stock - int(v.Count)
		//更新图书的信息
		dao.UpdateBook(book)
	}
	//清空购物车
	dao.DeleteCartByCartID(cart.CartID)
	//将订单号设置到session中
	session.OrderID = orderID
	//解析模板
	t := template.Must(template.ParseFiles("views/pages/cart/checkout.html"))
	//执行
	t.Execute(w, session)
}

//GetOrders 获取所有订单
func GetOrders(w http.ResponseWriter, r *http.Request) {
	//调用dao中获取所有订单的函数
	orders, _ := dao.GetOrders()
	//解析模板
	t := template.Must(template.ParseFiles("views/pages/order/order_manager.html"))
	//执行
	t.Execute(w, orders)
}

//GetOrderInfo 获取订单对应的订单项
func GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	//获取订单号
	orderID := r.FormValue("orderId")
	//根据订单号调用dao中获取所有订单项的函数
	orderItems, _ := dao.GetOrderItemsByOrderID(orderID)
	//解析模板
	t := template.Must(template.ParseFiles("views/pages/order/order_info.html"))
	//执行
	t.Execute(w, orderItems)
}



//SendOrder 发货
func SendOrder(w http.ResponseWriter, r *http.Request) {
	//获取要发货的订单号
	orderID := r.FormValue("orderId")
	//调用dao中的更新订单状态的函数
	dao.UpdateOrderState(orderID, 1)
	//调用GetOrders函数再次查询一下所有的订单
	GetOrders(w, r)
}

//TakeOrder 收货
func TakeOrder(w http.ResponseWriter, r *http.Request) {
	//获取要收货的订单号
	orderID := r.FormValue("orderId")
	//调用dao中的更新订单状态的函数
	dao.UpdateOrderState(orderID, 2)
	//调用获取我的订单的函数再次查询我的订单
	GetMyOrders(w, r)
}

//GetMyOrders 获取我的订单
func GetMyOrders(w http.ResponseWriter, r *http.Request) {
	//获取session
	_, session := dao.IsLogin(r)
	//获取用户的id
	userID := session.UserID
	//调用dao中获取用户的所有订单的函数
	orders, _ := dao.GetMyOrders(userID)
	//将订单设置到session中
	session.Orders = orders
	//解析模板
	t := template.Must(template.ParseFiles("views/pages/order/order.html"))
	//执行
	t.Execute(w, session)
}