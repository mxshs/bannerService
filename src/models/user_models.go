package models

import "mxshs/bannerService/src/domain"

type CreateUserModel struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Role domain.Role `json:"role"`
}

type LoginUserModel struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
