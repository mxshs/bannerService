package models

import "mxshs/bannerService/src/domain"

type CreateBannerModel struct {
    Tags []int `json:"tag_ids"`
    Feature int `json:"feature_id"`
    Content domain.BannerContent `json:"content"`
    IsActive bool `json:"is_active"`
}

type BannerCreatedModel struct {
    BannerId int `json:"banner_id"`
}

type UpdateBannerModel struct {
    Id int `json:"banner_id"`
    Tags []int `json:"tag_ids"`
    Feature *int `json:"feature_id"`
    Content domain.BannerContent `json:"content"`
    IsActive *bool `json:"is_active"`
}
