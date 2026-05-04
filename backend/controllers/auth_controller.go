package controllers

import (
	"smartMobility/models"
	"smartMobility/utils"

	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func UserRegister(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var existingUser models.User
	if err := utils.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	if req.Phone != "" {
		var userByPhone models.User
		if err := utils.DB.Where("phone = ?", req.Phone).First(&userByPhone).Error; err == nil {
			utils.BadRequest(c, "手机号已被注册")
			return
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   1,
	}

	if err := utils.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "注册失败")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"real_name": user.RealName,
			"phone":    user.Phone,
			"email":    user.Email,
		},
	})
}

func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := utils.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.BadRequest(c, "账号已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"real_name": user.RealName,
			"phone":    user.Phone,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var admin models.Admin
	if err := utils.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if admin.Status != 1 {
		utils.BadRequest(c, "账号已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"admin": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"real_name": admin.RealName,
		},
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"real_name": user.RealName,
		"phone":    user.Phone,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"status":   user.Status,
	})
}

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if realName, ok := req["real_name"].(string); ok {
		updates["real_name"] = realName
	}
	if phone, ok := req["phone"].(string); ok {
		updates["phone"] = phone
	}
	if email, ok := req["email"].(string); ok {
		updates["email"] = email
	}
	if avatar, ok := req["avatar"].(string); ok {
		updates["avatar"] = avatar
	}

	if err := utils.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

func UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if role == "user" {
		var user models.User
		if err := utils.DB.First(&user, userID).Error; err != nil {
			utils.NotFound(c, "用户不存在")
			return
		}

		if !utils.CheckPassword(req.OldPassword, user.Password) {
			utils.BadRequest(c, "原密码错误")
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			utils.InternalServerError(c, "密码加密失败")
			return
		}

		if err := utils.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
			utils.InternalServerError(c, "更新密码失败")
			return
		}
	} else if role == "admin" {
		var admin models.Admin
		if err := utils.DB.First(&admin, userID).Error; err != nil {
			utils.NotFound(c, "管理员不存在")
			return
		}

		if !utils.CheckPassword(req.OldPassword, admin.Password) {
			utils.BadRequest(c, "原密码错误")
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			utils.InternalServerError(c, "密码加密失败")
			return
		}

		if err := utils.DB.Model(&admin).Update("password", hashedPassword).Error; err != nil {
			utils.InternalServerError(c, "更新密码失败")
			return
		}
	}

	utils.Success(c, gin.H{"message": "密码更新成功"})
}
