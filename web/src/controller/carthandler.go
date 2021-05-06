package controller

import (
	"dao"
	"encoding/json"
	"fmt"
	"html/template"
	"model"
	"net/http"
	"strconv"
	"utils"
)

//AddBook2Cart 添加图书到购物车
func AddBook2Cart(w http.ResponseWriter, r *http.Request) {
	//判断是否登录
	flag, session := dao.IsLogin(r)
	if flag {
		//已经登录
		//获取要添加的图书的id
		bookID := r.FormValue("bookId")
		//根据图书的id获取图书信息
		book, _ := dao.GetBookByID(bookID)

		//获取用户的id
		userID := session.UserID
		//判断数据库中是否有当前用户的购物车
		cart, _ := dao.GetCartByUserID(userID)
		if cart != nil {
			//当前用户已经有购物车，此时需要判断购物车中是否有当前这本图书
			carItem, _ := dao.GetCartItemByBookIDAndCartID(bookID, cart.CartID)
			if carItem != nil {
				//购物车的购物项中已经有该图书，只需要将该图书所对应的购物项中的数量加1即可
				//1.获取购物车切片中的所有的购物项
				cts := cart.CartItems
				//2.遍历得到每一个购物项
				for _, v := range cts {
					fmt.Println("当前购物项中是否有Book：", v)
					fmt.Println("查询到的Book是：", carItem.Book)
					//3.找到当前的购物项
					if v.Book.ID == carItem.Book.ID {
						//将购物项中的图书的数量加1
						v.Count = v.Count + 1
						//更新数据库中该购物项的图书的数量
						dao.UpdateBookCount(v)
					}
				}
			} else {
				//购物车的购物项中还没有该图书，此时需要创建一个购物项并添加到数据库中
				//创建购物车中的购物项
				cartItem := &model.CartItem{
					Book:   book,
					Count:  1,
					CartID: cart.CartID,
				}
				//将购物项添加到当前cart的切片中
				cart.CartItems = append(cart.CartItems, cartItem)
				//将新创建的购物项添加到数据库中
				dao.AddCartItem(cartItem)
			}
			//不管之前购物车中是否有当前图书对应的购物项，都需要更新购物车中的图书的总数量和总金额
			dao.UpdateCart(cart)
		} else {
			//证明当前用户还没有购物车，需要创建一个购物车并添加到数据库中
			//1.创建购物车
			//生成购物车的id
			cartID := utils.CreateUUID()
			cart := &model.Cart{
				CartID: cartID,
				UserID: userID,
			}
			//2.创建购物车中的购物项
			//声明一个CartItem类型的切片
			var cartItems []*model.CartItem
			cartItem := &model.CartItem{
				Book:   book,
				Count:  1,
				CartID: cartID,
			}
			//将购物项添加到切片中
			cartItems = append(cartItems, cartItem)
			//3将切片设置到cart中
			cart.CartItems = cartItems
			//4.将购物车cart保存到数据库中
			dao.AddCart(cart)
		}
		w.Write([]byte("您刚刚将" + book.Title + "添加到了购物车！"))
	} else {
		//没有登录
		w.Write([]byte("请先登录！"))
	}
}

//GetCartInfo 根据用户的id获取购物车信息
func GetCartInfo(w http.ResponseWriter, r *http.Request) {
	_, session := dao.IsLogin(r)
	//获取用户的id
	userID := session.UserID
	//根据用户的id从数据库中获取对应的购物车
	cart, _ := dao.GetCartByUserID(userID)
	if cart != nil {
		//将购物车设置到session中
		session.Cart = cart
		//解析模板文件
		t := template.Must(template.ParseFiles("views/pages/cart/cart.html"))
		//执行
		t.Execute(w, session)
	} else {
		//该用户还没有购物车
		//解析模板文件
		t := template.Must(template.ParseFiles("views/pages/cart/cart.html"))
		//执行
		t.Execute(w, session)
	}
}

//DeleteCart 清空购物车
func DeleteCart(w http.ResponseWriter, r *http.Request) {
	//获取要删除的购物车的id
	cartID := r.FormValue("cartId")
	//清空购物车
	dao.DeleteCartByCartID(cartID)
	//调用GetCartInfo函数再次查询购物车信息
	GetCartInfo(w, r)
}


//UpdateCartItem 更新购物项
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	//获取要更新的购物项的id
	cartItemID := r.FormValue("cartItemId")
	//将购物项的id转换为int64
	iCartItemID, _ := strconv.ParseInt(cartItemID, 10, 64)
	//获取用户输入的图书的数量
	bookCount := r.FormValue("bookCount")
	iBookCount, _ := strconv.ParseInt(bookCount, 10, 64)
	//获取session
	_, session := dao.IsLogin(r)
	//获取用户的id
	userID := session.UserID
	//获取该用户的购物车
	cart, _ := dao.GetCartByUserID(userID)
	//获取购物车中的所有的购物项
	cartItems := cart.CartItems
	//遍历得到每一个购物项
	for _, v := range cartItems {
		//寻找要更新的购物项
		if v.CartItemID == iCartItemID {
			//这个就是我们要更新的购物项
			//将当前购物项中的图书的数量设置为用户输入的值
			v.Count = iBookCount
			//更新数据库中该购物项的图书的数量和金额小计
			dao.UpdateBookCount(v)
		}
	}
	//更新购物车中的图书的总数量和总金额
	dao.UpdateCart(cart)
	//调用获取购物项信息的函数再次查询购物车信息
	cart, _ = dao.GetCartByUserID(userID)
	// GetCartInfo(w, r)
	//获取购物车中图书的总数量
	totalCount := cart.TotalCount
	//获取购物车中图书的总金额
	totalAmount := cart.TotalAmount
	var amount float64
	//获取购物车中更新的购物项中的金额小计
	cIs := cart.CartItems
	for _, v := range cIs {
		if iCartItemID == v.CartItemID {
			//这个就是我们寻找的购物项，此时获取当前购物项中的金额小计
			amount = v.Amount
		}
	}
	//创建Data结构
	data := model.Data{
		Amount:      amount,
		TotalAmount: totalAmount,
		TotalCount:  totalCount,
	}
	//将data转换为json字符串
	json, _ := json.Marshal(data)
	//响应到浏览器
	w.Write(json)
}

//DeleteCartItem 删除购物项
func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	//获取要删除的购物项的id
	cartItemID := r.FormValue("cartItemId")
	//将购物项的id转换为int64
	iCartItemID, _ := strconv.ParseInt(cartItemID, 10, 64)
	//获取session
	_, session := dao.IsLogin(r)
	//获取用户的id
	userID := session.UserID
	//获取该用户的购物车
	cart, _ := dao.GetCartByUserID(userID)
	//获取购物车中的所有的购物项
	cartItems := cart.CartItems
	//遍历得到每一个购物项
	for k, v := range cartItems {
		//寻找要删除的购物项
		if v.CartItemID == iCartItemID {
			//这个就是我们要删除的购物项
			//将当前购物项从切片中移出
			cartItems = append(cartItems[:k], cartItems[k+1:]...)
			//将删除购物项之后的切片再次赋给购物车中的切片
			cart.CartItems = cartItems
			//将当前购物项从数据库中删除
			dao.DeleteCartItemByID(cartItemID)
		}
	}
	//更新购物车中的图书的总数量和总金额
	dao.UpdateCart(cart)
	//调用获取购物项信息的函数再次查询购物车信息
	GetCartInfo(w, r)
}
