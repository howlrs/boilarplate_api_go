package routes

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *Client) Signup(c echo.Context) error {
	// リクエストのバインド
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
	}

	// パスワードのハッシュ化
	if err := user.ToEncryptPassword(); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to encrypt password")
	}

	// set database
	user.ID = xid.New().String()
	if _, err := p.firestore.Collection(user.ToCollection(p.IsTest())).Doc(user.ID).Set(c.Request().Context(), user); err != nil {
		// すでにKeyが存在する場合はエラーを返す
		if status.Code(err) == codes.AlreadyExists {
			return responseHandler(c, http.StatusConflict, nil, err, "Failed to set user, already exists")
		}

		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to set user")
	}

	return responseHandler(c, http.StatusOK, user, nil, "success, create user")
}

func (p *Client) Signin(c echo.Context) error {
	user := &models.User{}
	dbUser := &models.User{}
	if !p.IsTest() {
		// リクエストのバインド
		if err := c.Bind(user); err != nil {
			return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
		}
		if user.Email == "" || user.Password == "" {
			return responseHandler(c, http.StatusBadRequest, nil, nil, "Failed to bind request")
		}

		// read database
		// KeyはEmailを想定
		doc, err := p.firestore.Collection(user.ToCollection(p.IsTest())).Doc(user.Email).Get(c.Request().Context())
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return responseHandler(c, http.StatusNotFound, nil, err, "Failed to get user, not exists")
			}

			return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get user")
		}

		if err := doc.DataTo(dbUser); err != nil {
			return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get user")
		}

		// パスワードの検証
		if err := dbUser.IsVerifyPassword(user.Password); err != nil {
			return responseHandler(c, http.StatusUnauthorized, nil, err, "Failed to verify password")
		}

		// [Important] パスワードは返さない
		dbUser.Password = ""
	} else {
		// テスト用のユーザ情報を設定
		dbUser = &models.User{
			ID:       "test",
			Email:    "",
			Password: "test",
		}
	}

	// 期限を指定し
	// ユーザ情報からJWTトークンを生成
	expireAt := time.Now().Add(time.Hour * 24 * 7)
	isAdmin := false
	token, err := models.NewClaims(dbUser, isAdmin, expireAt).ToJwtToken()
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to create token")
	}

	return responseHandler(c, http.StatusOK, echo.Map{
		"token":      token,
		"token_type": "bearer",
		"user":       dbUser,
	}, nil, "success, create jwt token")
}
