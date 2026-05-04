package controllers

import (
	"smartMobility/models"
	"smartMobility/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetReviews(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	routeID := c.Query("route_id")

	var total int64
	query := utils.DB.Model(&models.Review{}).
		Preload("User").
		Preload("Order").
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation")

	if role == "user" {
		query = query.Where("user_id = ?", userID)
	}
	if routeID != "" {
		query = query.Where("route_id = ?", routeID)
	}

	query.Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews)

	utils.Page(c, reviews, total, page, pageSize)
}

func CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OrderID uint   `json:"order_id" binding:"required"`
		Rating  int    `json:"rating" binding:"required"`
		Content string `json:"content"`
		Images  string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		utils.BadRequest(c, "评分只能是1-5星")
		return
	}

	var order models.Order
	if err := utils.DB.First(&order, req.OrderID).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.UserID != userID {
		utils.Forbidden(c, "无权评价此订单")
		return
	}

	if order.Status != 1 && order.Status != 3 {
		utils.BadRequest(c, "只能对已支付的订单进行评价")
		return
	}

	var existingReview models.Review
	if err := utils.DB.Where("order_id = ?", req.OrderID).First(&existingReview).Error; err == nil {
		utils.BadRequest(c, "该订单已评价")
		return
	}

	var ticket models.Ticket
	if err := utils.DB.First(&ticket, order.TicketID).Error; err != nil {
		utils.NotFound(c, "车票不存在")
		return
	}

	review := models.Review{
		OrderID: req.OrderID,
		UserID:  userID,
		RouteID: ticket.RouteID,
		Rating:  req.Rating,
		Content: req.Content,
		Images:  req.Images,
		Status:  1,
	}

	if err := utils.DB.Create(&review).Error; err != nil {
		utils.InternalServerError(c, "评价失败")
		return
	}

	utils.Success(c, review)
}

func UpdateReviewStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var review models.Review
	if err := utils.DB.First(&review, id).Error; err != nil {
		utils.NotFound(c, "评价不存在")
		return
	}

	review.Status = req.Status
	if err := utils.DB.Save(&review).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

func GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.GetString("role")

	var total int64
	query := utils.DB.Model(&models.Announcement{})
	
	if role != "admin" {
		query = query.Where("status = 1")
	}
	
	query.Count(&total)

	var announcements []models.Announcement
	offset := (page - 1) * pageSize
	query.Order("type DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&announcements)

	utils.Page(c, announcements, total, page, pageSize)
}

func GetAnnouncementByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var announcement models.Announcement
	if err := utils.DB.First(&announcement, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	utils.DB.Model(&announcement).Update("views", announcement.Views+1)

	utils.Success(c, announcement)
}

func CreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Create(&announcement).Error; err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	utils.Success(c, announcement)
}

func UpdateAnnouncement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var announcement models.Announcement
	if err := utils.DB.First(&announcement, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if title, ok := req["title"].(string); ok {
		updates["title"] = title
	}
	if content, ok := req["content"].(string); ok {
		updates["content"] = content
	}
	if typeVal, ok := req["type"].(float64); ok {
		updates["type"] = int(typeVal)
	}
	if status, ok := req["status"].(float64); ok {
		updates["status"] = int(status)
	}

	if err := utils.DB.Model(&announcement).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, announcement)
}

func DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Delete(&models.Announcement{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	isRead := c.Query("is_read")

	var total int64
	query := utils.DB.Model(&models.Notification{}).Where("user_id = ?", userID)
	
	if isRead != "" {
		query = query.Where("is_read = ?", isRead)
	}
	
	query.Count(&total)

	var notifications []models.Notification
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications)

	utils.Page(c, notifications, total, page, pageSize)
}

func MarkNotificationRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var notification models.Notification
	if err := utils.DB.First(&notification, id).Error; err != nil {
		utils.NotFound(c, "通知不存在")
		return
	}

	if notification.UserID != userID {
		utils.Forbidden(c, "无权操作此通知")
		return
	}

	notification.IsRead = 1
	if err := utils.DB.Save(&notification).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "已标记为已读"})
}

func GetBanners(c *gin.Context) {
	role := c.GetString("role")

	var banners []models.Banner
	query := utils.DB.Model(&models.Banner{})
	
	if role != "admin" {
		query = query.Where("status = 1")
	}
	
	query.Order("sort ASC, created_at DESC").Find(&banners)

	utils.Success(c, banners)
}

func CreateBanner(c *gin.Context) {
	var banner models.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Create(&banner).Error; err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	utils.Success(c, banner)
}

func UpdateBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var banner models.Banner
	if err := utils.DB.First(&banner, id).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if title, ok := req["title"].(string); ok {
		updates["title"] = title
	}
	if image, ok := req["image"].(string); ok {
		updates["image"] = image
	}
	if link, ok := req["link"].(string); ok {
		updates["link"] = link
	}
	if sort, ok := req["sort"].(float64); ok {
		updates["sort"] = int(sort)
	}
	if status, ok := req["status"].(float64); ok {
		updates["status"] = int(status)
	}

	if err := utils.DB.Model(&banner).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, banner)
}

func DeleteBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Delete(&models.Banner{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func GetMenus(c *gin.Context) {
	var menus []models.Menu
	utils.DB.Order("sort ASC, created_at ASC").Find(&menus)

	utils.Success(c, menus)
}

func CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Create(&menu).Error; err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	utils.Success(c, menu)
}

func UpdateMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var menu models.Menu
	if err := utils.DB.First(&menu, id).Error; err != nil {
		utils.NotFound(c, "菜单不存在")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if name, ok := req["name"].(string); ok {
		updates["name"] = name
	}
	if path, ok := req["path"].(string); ok {
		updates["path"] = path
	}
	if icon, ok := req["icon"].(string); ok {
		updates["icon"] = icon
	}
	if component, ok := req["component"].(string); ok {
		updates["component"] = component
	}
	if sort, ok := req["sort"].(float64); ok {
		updates["sort"] = int(sort)
	}
	if status, ok := req["status"].(float64); ok {
		updates["status"] = int(status)
	}

	if err := utils.DB.Model(&menu).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, menu)
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Delete(&models.Menu{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func GetStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var ticketCount int64
	var ticketRevenue float64
	ticketQuery := utils.DB.Model(&models.Order{}).Where("status = 1")
	
	if startDate != "" {
		ticketQuery = ticketQuery.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		ticketQuery = ticketQuery.Where("created_at <= ?", endDate+" 23:59:59")
	}
	
	ticketQuery.Count(&ticketCount)
	ticketQuery.Select("COALESCE(SUM(total_price), 0)").Scan(&ticketRevenue)

	var userCount int64
	utils.DB.Model(&models.User{}).Count(&userCount)

	var orderCount int64
	utils.DB.Model(&models.Order{}).Count(&orderCount)

	var pendingOrders int64
	utils.DB.Model(&models.Order{}).Where("status = 0").Count(&pendingOrders)

	utils.Success(c, gin.H{
		"ticket_count":      ticketCount,
		"ticket_revenue":    ticketRevenue,
		"user_count":        userCount,
		"order_count":       orderCount,
		"pending_orders":    pendingOrders,
	})
}

func GetUserManagement(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	var total int64
	query := utils.DB.Model(&models.User{})
	
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	utils.Page(c, users, total, page, pageSize)
}

func UpdateUserStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := utils.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	user.Status = req.Status
	if err := utils.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}
