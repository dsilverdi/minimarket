package api

import (
	"encoding/json"
	"log"
	"minimarket/gateway"
	"minimarket/pkg/rest"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ServerEndpoint struct {
	RegisterEndpoint     echo.HandlerFunc
	AuthorizeEndpoint    echo.HandlerFunc
	GetProductEndpoint   echo.HandlerFunc
	WriteCommentEndpoint echo.HandlerFunc
}

func NewEndpoint(svc gateway.Service) ServerEndpoint {
	return ServerEndpoint{
		RegisterEndpoint:     RegisterEndpoint(svc),
		AuthorizeEndpoint:    AuthorizeEndpoint(svc),
		GetProductEndpoint:   GetProductEndpoint(svc),
		WriteCommentEndpoint: WriteCommentEndpoint(svc),
	}
}

func RegisterEndpoint(svc gateway.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var payload UserRequestBody
		if err := c.Bind(&payload); err != nil {
			return err
		}

		if payload.Email == "" || payload.Password == "" {
			return c.JSON(http.StatusBadRequest, rest.HTTPResponse{
				Message: "empty payload received",
			})
		}

		err := svc.Register(ctx, payload.Email, payload.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, rest.HTTPResponse{
				Message: "error",
				Errors:  err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, rest.HTTPResponse{
			Message: "Success",
			Data: map[string]interface{}{
				"email": payload.Email,
			},
		})
	}
}

func AuthorizeEndpoint(svc gateway.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var payload UserRequestBody
		if err := c.Bind(&payload); err != nil {
			return err
		}

		if payload.Email == "" || payload.Password == "" {
			return c.JSON(http.StatusBadRequest, rest.HTTPResponse{
				Message: "empty payload received",
			})
		}

		val, err := rediscl.Get(ctx, payload.Email).Bytes()
		if err != nil {

			auth, err := svc.Authorize(ctx, payload.Email, payload.Password)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, rest.HTTPResponse{
					Message: "error",
					Errors:  err.Error(),
				})
			}

			err = rediscl.Set(ctx, payload.Email, auth.Token, 2*24*time.Hour).Err()
			if err != nil {
				log.Print(err)
			}

			return c.JSON(http.StatusOK, rest.HTTPResponse{
				Message: "Success",
				Data: map[string]interface{}{
					"token": auth.Token,
				},
			})
		}

		return c.JSON(http.StatusOK, rest.HTTPResponse{
			Message: "cached",
			Data: map[string]interface{}{
				"token": string(val),
			},
		})

	}
}

func GetProductEndpoint(svc gateway.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		urlkey := c.Request().URL.String()

		querypage := rest.QueryParameter{}

		per_page := c.QueryParam("limit")
		page := c.QueryParam("page")

		var err error
		if per_page == "" {
			querypage.Limit = 15
		} else {
			querypage.Limit, err = strconv.Atoi(per_page)
			if err != nil {
				return err
			}
		}

		if page == "" {
			querypage.Page = 1
		} else {
			querypage.Page, err = strconv.Atoi(page)
			if err != nil {
				return err
			}
		}

		name := c.QueryParam("name")
		category := c.QueryParam("category")

		prod, err := svc.GetProducts(ctx, querypage, name, category)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, rest.HTTPResponse{
				Message: "error",
				Errors:  err.Error(),
			})
		}

		buf, err := json.Marshal(prod)
		if err != nil {
			panic(err)
		}

		err = rediscl.Set(ctx, urlkey, buf, 5*time.Minute).Err()
		if err != nil {
			log.Print(err)
		}

		return c.JSON(http.StatusOK, rest.HTTPResponse{
			Message: "success",
			PerPage: querypage.Limit,
			Page:    querypage.Page,
			Data:    prod,
		})
	}
}

func WriteCommentEndpoint(svc gateway.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userinfo := c.Get("userinfo")

		var userdata UserClaim
		bodyBytes, _ := json.Marshal(userinfo)
		json.Unmarshal(bodyBytes, &userdata)

		var payload CommentRequestBody
		if err := c.Bind(&payload); err != nil {
			return err
		}

		if payload.ProductID == 0 || payload.Message == "" {
			return c.JSON(http.StatusBadRequest, rest.HTTPResponse{
				Message: "empty payload received",
			})
		}

		err := svc.WriteComment(ctx, payload.ProductID, payload.ReplyID, payload.Message, userdata.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, rest.HTTPResponse{
				Message: "error",
				Errors:  err.Error(),
			})
		}

		return c.JSON(http.StatusOK, rest.HTTPResponse{
			Message: "Success Sending Message",
			Data: map[string]interface{}{
				"product_id": payload.ProductID,
				"message":    payload.Message,
				"reply_id":   payload.ReplyID,
				"owner":      userdata.Email,
			},
		})
	}
}

type UserClaim struct {
	Email string
	Exp   int64
}
