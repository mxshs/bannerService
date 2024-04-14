package handlers

import (
	"mxshs/bannerService/src/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bs *BannerServer) HandleAuth(permission domain.Role) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        res := ctx.GetHeader("token")
        if len(res) == 0 {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        err := bs.us.ValidateToken(res, permission)
        if err != nil {
            ctx.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        ctx.Next()
    }
}
