package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

//func reassignedPriorities(priorities []int32) []int32 {
//	// Write your code here
//	hashMap := make(map[int32]int32)
//	n := len(priorities)
//	tempArr := make([]int32, n)
//	for i := 0; i < n; i++ {
//		tempArr[i] = priorities[i]
//	}
//	sort.Slice(tempArr, func(i, j int) bool { return tempArr[i] < tempArr[j] })
//	var priority int32
//	priority = 1
//	for i := 0; i < n; i++ {
//		if i > 0 {
//			if tempArr[i-1] != tempArr[i] {
//				priority += 1
//			}
//		}
//		hashMap[tempArr[i]] = priority
//	}
//
//	for i := 0; i < n; i++ {
//		priorities[i] = hashMap[priorities[i]]
//	}
//
//	return priorities
//}
//
//func frequencyOfMaxValue(numbers []int32, q []int32) []int32 {
//	n := len(numbers)
//	answer := make([]int32, n)
//	for i := 0; i < n; i++ {
//		answer[i] = -1
//	}
//	maxvalue, count := -1, 1
//	for i := n - 1; i >= 0; i-- {
//		if int(numbers[i]) == maxvalue {
//			count += 1
//		}
//		if int(numbers[i]) > maxvalue {
//			maxvalue = int(numbers[i])
//			count = 1
//		}
//		answer[i] = int32((count))
//	}
//	return answer
//}
//
//func main() {
//	numbers := []int32{2, 1, 2}
//	q := []int32{1, 2, 3}
//
//	fmt.Println(frequencyOfMaxValue(numbers, q))
//}

func main() {
	dsn := "sql12610918:CxakkFFBWc@tcp(sql12.freemysqlhosting.net:3306)/sql12610918?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	//service
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	//handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	//router
	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaignById)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransaction)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransactions)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.Run()
	//
	//fmt.Println("Connection to Database is good")
	//
	//var users []user.User
	//length := len(users)nsom
	//fmt.Println(length)
	//db.Find(&users)
	//length = len(users)
	//fmt.Println(length)
	//
	//for _, user := range users {
	//	fmt.Println(user.Name)
	//	fmt.Println(user.Email)
	//	fmt.Println("=============")
	//}

	//router := gin.Default()
	//router.GET("/handler", handler)
	//
	//router.Run()
}

//
//func handler(c *gin.Context) {
//	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	var users []user.User
//	db.Find(&users)
//
//	c.JSON(http.StatusOK, users)
//}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
