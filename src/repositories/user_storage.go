package repositories

import (
	"mxshs/bannerService/src/domain"
	"time"
)

type UserStorage interface {
    GetUser(id int) (*domain.User, error)
    LoginUser(username, password string) (*domain.User, error)
    CreateUser(username, password string, role domain.Role) (*domain.User, error)
    UpdateUser(id int, username, password *string, role *domain.Role) (*domain.User, error)
    UpdateTokens(id int, expirationDate time.Time) (error)
}
