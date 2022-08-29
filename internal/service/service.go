package service

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/storage"
	"github.com/whyslove/game-order-bot/internal/types"
)

type service struct {
	db                storage.DatabaseI
	Matches           []types.MatchQueue
	CurrentMatchIndex int
}

func NewService(db storage.DatabaseI) *service {
	matches, err := db.GetMatchesQueue(time.Now())
	if err != nil {
		log.Error().Msgf("error while preloading matches from db, err: %s", err.Error())
		return &service{db: db, Matches: make([]types.MatchQueue, 0)}
	}
	return &service{db: db, Matches: matches}

}

type ServiceI interface {
	CreateTeam(name string, ownerID int64, ownerTag string) error
	GetTodayTeams() ([]types.Team, error)
	GetAllTodayMatches() []types.MatchQueue
	SetMatchPlayed()
	DeleteTeam(teamID int64) error
	GetMyTeams(ownderID int64) ([]types.Team, error)
	DeleteAllInformationToday() error
}
