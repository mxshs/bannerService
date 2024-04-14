package redis

import (
	"context"
	"fmt"
	"mxshs/bannerService/src/domain"

	"github.com/redis/go-redis/v9"
)

func (rc *RedisCache) GetBanner(tagId int, featureId int) (domain.BannerContent, error) {
    content, err := rc.rdb.Get(context.Background(), fmt.Sprintf("banner:%d:%d", tagId, featureId)).Result()
    switch err {
    case nil:
        return domain.BannerContent(content), nil
    case redis.Nil:
        return nil, domain.CacheMiss
    default:
        return nil, err
    }
}

func (rc *RedisCache) SetBanner(tagId, featureId int, content domain.BannerContent) error {
    _, err := rc.rdb.Set(context.Background(), fmt.Sprintf("banner:%d:%d", tagId, featureId), []byte(content), rc.ttl).Result()

    return err
}

func (rc *RedisCache) InvalidateBanner(tagId, featureId int) error {
    _, err := rc.rdb.Del(context.Background(), fmt.Sprintf("banner:%d:%d", tagId, featureId)).Result()

    return err
}

