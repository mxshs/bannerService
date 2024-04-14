package handlers

import (
	"mxshs/bannerService/src/domain"
	"mxshs/bannerService/src/repositories"
	"mxshs/bannerService/src/services"

	"github.com/gin-gonic/gin"
)

const DEFAULT_PAGINATION_LIMIT = 100

type BannerServer struct {
    us *services.UserService
    bs *services.BannerService
}

func NewBannerServer(
    us repositories.UserStorage, signingToken []byte,
    bs repositories.BannerStorage, bc repositories.BannerCache,
    ) *BannerServer {

    server := &BannerServer{
        us: services.NewUserService(us, signingToken),
        bs: services.NewBannerService(bs, bc),
    }

    return server
}

func (bs *BannerServer) Run() {
    handler := gin.Default()
    
    handler.Static("/docs/swagger/", "docs/swaggerui")

    handler.POST("/signup/", bs.Signup)

    handler.POST("/login/", bs.Login)

    handler.POST("/banner/", bs.HandleAuth(domain.ADM), bs.CreateBanner)

    handler.GET("/user_banner", bs.HandleAuth(domain.USR), bs.GetBanner)

    handler.GET("/banner", bs.HandleAuth(domain.ADM), bs.GetBanners)

    handler.PATCH("/banner/:id", bs.HandleAuth(domain.ADM), bs.UpdateBanner)

    handler.DELETE("/banner/:id", bs.HandleAuth(domain.ADM), bs.DeleteBanner)

    handler.Run(":3000")
}
