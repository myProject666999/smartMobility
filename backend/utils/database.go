package utils

import (
	"fmt"
	"log"

	"smartMobility/config"
	"smartMobility/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() error {
	dsn := config.GetDBConnString()
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	
	log.Println("数据库连接成功")
	
	if err := models.AutoMigrate(DB); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}
	
	log.Println("数据库迁移完成")
	
	if err := InitSeedData(); err != nil {
		log.Printf("初始化种子数据失败: %v", err)
	}
	
	return nil
}

func InitSeedData() error {
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount > 0 {
		return nil
	}
	
	hashedPassword, err := HashPassword("admin123")
	if err != nil {
		return err
	}
	
	defaultAdmin := models.Admin{
		Username: "admin",
		Password: hashedPassword,
		RealName: "超级管理员",
		Status: 1,
	}
	
	if err := DB.Create(&defaultAdmin).Error; err != nil {
		return err
	}
	
	defaultMenus := []models.Menu{
		{Name: "统计分析", Path: "/statistics", Icon: "BarChartOutlined", Component: "Statistics", Sort: 1},
		{Name: "用户管理", Path: "/users", Icon: "UserOutlined", Component: "UserManage", Sort: 2},
		{Name: "路线管理", Path: "/routes", Icon: "CarOutlined", Component: "RouteManage", Sort: 3},
		{Name: "车站管理", Path: "/stations", Icon: "EnvironmentOutlined", Component: "StationManage", Sort: 4},
		{Name: "车票管理", Path: "/tickets", Icon: "TicketOutlined", Component: "TicketManage", Sort: 5},
		{Name: "订单管理", Path: "/orders", Icon: "ShoppingCartOutlined", Component: "OrderManage", Sort: 6},
		{Name: "评价管理", Path: "/reviews", Icon: "StarOutlined", Component: "ReviewManage", Sort: 7},
		{Name: "通知管理", Path: "/notifications", Icon: "BellOutlined", Component: "NotificationManage", Sort: 8},
		{Name: "轮播图管理", Path: "/banners", Icon: "PictureOutlined", Component: "BannerManage", Sort: 9},
		{Name: "公告管理", Path: "/announcements", Icon: "ReadOutlined", Component: "AnnouncementManage", Sort: 10},
		{Name: "菜单管理", Path: "/menus", Icon: "MenuOutlined", Component: "MenuManage", Sort: 11},
	}
	
	for _, menu := range defaultMenus {
		if err := DB.Create(&menu).Error; err != nil {
			log.Printf("创建菜单失败: %v", err)
		}
	}
	
	log.Println("种子数据初始化完成")
	return nil
}
