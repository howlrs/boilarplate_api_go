package routes

import (
	"backend/models"
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Client struct {
	isTest    bool
	firestore *firestore.Client
}

func NewClient(isTest bool) *Client {
	ctx := context.Background()
	projectId := os.Getenv("PROJECT_ID")
	db, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		panic(err)
	}

	return &Client{
		isTest:    isTest,
		firestore: db,
	}
}

// IsTest: テストモードかどうかを返す
// コレクション名を切り替える際やテスト用のデータを返す際に使用
func (p *Client) IsTest() bool {
	return p.isTest
}

// Endpoint sets up the routes for the application.
func Endpoint(e *echo.Echo, isTest bool) {
	p := NewClient(isTest)

	// version 1
	v1 := e.Group("/api/v1")

	// public routes
	v1Public := v1.Group("/public")
	v1Public.GET("/health", publicHealth)
	v1Public.POST("/signup", p.Signup)
	v1Public.POST("/signin", p.Signin)
	v1Public.GET("/reservation", p.ReadReservation)
	v1Public.POST("/reservation", p.CreateReservation)
	v1Public.PUT("/reservation", p.UpdateReservation)
	v1Public.DELETE("/reservation", p.DeleteReservation)

	// private routes
	v1Private := v1.Group("/private")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.Claims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	v1Private.Use(echojwt.WithConfig(config))
	// -H "Authorization: Bearer <token>"を付与してリクエスト
	v1Private.GET("/health", privateHealth)
}

func publicHealth(c echo.Context) error {
	return responseHandler(c, http.StatusOK, echo.Map{"message": "OK"}, nil, "success, public health")
}

func privateHealth(c echo.Context) error {
	// コンテキストからユーザ情報を取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.Claims)
	fmt.Println(claims)

	return responseHandler(c, http.StatusOK, echo.Map{"message": "OK"}, nil, "success, private health")
}
