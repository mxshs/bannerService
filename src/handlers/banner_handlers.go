package handlers

import (
	"mxshs/bannerService/src/models"
	"mxshs/bannerService/src/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (bs *BannerServer) CreateBanner(ctx *gin.Context) {
    var createBannerModel models.CreateBannerModel

    err := ctx.Bind(&createBannerModel)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }

    id, err := bs.bs.CreateBanner(createBannerModel.Tags, createBannerModel.Feature, createBannerModel.Content, createBannerModel.IsActive)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }

    ctx.JSON(http.StatusCreated, models.BannerCreatedModel{
        BannerId: id,
    })
}

func (bs *BannerServer) GetBanner(ctx *gin.Context) {
    tagId := ctx.Query("tag_id")
    if len(tagId) == 0 {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: "missing banner tag",
            },
        )
        return
    }

    tid, err := strconv.Atoi(tagId)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: "invalid banner tag",
            },
        )
        return
    }

    featureId := ctx.Query("feature_id")
    if len(featureId) == 0 {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: "missing banner feature",
            },
        )
        return
    }

    fid, err := strconv.Atoi(featureId)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: "invalid banner feature",
            },
        )
        return
    }

    useLastRevision := ctx.Query("use_last_revision")
    forceUpdate := len(useLastRevision) != 0 && useLastRevision == "true"

    res, err := bs.bs.GetBanner(tid, fid, forceUpdate)
    switch err {
    case nil:
        ctx.Data(http.StatusOK, "application/json", res)
    case domain.BannerNotFound:
        ctx.AbortWithStatus(http.StatusNotFound)
    case domain.InactiveBanner:
        ctx.AbortWithStatus(http.StatusForbidden)
    default:
        ctx.JSON(
            http.StatusInternalServerError,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
    }
}

func (bs *BannerServer) GetBanners(ctx *gin.Context) {
    var tagId, featureId *int
    var err error

    qtid := ctx.Query("tag_id")
    if len(qtid) != 0 {
        tid, err := strconv.Atoi(qtid)
        if err != nil {
            ctx.JSON(
                http.StatusBadRequest,
                models.ErrorModel{
                    Error: "invalid banner tag",
                },
            )
            return
        }
        tagId = &tid
    }

    qfid := ctx.Query("feature_id")
    if len(qfid) != 0 {
        fid, err := strconv.Atoi(qfid)
        if err != nil {
            ctx.JSON(
                http.StatusBadRequest,
                models.ErrorModel{
                    Error: "invalid banner feature",
                },
            )
            return
        }
        featureId = &fid
    }

    limit, offset := DEFAULT_PAGINATION_LIMIT, 0

    qlimit := ctx.Query("limit")
    if len(qlimit) != 0 {
        limit, err = strconv.Atoi(qlimit)
        if err != nil {
            ctx.JSON(
                http.StatusBadRequest,
                models.ErrorModel{
                    Error: "invalid limit",
                },
            )
            return
        }
    }

    qoffset := ctx.Query("offset")
    if len(qoffset) != 0 {
        offset, err = strconv.Atoi(qoffset)
        if err != nil {
            ctx.JSON(
                http.StatusBadRequest,
                models.ErrorModel{
                    Error: "invalid offset",
                },
            )
            return
        }
    }

    banners, err := bs.bs.GetBanners(tagId, featureId, limit, offset)
    if err != nil {
        ctx.JSON(
            http.StatusInternalServerError,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }

    ctx.JSON(http.StatusOK, banners)
}

func (bs *BannerServer) UpdateBanner (ctx *gin.Context) {
    bid, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: "invalid banner id",
            },
        )
        return
    }

    var updateBannerModel models.UpdateBannerModel

    if err = ctx.Bind(&updateBannerModel); err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
        return
    }

    err = bs.bs.UpdateBanner(
        bid, 
        updateBannerModel.Tags,
        updateBannerModel.Feature, 
        updateBannerModel.Content, 
        updateBannerModel.IsActive,
    )
    switch err {
    case nil:
        ctx.Status(http.StatusOK)
    case domain.BannerNotFound:
        ctx.AbortWithStatus(http.StatusNotFound)
    default:
        ctx.JSON(
            http.StatusInternalServerError,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
    }
}

func (bs *BannerServer) DeleteBanner (ctx *gin.Context) {
    bid, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            models.ErrorModel{
                Error: "invalid banner id",
            },
        )
        return
    }

    err = bs.bs.DeleteBanner(bid)
    switch err {
    case nil:
        ctx.Status(http.StatusNoContent)
    case domain.BannerNotFound:
        ctx.AbortWithStatus(http.StatusNotFound)
    default:
        ctx.JSON(
            http.StatusInternalServerError,
            models.ErrorModel{
                Error: err.Error(),
            },
        )
    }
}
