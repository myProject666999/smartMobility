package routes

import (
	"smartMobility/controllers"
	"smartMobility/middleware"
	"smartMobility/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/user/register", controllers.UserRegister)
			auth.POST("/user/login", controllers.UserLogin)
			auth.POST("/admin/login", controllers.AdminLogin)
		}

		api.GET("/stations", controllers.GetStations)
		api.GET("/stations/:id", controllers.GetStationByID)
		api.GET("/routes", controllers.GetRoutes)
		api.GET("/routes/:id", controllers.GetRouteByID)
		api.GET("/routes/search", controllers.SearchRoutes)
		api.GET("/tickets", controllers.GetTickets)
		api.GET("/tickets/:id", controllers.GetTicketByID)
		api.GET("/announcements", controllers.GetAnnouncements)
		api.GET("/announcements/:id", controllers.GetAnnouncementByID)
		api.GET("/banners", controllers.GetBanners)
		api.GET("/reviews", controllers.GetReviews)

		user := api.Group("/user")
		user.Use(middleware.JWTAuth(), middleware.UserAuth())
		{
			user.GET("/info", controllers.GetUserInfo)
			user.PUT("/info", controllers.UpdateUserInfo)
			user.PUT("/password", controllers.UpdatePassword)

			user.GET("/notifications", controllers.GetNotifications)
			user.PUT("/notifications/:id/read", controllers.MarkNotificationRead)

			user.POST("/orders", controllers.CreateOrder)
			user.GET("/orders", controllers.GetOrders)
			user.GET("/orders/:id", controllers.GetOrderByID)
			user.POST("/orders/:id/pay", controllers.PayOrder)
			user.POST("/orders/:id/cancel", controllers.CancelOrder)
			user.POST("/orders/:id/refund", controllers.RefundOrder)

			user.POST("/reviews", controllers.CreateReview)
			user.GET("/reviews", controllers.GetReviews)

			user.GET("/ws", websocket.HandleWebSocket)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			admin.PUT("/password", controllers.UpdatePassword)

			admin.GET("/statistics", controllers.GetStatistics)

			admin.GET("/users", controllers.GetUserManagement)
			admin.PUT("/users/:id/status", controllers.UpdateUserStatus)

			admin.GET("/stations", controllers.GetStations)
			admin.GET("/stations/:id", controllers.GetStationByID)
			admin.POST("/stations", controllers.CreateStation)
			admin.PUT("/stations/:id", controllers.UpdateStation)
			admin.DELETE("/stations/:id", controllers.DeleteStation)

			admin.GET("/routes", controllers.GetRoutes)
			admin.GET("/routes/:id", controllers.GetRouteByID)
			admin.POST("/routes", controllers.CreateRoute)
			admin.PUT("/routes/:id", controllers.UpdateRoute)
			admin.DELETE("/routes/:id", controllers.DeleteRoute)

			admin.GET("/tickets", controllers.GetTickets)
			admin.GET("/tickets/:id", controllers.GetTicketByID)
			admin.POST("/tickets", controllers.CreateTicket)
			admin.PUT("/tickets/:id", controllers.UpdateTicket)
			admin.DELETE("/tickets/:id", controllers.DeleteTicket)

			admin.GET("/orders", controllers.GetOrders)
			admin.GET("/orders/:id", controllers.GetOrderByID)
			admin.POST("/orders/:id/cancel", controllers.CancelOrder)
			admin.POST("/orders/:id/refund", controllers.RefundOrder)

			admin.GET("/reviews", controllers.GetReviews)
			admin.PUT("/reviews/:id/status", controllers.UpdateReviewStatus)

			admin.GET("/announcements", controllers.GetAnnouncements)
			admin.GET("/announcements/:id", controllers.GetAnnouncementByID)
			admin.POST("/announcements", controllers.CreateAnnouncement)
			admin.PUT("/announcements/:id", controllers.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", controllers.DeleteAnnouncement)

			admin.GET("/banners", controllers.GetBanners)
			admin.POST("/banners", controllers.CreateBanner)
			admin.PUT("/banners/:id", controllers.UpdateBanner)
			admin.DELETE("/banners/:id", controllers.DeleteBanner)

			admin.GET("/menus", controllers.GetMenus)
			admin.POST("/menus", controllers.CreateMenu)
			admin.PUT("/menus/:id", controllers.UpdateMenu)
			admin.DELETE("/menus/:id", controllers.DeleteMenu)

			admin.GET("/ws", websocket.HandleWebSocket)
		}
	}
}
