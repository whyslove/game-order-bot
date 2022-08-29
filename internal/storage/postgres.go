package storage

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/whyslove/game-order-bot/internal/types"
)

type DatabaseI interface {
	GetTeamsForDay(day time.Time) ([]types.Team, error)
	CreateTeam(name string, ownderID int64, ownerTag string, dateCreated time.Time, deleted bool) (int64, error)
	StartTransaction() (*sqlx.Tx, error)
	CommitTransaction(tx *sqlx.Tx) error
	RollbackTransaction(tx *sqlx.Tx) error
	GetTeam(tx *sqlx.Tx, teamID int64, date time.Time) (types.Team, error)
	DeleteTeam(teamID int64, date time.Time) error
	SetMatchesQueue(date time.Time, matches []types.MatchQueue) error
	GetMatchesQueue(date time.Time) ([]types.MatchQueue, error)
	GetMyTeams(ownderID int64, day time.Time) ([]types.Team, error)
	DeleteAllMatches(tx *sqlx.Tx, day time.Time) error
	DeleteAllTeams(tx *sqlx.Tx, day time.Time) error
}

type gameOrderPostgres struct {
	db *sqlx.DB
}

const (
	adminsTable  = "admins"
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
