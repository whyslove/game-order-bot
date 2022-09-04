package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/whyslove/game-order-bot/internal/types"
)

type DatabaseI interface {
	GetTeamsForDay(day time.Time) ([]types.Team, error)
	CreateTeam(name string, ownderID int64, ownerTag string, members string, dateCreated time.Time, deleted bool) (int64, error)
	StartTransaction() (*sqlx.Tx, error)
	CommitTransaction(tx *sqlx.Tx) error
	RollbackTransaction(tx *sqlx.Tx) error
	GetTeamTx(tx *sqlx.Tx, teamID int64, date time.Time) (types.Team, error)
	GetTeam(teamID int64, date time.Time) (types.Team, error)
	DeleteTeam(teamID int64, date time.Time) error
	UpdateTeamMembers(teamID int64, date time.Time, members string) error
	SetMatchesQueue(date time.Time, matches []types.MatchQueue) error
	GetMatchesQueue(date time.Time) ([]types.MatchQueue, error)
	GetMyTeams(ownderID int64, day time.Time) ([]types.Team, error)
	DeleteAllMatches(tx *sqlx.Tx, day time.Time) error
	DeleteAllTeams(tx *sqlx.Tx, day time.Time) error

	GetUser(userID int64) (types.User, error)
	SaveUserTx(tx *sqlx.Tx, user types.User) error
	GetTokenTx(tx *sqlx.Tx, token string) error
	UpdateTokenTx(tx *sqlx.Tx, token string) error
}

type gameOrderPostgres struct {
	db *sqlx.DB
}

const (
	usersTable   = "users"
	teamsTable   = "teams"
	matchesTable = "matches"
	tokensTable  = "tokens"
)

func NewPostgresDb(host, port, user, dbname, password, sslmode string) (*gameOrderPostgres, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	gop := gameOrderPostgres{
		db: db,
	}
	return &gop, nil

}

func (gop *gameOrderPostgres) StartTransaction() (*sqlx.Tx, error) {
	return gop.db.BeginTxx(context.Background(), nil)
}

func (gop *gameOrderPostgres) CommitTransaction(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (gop *gameOrderPostgres) RollbackTransaction(tx *sqlx.Tx) error {
	return tx.Rollback()
}
