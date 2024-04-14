package domain

const CacheMiss = BaseError("no cached value")
const InactiveBanner = BaseError("cannot access requested banner")
const BannerNotFound = BaseError("banner not found")
const UserNotFound = BaseError("user not found")

type BaseError string

func (be BaseError) Error() string {
    return string(be)
}
