package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password    string         `json:"-" gorm:"size:255;not null"`
	RealName    string         `json:"real_name" gorm:"size:50"`
	Phone       string         `json:"phone" gorm:"uniqueIndex;size:20"`
	Email       string         `json:"email" gorm:"size:100"`
	Avatar      string         `json:"avatar" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1"` // 1:正常 0:禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Admin struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	RealName  string         `json:"real_name" gorm:"size:50"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Email     string         `json:"email" gorm:"size:100"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Station struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Address   string         `json:"address" gorm:"size:255"`
	Latitude  float64        `json:"latitude" gorm:"type:decimal(10,7)"`
	Longitude float64        `json:"longitude" gorm:"type:decimal(10,7)"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常 0:停用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Route struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	RouteNumber string        `json:"route_number" gorm:"size:50;not null"`
	StartStationID uint       `json:"start_station_id" gorm:"not null"`
	EndStationID uint         `json:"end_station_id" gorm:"not null"`
	StartStation *Station     `json:"start_station" gorm:"foreignKey:StartStationID"`
	EndStation *Station       `json:"end_station" gorm:"foreignKey:EndStationID"`
	Distance    float64        `json:"distance" gorm:"type:decimal(10,2)"`
	Duration    int            `json:"duration"` // 预计时长(分钟)
	Status      int            `json:"status" gorm:"default:1"` // 1:正常 0:停运
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type RouteStation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RouteID   uint      `json:"route_id" gorm:"not null;uniqueIndex:idx_route_station"`
	StationID uint      `json:"station_id" gorm:"not null;uniqueIndex:idx_route_station"`
	Sequence  int       `json:"sequence" gorm:"not null"`
	ArriveTime string    `json:"arrive_time" gorm:"size:10"`
	DepartTime string    `json:"depart_time" gorm:"size:10"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Ticket struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	RouteID     uint           `json:"route_id" gorm:"not null"`
	Route       *Route         `json:"route" gorm:"foreignKey:RouteID"`
	DepartDate  time.Time      `json:"depart_date" gorm:"type:date;not null"`
	DepartTime  string         `json:"depart_time" gorm:"size:10;not null"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	SeatsTotal  int            `json:"seats_total" gorm:"not null"`
	SeatsSold   int            `json:"seats_sold" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:在售 0:停售
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderNo     string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        *User          `json:"user" gorm:"foreignKey:UserID"`
	TicketID    uint           `json:"ticket_id" gorm:"not null"`
	Ticket      *Ticket        `json:"ticket" gorm:"foreignKey:TicketID"`
	Quantity    int            `json:"quantity" gorm:"not null"`
	TotalPrice  float64        `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Status      int            `json:"status" gorm:"default:0"` // 0:待支付 1:已支付 2:已取消 3:已完成
	PayTime     *time.Time     `json:"pay_time"`
	PayMethod   string         `json:"pay_method" gorm:"size:20"`
	Remark      string         `json:"remark" gorm:"size:255"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	Order     *Order         `json:"order" gorm:"foreignKey:OrderID"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      *User          `json:"user" gorm:"foreignKey:UserID"`
	RouteID   uint           `json:"route_id" gorm:"not null"`
	Route     *Route         `json:"route" gorm:"foreignKey:RouteID"`
	Rating    int            `json:"rating" gorm:"not null"` // 1-5星
	Content   string         `json:"content" gorm:"size:500"`
	Images    string         `json:"images" gorm:"size:1000"` // 图片URL,逗号分隔
	Status    int            `json:"status" gorm:"default:1"` // 1:显示 0:隐藏
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Announcement struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Type      int            `json:"type" gorm:"default:0"` // 0:普通公告 1:重要公告
	Status    int            `json:"status" gorm:"default:1"` // 1:显示 0:隐藏
	Views     int            `json:"views" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Type      int            `json:"type" gorm:"default:0"` // 0:系统通知 1:订单通知 2:活动通知
	IsRead    int            `json:"is_read" gorm:"default:0"` // 0:未读 1:已读
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Banner struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Image     string         `json:"image" gorm:"size:255;not null"`
	Link      string         `json:"link" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:显示 0:隐藏
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Menu struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ParentID  uint           `json:"parent_id" gorm:"default:0"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Path      string         `json:"path" gorm:"size:255"`
	Icon      string         `json:"icon" gorm:"size:100"`
	Component string         `json:"component" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:显示 0:隐藏
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Admin{},
		&Station{},
		&Route{},
		&RouteStation{},
		&Ticket{},
		&Order{},
		&Review{},
		&Announcement{},
		&Notification{},
		&Banner{},
		&Menu{},
	)
}
