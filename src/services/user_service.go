package services

import (
	"fmt"
	"mxshs/bannerService/src/domain"
	"mxshs/bannerService/src/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
    UserRepo repositories.UserStorage
    signingToken []byte
}

func NewUserService(ur repositories.UserStorage, token []byte) *UserService {
    return &UserService{
        UserRepo: ur,
        signingToken: token,
    }
}

func (us *UserService) Signup(username, password string, role domain.Role) (*domain.TokenPair, error) {
    user, err := us.UserRepo.CreateUser(username, password, role)
    if err != nil {
        return nil, err
    }

    tokens, err := us.createDefaultTokenPair(user.ID, role)
    if err != nil {
        return nil, err
    }

    return tokens, nil
}

func (us *UserService) LoginUser(username, password string) (*domain.TokenPair, error) {
    user, err := us.UserRepo.LoginUser(username, password)
    if err != nil {
        return nil, err
    }

    if user == nil {
        return nil, domain.UserNotFound
    }

    tokens, err := us.createDefaultTokenPair(user.ID, user.Role)
    if err != nil {
        return nil, err
    }

    return tokens, nil
}

func (us *UserService) ValidateToken(token string, role domain.Role) error {
    claims, err := us.verifyToken(token)
    if err != nil {
        return err
    }

    if claims.Role < role {
        return fmt.Errorf("this token has insufficient permissions")
    }

    return nil
}

func (us *UserService) RefreshToken(token string) (*domain.TokenPair, error) {
    claims, err := us.verifyToken(token)
    if err != nil {
        return nil, err
    }

    user, err := us.UserRepo.GetUser(claims.UID)
    if err != nil {
        return nil, err
    }

    if user == nil {
        return nil, domain.UserNotFound
    }

    if user.TokenExpiresAt > claims.ExpiresAt.Unix() {
        return nil, fmt.Errorf("provided token was already refreshed")
    }

    tokens, err := us.createDefaultTokenPair(claims.UID, claims.Role)
    if err != nil {
        return nil, err
    }

    return tokens, nil
}

func (us *UserService) createDefaultTokenPair(uid int, role domain.Role) (*domain.TokenPair, error) {
    aeat := time.Now().Add(time.Hour)
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.TokenClaims{
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(aeat),
        },
    })
    accessString, err := accessToken.SignedString(us.signingToken)
    if err != nil {
        return nil, err
    }

    reat := time.Now().Add(24 * time.Hour)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.TokenClaims{
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(reat),
        },
    })
    refreshString , err := refreshToken.SignedString(us.signingToken)
    if err != nil {
        return nil, fmt.Errorf("failed to sign token with provided key")
    }

    err = us.UserRepo.UpdateTokens(uid, reat)
    if err != nil {
        return nil, fmt.Errorf("failed to update token expiration date")
    }

    return &domain.TokenPair{
        AccessToken: accessString,
        RefreshToken: refreshString,
    }, nil
}

func (us *UserService) verifyToken(token string) (*domain.TokenClaims, error) {
    result, err := jwt.ParseWithClaims(token, &domain.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
        return us.signingToken, nil
    })
    if err != nil || !result.Valid {
        return nil, fmt.Errorf("received invalid jwt key")
    }

    claims, ok := result.Claims.(*domain.TokenClaims)
    if !ok {
        return nil, fmt.Errorf("received jwt key with unexpected claims")
    } else if claims.ExpiresAt.Before(time.Now()) {
        return nil, fmt.Errorf("received expired jwt key")
    }

    return claims, nil
}
