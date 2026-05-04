package controllers

import (
	"smartMobility/models"
	"smartMobility/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	var total int64
	query := utils.DB.Model(&models.Station{})
	
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	
	query.Count(&total)

	var stations []models.Station
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&stations)

	utils.Page(c, stations, total, page, pageSize)
}

func GetStationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var station models.Station
	if err := utils.DB.First(&station, id).Error; err != nil {
		utils.NotFound(c, "车站不存在")
		return
	}

	utils.Success(c, station)
}

func CreateStation(c *gin.Context) {
	var station models.Station
	if err := c.ShouldBindJSON(&station); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Create(&station).Error; err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	utils.Success(c, station)
}

func UpdateStation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var station models.Station
	if err := utils.DB.First(&station, id).Error; err != nil {
		utils.NotFound(c, "车站不存在")
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
	if address, ok := req["address"].(string); ok {
		updates["address"] = address
	}
	if latitude, ok := req["latitude"].(float64); ok {
		updates["latitude"] = latitude
	}
	if longitude, ok := req["longitude"].(float64); ok {
		updates["longitude"] = longitude
	}
	if status, ok := req["status"].(float64); ok {
		updates["status"] = int(status)
	}

	if err := utils.DB.Model(&station).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, station)
}

func DeleteStation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Delete(&models.Station{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func GetRoutes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	startStation := c.Query("start_station")
	endStation := c.Query("end_station")

	var total int64
	query := utils.DB.Model(&models.Route{}).Preload("StartStation").Preload("EndStation")
	
	if keyword != "" {
		query = query.Where("name LIKE ? OR route_number LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startStation != "" {
		query = query.Joins("JOIN stations AS s1 ON routes.start_station_id = s1.id").Where("s1.name LIKE ?", "%"+startStation+"%")
	}
	if endStation != "" {
		query = query.Joins("JOIN stations AS s2 ON routes.end_station_id = s2.id").Where("s2.name LIKE ?", "%"+endStation+"%")
	}
	
	query.Count(&total)

	var routes []models.Route
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&routes)

	utils.Page(c, routes, total, page, pageSize)
}

func GetRouteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var route models.Route
	if err := utils.DB.Preload("StartStation").Preload("EndStation").First(&route, id).Error; err != nil {
		utils.NotFound(c, "路线不存在")
		return
	}

	var routeStations []models.RouteStation
	utils.DB.Where("route_id = ?", id).Order("sequence ASC").Preload("Station").Find(&routeStations)

	utils.Success(c, gin.H{
		"route":          route,
		"route_stations": routeStations,
	})
}

func CreateRoute(c *gin.Context) {
	var req struct {
		models.Route
		StationIDs []uint `json:"station_ids"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := utils.DB.Create(&req.Route).Error; err != nil {
		utils.InternalServerError(c, "创建失败")
		return
	}

	for i, stationID := range req.StationIDs {
		routeStation := models.RouteStation{
			RouteID:   req.Route.ID,
			StationID: stationID,
			Sequence:  i + 1,
		}
		utils.DB.Create(&routeStation)
	}

	utils.Success(c, req.Route)
}

func UpdateRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var route models.Route
	if err := utils.DB.First(&route, id).Error; err != nil {
		utils.NotFound(c, "路线不存在")
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
	if routeNumber, ok := req["route_number"].(string); ok {
		updates["route_number"] = routeNumber
	}
	if distance, ok := req["distance"].(float64); ok {
		updates["distance"] = distance
	}
	if duration, ok := req["duration"].(float64); ok {
		updates["duration"] = int(duration)
	}
	if status, ok := req["status"].(float64); ok {
		updates["status"] = int(status)
	}

	if err := utils.DB.Model(&route).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, route)
}

func DeleteRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	utils.DB.Where("route_id = ?", id).Delete(&models.RouteStation{})
	
	if err := utils.DB.Delete(&models.Route{}, id).Error; err != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

func SearchRoutes(c *gin.Context) {
	startStation := c.Query("start_station")
	endStation := c.Query("end_station")
	date := c.Query("date")

	if startStation == "" || endStation == "" {
		utils.BadRequest(c, "请输入起点和终点")
		return
	}

	var routes []models.Route
	query := utils.DB.Model(&models.Route{}).
		Preload("StartStation").Preload("EndStation").
		Joins("JOIN stations AS s1 ON routes.start_station_id = s1.id").
		Joins("JOIN stations AS s2 ON routes.end_station_id = s2.id").
		Where("s1.name LIKE ? AND s2.name LIKE ? AND routes.status = 1", "%"+startStation+"%", "%"+endStation+"%")

	query.Find(&routes)

	var routeIDs []uint
	for _, r := range routes {
		routeIDs = append(routeIDs, r.ID)
	}

	var tickets []models.Ticket
	if len(routeIDs) > 0 {
		ticketQuery := utils.DB.Model(&models.Ticket{}).
			Preload("Route").Preload("Route.StartStation").Preload("Route.EndStation").
			Where("route_id IN ? AND status = 1", routeIDs)
		
		if date != "" {
			ticketQuery = ticketQuery.Where("depart_date = ?", date)
		}
		
		ticketQuery.Find(&tickets)
	}

	utils.Success(c, gin.H{
		"routes":  routes,
		"tickets": tickets,
	})
}
