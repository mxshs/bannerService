package pg

import (
	"fmt"
	"mxshs/bannerService/src/domain"
	"time"

	"github.com/lib/pq"
)

func (db *PgDB) GetBanner(tagId int, featureId int) (*domain.Banner, error) {
    res, err := db.db.Query(
        `
        SELECT content, is_active
        FROM banners
        RIGHT JOIN (
            SELECT banner_features.bid FROM banner_features
            INNER JOIN banner_tags
            ON tid = $1 AND banner_features.bid = banner_tags.bid
            WHERE fid = $2
        ) as filtered
        ON banners.bid = filtered.bid
        LIMIT 1;
        `,
        tagId,
        featureId,
    )
    if err != nil {
        return nil, err
    }

    defer res.Close()

    if !res.Next() {
        return nil, nil
    }

    var banner domain.Banner

    err = res.Scan(&banner.Content, &banner.IsActive)
    if err != nil {
        return nil, err
    }

    return &banner, nil
}

func (db *PgDB) GetBanners(tagId *int, featureId *int, limit int, offset int) ([]*domain.BannerDetails, error) {
    res, err := db.db.Query(
        `
        SELECT banners.bid, banners.content, is_active, created_at, updated_at
        FROM banners
        RIGHT JOIN (
            SELECT DISTINCT banner_features.bid FROM banner_features
            INNER JOIN banner_tags
            ON tid = COALESCE($1, tid) AND banner_features.bid = banner_tags.bid
            WHERE fid = COALESCE($2, fid)
        ) as filtered
        ON banners.bid = filtered.bid
        LIMIT $3 OFFSET $4;
        `,
        tagId,
        featureId,
        limit,
        offset,
    )
    if err != nil {
        return nil, err
    }

    defer res.Close()

    banners := make([]*domain.BannerDetails, 0, limit)

    for idx := 0; idx < limit; idx++ {
        if !res.Next() {
            break
        }

        var banner domain.BannerDetails

        err = res.Scan(&banner.Id, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
        if err != nil {
            return nil, err
        }

        tags, err := db.db.Query(
            `
            SELECT tid FROM banner_tags
            WHERE bid = $1;
            `,
            banner.Id,
        )
        if err != nil {
            return nil, err
        }

        for tags.Next() {
            var tag int

            if err = tags.Scan(&tag); err != nil {
                return nil, err
            }

            banner.Tags = append(banner.Tags, tag)
        }

        tags.Close()

        feature, err := db.db.Query(
            `
            SELECT fid FROM banner_features
            WHERE bid = $1;
            `,
            banner.Id,
        )
        if err != nil {
            return nil, err
        }

        if feature.Next() {
            var featureId int

            if err = feature.Scan(&featureId); err != nil {
                return nil, err
            }

            banner.Feature = &featureId
        }

        feature.Close()

        banners = append(banners, &banner)
    }

    return banners, nil
}


func (db *PgDB) CreateBanner(tagIds []int, featureId int, content domain.BannerContent, isActive bool) (int, error) {
    transaction, err := db.db.Begin()
    if err != nil {
        return 0, err
    }

    defer func() {
        if err := transaction.Rollback(); err != nil {
            fmt.Printf("[ERROR] failed transaction rollback: %s\n", err.Error())
        }
    }()

    res, err := transaction.Query(
        `
        INSERT INTO banners(content, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING banners.bid;
        `,
        content,
        isActive,
        time.Now(),
        time.Now(),
    )
    if err != nil {
        return 0, err
    }

    defer res.Close()

    if !res.Next() {
        return 0, fmt.Errorf("unexpected empty result set after successful insertion")
    }

    var bannerId int

    err = res.Scan(&bannerId)
    if err != nil {
        return 0, err
    }

    for _, id := range tagIds {
        _, err = transaction.Exec(
            `
            INSERT INTO banner_tags
            VALUES ($1, $2);
            `,
            bannerId,
            id,
        )
        if err != nil {
            return 0, err
        }
    }

    _, err = transaction.Exec(
        `
        INSERT INTO banner_features
        VALUES ($1, $2);
        `,
        bannerId,
        featureId,
    )
    if err != nil {
        return 0, err
    }

    return bannerId, transaction.Commit()
}

func (db *PgDB) UpdateBanner(bannerId int, tagIds []int, featureId *int,  content domain.BannerContent, isActive *bool) (bool, error) {
    transaction, err := db.db.Begin()
    if err != nil {
        return false, err
    }

    defer func() {
        if err := transaction.Rollback(); err != nil {
            fmt.Printf("[ERROR] failed transaction rollback: %s\n", err.Error())
        }
    }()

    if content != nil {
        res, err := transaction.Exec(
            `
            UPDATE banners
            SET content = $1,
            updated_at = $2
            WHERE bid = $3;
            `,
            content,
            time.Now(),
            bannerId,
        )
        if err != nil {
            return false, err
        }

        cnt, _ := res.RowsAffected()
        if cnt == 0 {
            return false, nil
        }
    }

    if isActive != nil {
        _, err := transaction.Exec(
            `
            UPDATE banners
            SET is_active = $1,
            updated_at = $2
            WHERE bid = $3;
            `,
            isActive,
            time.Now(),
            bannerId,
        )
        if err != nil {
            return false, err
        }
    }

    if tagIds != nil {
        _, err = transaction.Exec(
            `
            DELETE FROM banner_tags
            WHERE bid = $1;
            `,
            bannerId,
        )
        if err != nil {
            return false, err
        }

        row, err := transaction.Prepare(pq.CopyIn("banner_tags", "bid", "tid"))
        if err != nil {
            return false, err
        }

        for _, id := range tagIds {
            if _, err = row.Exec(bannerId, id); err != nil {
                return false, err
            }
        }

        if _, err = row.Exec(); err != nil { 
            return false, err
        }

        if err = row.Close(); err != nil {
            return false, err
        }
    }

    if featureId != nil {
        _, err = transaction.Exec(
            `
            DELETE FROM banner_features
            WHERE bid = $1;
            `,
            bannerId,
        )
        if err != nil {
            return false, err
        }

        _, err = transaction.Exec(
            `
            INSERT INTO banner_features
            VALUES ($1, $2);
            `,
            bannerId,
            featureId,
        )
        if err != nil {
            return false, err
        }
    }

    return true, transaction.Commit()
}

func (db *PgDB) DeleteBanner(bannerId int) (bool, error) {
    res, err := db.db.Exec(
        `
        DELETE FROM banners
        WHERE bid = $1;
        `,
        bannerId,
    )
    if err != nil {
        return false, err
    }

    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return false, nil
    }

    return true, nil
}
