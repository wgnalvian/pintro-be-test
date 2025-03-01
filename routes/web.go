package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"wgnalvian.com/payment-server/config"
	"wgnalvian.com/payment-server/controller"
	"wgnalvian.com/payment-server/entity"
)

type Route struct {
	UserController *controller.UserController
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Split "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1] // Ambil tokennya saja

		jwtKey := []byte(config.LoadConfig().JWT_SECRET)
		claims := &entity.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Simpan email di context untuk digunakan di endpoint lain
		c.Set("email", claims.Email)
		c.Next()
	}
}

func (r *Route) Init() *gin.Engine {
	// Initialize the routes
	rt := gin.Default()

	rt.POST("/register", r.UserController.Register)
	rt.POST("/login", r.UserController.Login)
	rt.GET("/user", AuthMiddleware(), r.UserController.GetUserData)
	rt.GET("/transaction", AuthMiddleware(), r.UserController.GetTransactions)
	rt.POST("/transfer", AuthMiddleware(), r.UserController.Transfer)
	rt.POST("/token", AuthMiddleware(), r.UserController.TokenTopUp)
	return rt
}
