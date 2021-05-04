package model

//Cart 购物车结构体
type Cart struct {
	CartID      string      //购物车的id
	CartItems   []*CartItem //购物车中所有的购物项
	TotalCount  int64       //购物车中图书的总数量，通过计算得到
	TotalAmount float64     //购物车中图书的总金额，通过计算得到
	UserID      int         //当前购物车所属的用户
}

//GetTotalCount 获取购物车中图书的总数量
func (cart *Cart) GetTotalCount() int64 {
	var totalCount int64
	//遍历购物车中的购物项切片
	for _, v := range cart.CartItems {
		totalCount = totalCount + v.Count
	}
	return totalCount
}

//GetTotalAmount 获取购物车中图书的总金额
func (cart *Cart) GetTotalAmount() float64 {
	var totalAmount float64
	//遍历购物车中的购物项切片
	for _, v := range cart.CartItems {
		totalAmount = totalAmount + v.GetAmount()
	}
	return totalAmount

}
