package rest

import (
	"context"
	"net/http"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"

	"github.com/gin-gonic/gin"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, user *entity.Users) error
	LoginUser(ctx context.Context, user *entity.Login) error
}

type UserController struct {
	uuc UserUsecase
}

func NewUserController(uuc UserUsecase) *UserController {
	return &UserController{uuc}
}

func (uh *UserController) Register(ctx *gin.Context) {
	var payload dto.RegisterReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	data := entity.Users{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := uh.uuc.RegisterUser(ctx, &data); err != nil {
		ctx.Error(err)
		return
	}

	res := dto.RegisterRes{
		ID:    data.ID,
		Role:  data.Role,
		Email: data.Email,
		Name:  data.Name,
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    res,
	})
}

func (uh *UserController) Login(ctx *gin.Context) {
	var payload dto.LoginReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	data := entity.Login{
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := uh.uuc.LoginUser(ctx, &data); err != nil {
		ctx.Error(err)
		return
	}

	res := dto.LoginRes{
		ID:    data.ID,
		Email: data.Email,
		Name:  data.Name,
		Role:  data.Role,
		Token: data.Token,
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    res,
	})
}
