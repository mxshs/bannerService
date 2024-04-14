package handlers

import (
	"mxshs/bannerService/src/domain"
	"mxshs/bannerService/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bs *BannerServer) Signup (ctx *gin.Context) {
    var userModel models.CreateUserModel

    err := ctx.Bind(&userModel)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }

    tokens, err := bs.us.Signup(userModel.Username, userModel.Password, userModel.Role)
    switch err {
    case nil:
        ctx.JSON(http.StatusOK, tokens)
    default:
        ctx.JSON(
            http.StatusInternalServerError,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }
}

func (bs *BannerServer) Login (ctx *gin.Context) {
    var userModel models.LoginUserModel

    err := ctx.Bind(&userModel)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }

    tokens, err := bs.us.LoginUser(userModel.Username, userModel.Password)
    switch err {
    case nil: 
        ctx.JSON(http.StatusOK, tokens)
    case domain.UserNotFound:
        ctx.AbortWithStatus(http.StatusNotFound)
    default:
        ctx.JSON(
            http.StatusInternalServerError, 
            models.ErrorModel{
                Error: err.Error(),
            },
        )
    }
}
