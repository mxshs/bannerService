package repositories

import "mxshs/bannerService/src/domain"

type BannerStorage interface {
    GetBanner(tagId int, featureId int) (*domain.Banner, error)
    GetBanners(tagId *int, featureId *int, limit int, offset int) ([]*domain.BannerDetails, error)
    CreateBanner(tagIds []int, featureId int, content domain.BannerContent, isActive bool) (int, error)
    UpdateBanner(bannerId int, tagIds []int, featureId *int,  content domain.BannerContent, isActive *bool) (bool, error)
    DeleteBanner(bannerId int) (bool, error)
}

type BannerCache interface {
    GetBanner(tagId, featureId int) (domain.BannerContent, error)
    InvalidateBanner(tagId, featureId int) error
    SetBanner(tagId int, featureId int, banner domain.BannerContent) error
}
