package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/storage"
	"github.com/whyslove/game-order-bot/internal/types"
)

var NotAdminErr = errors.New("Чтобы выполнить это действие вы должны быть админом")
var BannedErr = errors.New("Ты забанен :)))))))))))))))д")
var CantDoThis = errors.New("Вы не можете сделать это")

type authorization struct {
	db  storage.DatabaseI
	svc ServiceI
}

func NewAuthotization(db storage.DatabaseI) *authorization {
	return &authorization{db: db}
}

func (auth *authorization) CreateTeam(userID int64, name string, ownerID int64, ownerTag string, members string) error {
	verdict, err := auth.auth(userID)
	if err != nil {
		return err
	}
	if !verdict {
		return BannedErr
	}

	return auth.svc.CreateTeam(userID, name, ownerID, ownerTag, members)
}
func (auth *authorization) GetTodayTeams(userID int64) ([]types.Team, error) {
	verdict, err := auth.authAdmin(userID)
	if err != nil {
		return nil, err
	}
	if !verdict {
		return nil, NotAdminErr
	}

	return auth.svc.GetTodayTeams(userID)
}

func (auth *authorization) GetAllTodayMatches() []types.MatchQueue {
	return auth.svc.GetAllTodayMatches()
}

func (auth *authorization) SetMatchPlayed(userID int64, leftStays bool) error {
	verdict, err := auth.auth(userID)
	if err != nil {
		return err
	}
	if !verdict {
		return CantDoThis
	}

	return auth.svc.SetMatchPlayed(userID, leftStays)
}

func (auth *authorization) DeleteTeam(userID int64, teamID int64) error {
	verdict, err := auth.checkTeamOwnership(userID, teamID)
	if err != nil {
		return err
	}
	if !verdict {
		return CantDoThis
	}
	return auth.svc.DeleteTeam(userID, teamID)
}

func (auth *authorization) GetMyTeams(userID int64) ([]types.Team, error) {
	verdict, err := auth.auth(userID)
	if err != nil {
		return nil, err
	}
	if !verdict {
		return nil, BannedErr
	}

	return auth.svc.GetMyTeams(userID)
}
func (auth *authorization) DeleteAllInformationToday(userID int64) error {
	verdict, err := auth.authAdmin(userID)
	if err != nil {
		return err
	}
	if !verdict {
		return NotAdminErr
	}

	return auth.svc.DeleteAllInformationToday(userID)
}

func (auth *authorization) SaveUser(userID int64, token string, user types.User) error {
	verdict, err := auth.auth(userID)
	if err != nil {
		return err
	}
	if !verdict {
		return BannedErr
	}

	return auth.svc.SaveUser(userID, token, user)
}

func (auth *authorization) UpdateTeamMembers(userID, teamID int64, members string) error {
	verdict, err := auth.checkTeamOwnership(userID, teamID)
	if err != nil {
		return fmt.Errorf("error in authorization in updateTeamMembers, err: %w", err)
	}
	if !verdict {
		return CantDoThis
	}
	return auth.svc.UpdateTeamMembers(userID, teamID, members)
}

func (auth *authorization) auth(userID int64) (bool, error) {
	user, err := auth.db.GetUser(userID)
	if err != nil {
		log.Error().Msgf("error in database while querying userID:%d, err: %s", userID, err.Error())
		return false, err
	}
	return user.IsBanned, nil
}

func (auth *authorization) authAdmin(userID int64) (bool, error) {
	user, err := auth.db.GetUser(userID)
	if err != nil {
		log.Error().Msgf("error in database while querying userID:%d, err: %s", userID, err.Error())
		return false, err
	}
	return user.IsAdmin, nil
}

func (auth *authorization) checkTeamOwnership(userID int64, teamID int64) (bool, error) {
	verdict, err := auth.authAdmin(userID)
	if err != nil || !verdict {
		teams, err := auth.db.GetMyTeams(userID, time.Now())
		if err != nil {
			return false, fmt.Errorf("error getting team for userID %d, err: %s", userID, err.Error())
		}

		for i := range teams {
			if teams[i].Id == teamID && teams[i].OwnerID == userID {
				return true, nil
			}
		}
		return false, nil
	}
	return true, nil
}

func (auth *authorization) CheckIsAdmin(userID int64) (bool, error) {
	verdict, err := auth.auth(userID)
	if err != nil {
		return false, err
	}
	if !verdict {
		return false, BannedErr
	}

	return auth.svc.CheckIsAdmin(userID)
}
