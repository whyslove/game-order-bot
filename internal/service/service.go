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
	CreateTeam(userID int64, name string, ownerID int64, ownerTag string, members string) error
	GetTodayTeams(userID int64) ([]types.Team, error)
	GetAllTodayMatches() []types.MatchQueue
	SetMatchPlayed(userID int64, leftStays bool) error
	DeleteTeam(userID int64, teamID int64) error
	UpdateTeamMembers(userID, teamID int64, members string) error
	GetMyTeams(userID int64) ([]types.Team, error)
	DeleteAllInformationToday(userID int64) error
	SaveUser(userID int64, token string, user types.User) error
	CheckIsAdmin(userID int64) (bool, error)
	RefreshMatches()
}
