package dao

import (
	"model"
	"utils"
)

//AddOrder 向数据库中插入订单
func AddOrder(order *model.Order) error {
	//写sql语句
	sql := "insert into orders(id,create_time,total_count,total_amount,state,user_id) values(?,?,?,?,?,?)"
	//执行
	_, err := utils.Db.Exec(sql, order.OrderID, order.CreateTime, order.TotalCount, order.TotalAmount, order.State, order.UserID)
	if err != nil {
		return err
	}
	return nil
}

//GetOrders 获取数据库中所有的订单
func GetOrders() ([]*model.Order, error) {
	//写sql语句
	sql := "select id,create_time,total_count,total_amount,state,user_id from orders"
	//执行
	rows, err := utils.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		rows.Scan(&order.OrderID, &order.CreateTime, &order.TotalCount, &order.TotalAmount, &order.State, &order.UserID)
		orders = append(orders, order)
	}
	return orders, nil
}

//GetMyOrders 获取我的订单
func GetMyOrders(userID int) ([]*model.Order, error) {
	//写sql语句
	sql := "select id,create_time,total_count,total_amount,state,user_id from orders where user_id = ?"
	//执行
	rows, err := utils.Db.Query(sql, userID)
	if err != nil {
		return nil, err
	}
	//声明一个切片
	var orders []*model.Order
	for rows.Next() {
		//创建Order
		order := &model.Order{}
		//给Order中的字段赋值
		rows.Scan(&order.OrderID, &order.CreateTime, &order.TotalCount, &order.TotalAmount, &order.State, &order.UserID)
		//将Order添加到切片中
		orders = append(orders, order)
	}
	return orders, nil
}

//UpdateOrderState 更新订单的状态，即发货和收货
func UpdateOrderState(orderID string, state int64) error {
	//写sql语句
	sql := "update orders set state = ? where id = ?"
	//执行
	_, err := utils.Db.Exec(sql, state, orderID)
	if err != nil {
		return err
	}
	return nil
}
