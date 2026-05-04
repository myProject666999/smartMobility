package controllers

import (
	"fmt"
	"smartMobility/models"
	"smartMobility/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	routeID := c.Query("route_id")
	date := c.Query("date")

	var total int64
	query := utils.DB.Model(&models.Ticket{}).Preload("Route").Preload("Route.StartStation").Preload("Route.EndStation")
	
	if routeID != "" {
		query = query.Where("route_id = ?", routeID)
	}
	if date != "" {
		query = query.Where("depart_date = ?", date)
	}
	
	query.Count(&total)

	var tickets []models.Ticket
	offset := (page - 1) * pageSize
	query.Order("depart_date DESC, depart_time ASC").Offset(offset).Limit(pageSize).Find(&tickets)

	utils.Page(c, tickets, total, page, pageSize)
}

func GetTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var ticket models.Ticket
	if err := utils.DB.Preload("Route").Preload("Route.StartStation").Preload("Route.EndStation").First(&ticket, id).Error; err != nil {
		utils.NotFound(c, "车票不存在")
		return
	}

	utils.Success(c, ticket)
}

func CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Create(&ticket).Error; err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	utils.Success(c, ticket)
}

func UpdateTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var ticket models.Ticket
	if err := utils.DB.First(&ticket, id).Error; err != nil {
		utils.NotFound(c, "车票不存在")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if price, ok := req["price"].(float64); ok {
		updates["price"] = price
	}
	if seatsTotal, ok := req["seats_total"].(float64); ok {
		updates["seats_total"] = int(seatsTotal)
	}
	if seatsSold, ok := req["seats_sold"].(float64); ok {
		updates["seats_sold"] = int(seatsSold)
	}
	if status, ok := req["status"].(float64); ok {
		updates["status"] = int(status)
	}

	if err := utils.DB.Model(&ticket).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, ticket)
}

func DeleteTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Delete(&models.Ticket{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		TicketID uint `json:"ticket_id" binding:"required"`
		Quantity int  `json:"quantity" binding:"required"`
		Remark   string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.Quantity <= 0 {
		utils.BadRequest(c, "购票数量无效")
		return
	}

	var ticket models.Ticket
	if err := utils.DB.First(&ticket, req.TicketID).Error; err != nil {
		utils.NotFound(c, "车票不存在")
		return
	}

	if ticket.Status != 1 {
		utils.BadRequest(c, "该车票已停售")
		return
	}

	if ticket.SeatsSold+req.Quantity > ticket.SeatsTotal {
		utils.BadRequest(c, "余票不足")
		return
	}

	orderNo := fmt.Sprintf("OD%s%d", time.Now().Format("20060102150405"), userID)

	order := models.Order{
		OrderNo:    orderNo,
		UserID:     userID,
		TicketID:   req.TicketID,
		Quantity:   req.Quantity,
		TotalPrice: float64(req.Quantity) * ticket.Price,
		Status:     0,
		Remark:     req.Remark,
	}

	if err := utils.DB.Create(&order).Error; err != nil {
		utils.InternalServerError(c, "创建订单失败")
		return
	}

	utils.Success(c, gin.H{
		"order_id": order.ID,
		"order_no": order.OrderNo,
		"total_price": order.TotalPrice,
	})
}

func GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var total int64
	query := utils.DB.Model(&models.Order{}).
		Preload("User").
		Preload("Ticket").
		Preload("Ticket.Route").
		Preload("Ticket.Route.StartStation").
		Preload("Ticket.Route.EndStation")

	if role == "user" {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Page(c, orders, total, page, pageSize)
}

func GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := utils.DB.Preload("User").
		Preload("Ticket").
		Preload("Ticket.Route").
		Preload("Ticket.Route.StartStation").
		Preload("Ticket.Route.EndStation").
		First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func PayOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := utils.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.UserID != userID {
		utils.Forbidden(c, "无权操作此订单")
		return
	}

	if order.Status != 0 {
		utils.BadRequest(c, "订单状态错误")
		return
	}

	var ticket models.Ticket
	if err := utils.DB.First(&ticket, order.TicketID).Error; err != nil {
		utils.NotFound(c, "车票不存在")
		return
	}

	if ticket.SeatsSold+order.Quantity > ticket.SeatsTotal {
		utils.BadRequest(c, "余票不足")
		return
	}

	now := time.Now()
	order.Status = 1
	order.PayTime = &now
	order.PayMethod = "alipay"

	ticket.SeatsSold += order.Quantity

	tx := utils.DB.Begin()
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "支付失败")
		return
	}

	if err := tx.Save(&ticket).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "支付失败")
		return
	}

	notification := models.Notification{
		UserID:  userID,
		Title:   "订单支付成功",
		Content: fmt.Sprintf("您的订单 %s 已支付成功，共支付 %.2f 元", order.OrderNo, order.TotalPrice),
		Type:    1,
		IsRead:  0,
	}
	tx.Create(&notification)

	tx.Commit()

	utils.Success(c, gin.H{"message": "支付成功"})
}

func CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := utils.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if role == "user" && order.UserID != userID {
		utils.Forbidden(c, "无权操作此订单")
		return
	}

	if order.Status != 0 {
		utils.BadRequest(c, "只能取消待支付的订单")
		return
	}

	order.Status = 2

	if err := utils.DB.Save(&order).Error; err != nil {
		utils.InternalServerError(c, "取消订单失败")
		return
	}

	utils.Success(c, gin.H{"message": "订单已取消"})
}

func RefundOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := utils.DB.Preload("Ticket").First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if role == "user" && order.UserID != userID {
		utils.Forbidden(c, "无权操作此订单")
		return
	}

	if order.Status != 1 {
		utils.BadRequest(c, "只能对已支付的订单申请退款")
		return
	}

	tx := utils.DB.Begin()

	order.Status = 2
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "退款失败")
		return
	}

	if order.Ticket != nil {
		order.Ticket.SeatsSold -= order.Quantity
		if order.Ticket.SeatsSold < 0 {
			order.Ticket.SeatsSold = 0
		}
		if err := tx.Save(order.Ticket).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "退款失败")
			return
		}
	}

	tx.Commit()

	utils.Success(c, gin.H{"message": "退款成功"})
}
