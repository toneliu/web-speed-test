package handlers

import (
	"net/http"
	"speedtest/pkg/database"
	"speedtest/pkg/middleware"
	"speedtest/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).Preload("Unit").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !database.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"is_admin":  user.IsAdmin,
			"unit_id":   user.UnitID,
			"unit":      user.Unit,
		},
	})
}

func GetUnits(c *gin.Context) {
	var units []models.Unit
	database.DB.Find(&units)
	c.JSON(http.StatusOK, units)
}

func CreateUnit(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unit := models.Unit{Name: req.Name}
	if err := database.DB.Create(&unit).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Unit already exists"})
		return
	}

	c.JSON(http.StatusOK, unit)
}

func UpdateUnit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var unit models.Unit
	if err := database.DB.First(&unit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	unit.Name = req.Name
	database.DB.Save(&unit)
	c.JSON(http.StatusOK, unit)
}

func DeleteUnit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&models.Unit{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Unit deleted"})
}

func GetUnitUsers(c *gin.Context) {
	unitID, _ := strconv.Atoi(c.Param("id"))
	var users []models.User
	database.DB.Where("unit_id = ?", unitID).Find(&users)
	c.JSON(http.StatusOK, users)
}

func CreateUnitUser(c *gin.Context) {
	unitID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	uid := uint(unitID)
	user := models.User{
		Username: req.Username,
		Password: string(hashedPwd),
		IsAdmin:  false,
		UnitID:   &uid,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ResetUserPassword(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPwd)
	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset"})
}

func DeleteUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	database.DB.Delete(&models.User{}, userID)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func SubmitSpeedTest(c *gin.Context) {
	var req struct {
		UnitID   uint    `json:"unit_id"`
		Download float64 `json:"download" binding:"required"`
		Upload   float64 `json:"upload"`
		Ping     float64 `json:"ping"`
		Jitter   float64 `json:"jitter"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var unitID uint
	isAdminVal, isAdminExists := c.Get("is_admin")
	userUnitIDInterface, userUnitIDExists := c.Get("unit_id")

	isAdmin := false
	if isAdminExists && isAdminVal != nil {
		if adminBool, ok := isAdminVal.(bool); ok {
			isAdmin = adminBool
		}
	}

	// 确定要使用的单位 ID
	if isAdmin && req.UnitID > 0 {
		// 管理员可以指定单位
		unitID = req.UnitID
	} else {
		// 普通用户使用自己的单位
		if userUnitIDExists && userUnitIDInterface != nil {
			if userUnitID, ok := userUnitIDInterface.(*uint); ok && userUnitID != nil {
				unitID = *userUnitID
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User not associated with a unit"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not associated with a unit"})
			return
		}
	}

	test := models.SpeedTest{
		UnitID:    unitID,
		Download:  req.Download,
		Upload:    req.Upload,
		Ping:      req.Ping,
		Jitter:    req.Jitter,
		IP:         c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if err := database.DB.Create(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save test"})
		return
	}

	database.DB.Preload("Unit").First(&test, test.ID)
	c.JSON(http.StatusOK, test)
}

func GetSpeedTests(c *gin.Context) {
	isAdminVal, isAdminExists := c.Get("is_admin")
	unitIDInterface, unitIDExists := c.Get("unit_id")

	isAdmin := false
	if isAdminExists && isAdminVal != nil {
		if adminBool, ok := isAdminVal.(bool); ok {
			isAdmin = adminBool
		}
	}

	var tests []models.SpeedTest
	query := database.DB.Preload("Unit").Order("created_at DESC")

	// 首先检查是否有查询参数指定 unit_id
	if unitIDQuery := c.Query("unit_id"); unitIDQuery != "" {
		qUnitID, _ := strconv.Atoi(unitIDQuery)
		// 只有管理员或者查询的是自己的单位
		if isAdmin {
			query = query.Where("unit_id = ?", qUnitID)
		} else if unitIDExists && unitIDInterface != nil {
			if userUnitID, ok := unitIDInterface.(*uint); ok && userUnitID != nil {
				if *userUnitID == uint(qUnitID) {
					query = query.Where("unit_id = ?", qUnitID)
				} else {
					// 非管理员不能查询其他单位
					query = query.Where("1=0") // 不返回任何记录
				}
			} else {
				query = query.Where("1=0")
			}
		} else {
			query = query.Where("1=0") // 不返回任何记录
		}
	} else if !isAdmin && unitIDExists && unitIDInterface != nil {
		// 非管理员默认只看自己单位的记录
		if userUnitID, ok := unitIDInterface.(*uint); ok && userUnitID != nil {
			query = query.Where("unit_id = ?", *userUnitID)
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	query.Limit(limit).Find(&tests)
	c.JSON(http.StatusOK, tests)
}

func GetTopologyLinks(c *gin.Context) {
	var links []models.TopologyLink
	database.DB.Preload("FromUnit").Preload("ToUnit").Find(&links)

	var units []models.Unit
	database.DB.Find(&units)

	c.JSON(http.StatusOK, gin.H{
		"links": links,
		"units": units,
	})
}

func CreateTopologyLink(c *gin.Context) {
	var req struct {
		FromUnitID uint    `json:"from_unit_id" binding:"required"`
		ToUnitID   uint    `json:"to_unit_id" binding:"required"`
		Bandwidth  float64 `json:"bandwidth"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link := models.TopologyLink{
		FromUnitID: req.FromUnitID,
		ToUnitID:   req.ToUnitID,
		Bandwidth:  req.Bandwidth,
	}

	database.DB.Create(&link)
	database.DB.Preload("FromUnit").Preload("ToUnit").First(&link, link.ID)
	c.JSON(http.StatusOK, link)
}

func DeleteTopologyLink(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&models.TopologyLink{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Link deleted"})
}

func GetStats(c *gin.Context) {
	var totalUnits int64
	database.DB.Model(&models.Unit{}).Count(&totalUnits)

	var totalTests int64
	database.DB.Model(&models.SpeedTest{}).Count(&totalTests)

	var recentTest models.SpeedTest
	database.DB.Preload("Unit").Order("created_at DESC").First(&recentTest)

	c.JSON(http.StatusOK, gin.H{
		"total_units": totalUnits,
		"total_tests": totalTests,
		"recent_test": recentTest,
	})
}
