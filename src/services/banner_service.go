package services

import (
	"fmt"
	"mxshs/bannerService/src/domain"
	"mxshs/bannerService/src/repositories"
)

type BannerService struct {
    bannerRepo repositories.BannerStorage
    bannerCache repositories.BannerCache
}

func NewBannerService(bannerRepo repositories.BannerStorage, bannerCache repositories.BannerCache) *BannerService {
    return &BannerService{
        bannerRepo: bannerRepo,
        bannerCache: bannerCache,
    }
}

func (bs *BannerService) GetBanner(tagId int, featureId int, forceUpdate bool) (domain.BannerContent, error) {
    if !forceUpdate {
        content, err := bs.tryGetCachedBanner(tagId, featureId)
        if err == nil {
            return content, nil
        }
    }

    banner, err := bs.bannerRepo.GetBanner(tagId, featureId)
    if err != nil {
        return nil, err
    }

    if banner == nil {
        return nil, domain.BannerNotFound
    }

    if !banner.IsActive {
        return nil, domain.InactiveBanner
    }

    err = bs.bannerCache.SetBanner(tagId, featureId, banner.Content)
    if err != nil {
        fmt.Println(err.Error())
    }

    return banner.Content, nil
}

func (bs *BannerService) tryGetCachedBanner(tagId int, featureId int) (domain.BannerContent, error) {
    content, err := bs.bannerCache.GetBanner(tagId, featureId)
    switch err {
    case nil:
        return content, nil
    case domain.CacheMiss:
        fmt.Printf("[DEBUG] missed cache for tag %d and feature %d\n", tagId, featureId)
        fallthrough
    default:
        return nil, err
    }
}

func (bs *BannerService) GetBanners(tagId *int, featureId *int, limit int, offset int) ([]*domain.BannerDetails, error) {
    res, err := bs.bannerRepo.GetBanners(tagId, featureId, limit, offset)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (bs *BannerService) CreateBanner(tagIds []int, featureId int, content domain.BannerContent, isActive bool) (int, error) {
    id, err := bs.bannerRepo.CreateBanner(tagIds, featureId, content, isActive)

    return id, err
}

func (bs *BannerService) UpdateBanner(bannerId int, tagIds []int, featureId *int, content domain.BannerContent, isActive *bool) (error) {
    found, err := bs.bannerRepo.UpdateBanner(bannerId, tagIds, featureId, content, isActive)
    if err != nil {
        return err
    } else if !found {
        return domain.BannerNotFound
    }

    return nil
}

func (bs *BannerService) DeleteBanner(bannerId int) (error) {
    found, err := bs.bannerRepo.DeleteBanner(bannerId)
    if err != nil {
        return err
    } else if !found {
        return domain.BannerNotFound
    }

    return nil
}
