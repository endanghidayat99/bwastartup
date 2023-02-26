package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMX0.3tDO2sf0eJuxK2SFuc2tY8Ue4fIYOYHMwqF2APskOY4")
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("ERROR")
		fmt.Println("ERROR")
	}

	if token.Valid {
		fmt.Println("VALID")
		fmt.Println("VALID")
		fmt.Println("VALID")
	}

	if !token.Valid {
		fmt.Println("INVALID")
		fmt.Println("INVALID")
		fmt.Println("INVALID")
	}

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)
	router.Run()
	//
	//fmt.Println("Connection to Database is good")
	//
	//var users []user.User
	//length := len(users)
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
