package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/whyslove/game-order-bot/internal/types"
)

func (gop *gameOrderPostgres) GetUser(userID int64) (types.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE tg_id = $1", usersTable)

	user := types.User{}
	err := gop.db.Get(&user, query, userID)
	if err != nil {
		return types.User{}, fmt.Errorf("error getting user from db with id: %d, err: %w", userID, err)
	}
	return user, nil
}

func (gop *gameOrderPostgres) SaveUserTx(tx *sqlx.Tx, user types.User) error {
	query := fmt.Sprintf("INSERT INTO %s (tg_id, name, is_admin, is_banned) VALUES ($1, $2, $3, $4)", usersTable)

	_, err := tx.Exec(query, user.TgID, user.Name, user.IsAdmin, user.IsBanned)
	if err != nil {
		return fmt.Errorf("error inserting into users user: %v, err: %w", user, err)
	}
	return nil
}

func (gop *gameOrderPostgres) GetTokenTx(tx *sqlx.Tx, token string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE token = $1 AND used = $2", tokensTable)

	tokenStruct := types.Token{}
	err := tx.Get(&tokenStruct, query, token, false)
	if err != nil {
		return fmt.Errorf("error getting token to check it token:%s, err:%w", token, err)
	}
	if tokenStruct.Token == "" {
		return fmt.Errorf("this token does not exits token: %s", token)
	}
	return nil
}

func (gop *gameOrderPostgres) UpdateTokenTx(tx *sqlx.Tx, token string) error {
	query := fmt.Sprintf("UPDATE %s SET used = true WHERE token = $1", tokensTable)

	_, err := tx.Exec(query, token)
	if err != nil {
		return fmt.Errorf("error while setting token user token: %s, err: %w", token, err)
	}
	return nil
}
