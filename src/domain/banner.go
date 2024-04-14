package domain

import "time"

type BannerContent []byte

func (bc *BannerContent) UnmarshalJSON(object []byte) error {
    *bc = object
    return nil
}

func (bc *BannerContent) MarshalJSON() ([]byte, error) {
    return []byte(*bc), nil
}

type Banner struct {
    Id int `json:"banner_id"`
    Content BannerContent `json:"content"`
    IsActive bool `json:"is_active"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

type BannerDetails struct {
    Id int `json:"banner_id"`
    Tags []int `json:"tag_ids"`
    Feature *int `json:"feature_id"`
    Content BannerContent `json:"content"`
    IsActive *bool `json:"is_active"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}
