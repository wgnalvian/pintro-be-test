package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"wgnalvian.com/payment-server/config"
	"wgnalvian.com/payment-server/dto"
	"wgnalvian.com/payment-server/entity"
	"wgnalvian.com/payment-server/service"
	"wgnalvian.com/payment-server/utils"
)

type UserController struct {
	UserService *service.UserService
}

func (u *UserController) Register(c *gin.Context) {

	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data not valid"})
		return
	}

	// validate not null
	if registerRequest.Email == "" || registerRequest.Password == "" || registerRequest.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data not valid"})
		return
	}

	hash, err := utils.HashPassword(registerRequest.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Check email exist
	isExist := u.UserService.CheckIfEmailExist(registerRequest.Email)

	if isExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist"})
		return
	}

	user := entity.User{
		ID:       uuid.New().String(),
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
		Password: hash,
	}

	err = u.UserService.Register(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func (u *UserController) GetUserData(c *gin.Context) {
	user, err := u.UserService.GetUserByEmail(c.GetString("email"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func (u *UserController) GetTransactions(c *gin.Context) {
	user, err := u.UserService.GetUserByEmail(c.GetString("email"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	transations, err := u.UserService.GetTransactions(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transations})
}

func (u *UserController) Transfer(c *gin.Context) {
	user, err := u.UserService.GetUserByEmail(c.GetString("email"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var transferRequest dto.TransferRequest

	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data not valid", "detail": err.Error()})
		return
	}

	if transferRequest.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Value must be greater than 0"})
		return
	}

	if user.Balance < transferRequest.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Get user op
	userOp, err := u.UserService.GetUserById(transferRequest.To)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if userOp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = u.UserService.Transfer(user.ID, userOp.ID, transferRequest.Amount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer success"})
}

func (u *UserController) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data not valid"})
		return
	}

	user, err := u.UserService.Login(loginRequest.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &entity.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(config.LoadConfig().JWT_SECRET)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Login success"})
}

func (u *UserController) TokenTopUp(c *gin.Context) {
	var tokenRequest dto.TokenRequest
	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data not valid", "detail": err.Error()})
		return
	}

	email := c.GetString("email")
	user, err := u.UserService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user != nil {
		err = u.UserService.TopUp(user.ID, tokenRequest.Amount)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	snapc := snap.Client{}

	snapc.New(config.LoadConfig().MIDTRANS_KEY, midtrans.Sandbox)

	reqSnap := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  uuid.New().String(),
			GrossAmt: int64(tokenRequest.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "CLIENT",
			Email: email,
		},
	}

	snapResp, _ := snapc.CreateTransaction(reqSnap)
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"token": snapResp})
}
