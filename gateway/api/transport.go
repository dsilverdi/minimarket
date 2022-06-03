package api

import (
	"encoding/json"
	"fmt"
	"minimarket/gateway"
	"minimarket/pkg/rest"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

var rediscl = redis.NewClient(&redis.Options{
	Addr: "rediscache:6379",
})

var JWT_SECRET_KEY = []byte("minimarket-signature-key")

func NewHttpAPI(svc gateway.Service) *echo.Echo {
	e := echo.New()

	ep := NewEndpoint(svc)

	e.POST("/user/register", ep.RegisterEndpoint) //done

	e.POST("/user/authorize", ep.AuthorizeEndpoint) //done

	e.GET("/product", productCache(ep.GetProductEndpoint)) //done

	e.POST("/product/comment", AuthMiddleware(ep.WriteCommentEndpoint)) //done

	return e
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get("Authorization")

		if !strings.Contains(authorizationHeader, "Bearer") {
			return c.String(http.StatusUnauthorized, "Unauthorized User")
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("signing method invalid")
			}

			return JWT_SECRET_KEY, nil
		})

		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		c.Set("userinfo", claims)

		return next(c)
	}
}

func productCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		urlkey := c.Request().URL.String()

		val, err := rediscl.Get(ctx, urlkey).Bytes()
		if err != nil {
			return next(c)
		}

		data := toJson(val)

		return c.JSON(http.StatusOK, rest.HTTPResponse{
			Message: "cached",
			// PerPage: querypage.Limit,
			// Page:    querypage.Page,
			Data: data,
		})
	}
}

func toJson(val []byte) []gateway.Product {
	var prod []gateway.Product
	err := json.Unmarshal(val, &prod)
	if err != nil {
		panic(err)
	}
	return prod
}
