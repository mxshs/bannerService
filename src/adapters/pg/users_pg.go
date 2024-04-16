package pg

import (
	"fmt"
	"mxshs/bannerService/src/domain"
	"time"
)

func (db *PgDB) GetUser(uid int) (*domain.User, error) {
    res, err := db.db.Query(
        `
        SELECT * FROM users
        WHERE uid = $1;
        `,
        uid,
    )
    if err != nil {
        return nil, err
    }

    defer res.Close()

    if !res.Next() {
        return nil, nil
    }

    var user domain.User

    err = res.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.TokenExpiresAt)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (db *PgDB) LoginUser(username, password string) (*domain.User, error) {
    res, err := db.db.Query(
        `
        SELECT * FROM users
        WHERE username = $1;
        `,
        username,
    )
    if err != nil {
        return nil, err
    }

    defer res.Close()

    if !res.Next() {
        return nil, nil
    }

    var user domain.User

    err = res.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.TokenExpiresAt)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (db *PgDB) CreateUser(username, password string, role domain.Role) (*domain.User, error) {
    res, err := db.db.Query(
        `
        INSERT INTO users(username, password, role, eat)
        VALUES($1, $2, $3, $4)
        RETURNING *;
        `,
        username,
        password,
        role,
        time.Now().Unix(),
    )
    if err != nil {
        return nil, err
    }

    defer res.Close()

    if !res.Next() {
        return nil, fmt.Errorf("unexpected empty result set after successful insertion")
    }

    var user domain.User

    err = res.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.TokenExpiresAt)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (db *PgDB) UpdateUser(id int, username, password *string, role *domain.Role) (*domain.User, error) {
    res, err := db.db.Query(
        `
        UPDATE users
        SET username = COALESCE($1, username),
        SET password = COALESCE($1, password),
        SET role = COALESCE($1, role),
        RETURNING *;
        `,
        username,
        password,
        role,
    )
    if err != nil {
        return nil, err
    }

    defer res.Close()

    if !res.Next() {
        return nil, fmt.Errorf("unexpected empty result set after successful update")
    }

    var user domain.User

    err = res.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.TokenExpiresAt)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (db *PgDB) UpdateTokens(id int, expirationDate time.Time) (error) {
    _, err := db.db.Exec(
        `
        UPDATE users
        SET eat = $1
        WHERE uid = $2;
        `,
        expirationDate.Unix(),
        id,
    )

    return err
}
