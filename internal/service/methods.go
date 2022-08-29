package service

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/types"
)

func GetSchedule() {

}

func (svc *service) CreateTeam(name string, ownerID int64, ownerTag string) error {
	nowDate := time.Now()
	deleted := false
	teamID, err := svc.db.CreateTeam(name, ownerID, ownerTag, nowDate, deleted)
	if err != nil {
		return fmt.Errorf("error while creatin team in db, err: %w", err)
	}
	log.Debug().Msgf("teamID %d", teamID)

	// Next in matches thing
	// Team 1 plays one match and runs off
	if len(svc.Matches) == 0 {
		newMatch := types.MatchQueue{
			Team1:       name,
			Team1ID:     teamID,
			Team2:       "",
			Current:     false,
			Played:      false,
			DateCreated: time.Now(),
		}
		svc.Matches = append(svc.Matches, newMatch)
		svc.CurrentMatchIndex = 0
	} else {
		if len(svc.Matches) == 1 {
			svc.Matches[0].Current = true
		}
		svc.Matches[len(svc.Matches)-1].Team2 = name
		svc.Matches[len(svc.Matches)-1].Team2ID = teamID

		newMatch := types.MatchQueue{
			Team1:       name,
			Team1ID:     teamID,
			Team2:       "",
			Current:     false,
			Played:      false,
			DateCreated: time.Now(),
		}
		// svc.Matches[len(svc.Matches)-1].NextGame = newMatch
		svc.Matches = append(svc.Matches, newMatch)
	}
	err = svc.db.SetMatchesQueue(time.Now(), svc.Matches)
	return err
}

func (svc *service) GetAllTodayMatches() []types.MatchQueue {
	return svc.Matches
}

// leftStays means team that was on the left will play one more game
func (svc *service) SetMatchPlayed(leftStays bool) {
	if len(svc.Matches) == 0 {
		return
	}
	playedMatch := svc.Matches[svc.CurrentMatchIndex]

	svc.Matches[svc.CurrentMatchIndex].Current = false
	svc.Matches[svc.CurrentMatchIndex].Played = true
	svc.Matches[svc.CurrentMatchIndex].Score = "25-25"

	if leftStays {
		svc.matchPlayedLeftStays(playedMatch)
	} else {
		svc.matchPlayedRightStays(playedMatch)
	}

	// Update index match playing
	svc.CurrentMatchIndex++
	svc.Matches[svc.CurrentMatchIndex].Current = true

	err := svc.db.SetMatchesQueue(time.Now(), svc.Matches)
	if err != nil {
		log.Error().Msgf("error setiing matches queue after match played %s", err.Error())
	}

}

func (svc *service) GetTodayTeams() ([]types.Team, error) {
	nowDate := time.Now()
	teams, err := svc.db.GetTeamsForDay(nowDate)
	if err != nil {
		return nil, fmt.Errorf("error while getting all teams fot date: %s, err: %w", nowDate.String(), err)
	}
	return teams, nil
}

func (svc *service) DeleteTeam(teamID int64) error {
	// GetTeam
	// TODO добавить проверку на принадлежность команды этому игроку
	if len(svc.Matches) <= 2 {
		return fmt.Errorf("Не хочу удалять когда длина 2 или меньше")
	}

	var teamWasFound bool
	//Delete from queue
	for i := svc.CurrentMatchIndex; i < len(svc.Matches)-1; i++ {
		//TODO Удаление команды, которая прямо сейчас вышла на поле

		if svc.Matches[i].Team2ID == teamID {
			teamWasFound = true
			for j := i; j < len(svc.Matches)-1; j++ {
				svc.Matches[j].Team2ID = svc.Matches[j+1].Team2ID
				svc.Matches[j].Team2 = svc.Matches[j+1].Team2

				svc.Matches[j+1].Team1ID = svc.Matches[j+1].Team2ID
				svc.Matches[j+1].Team1 = svc.Matches[j+1].Team2
			}
		}
	}
	if teamWasFound == false {
		return fmt.Errorf("this team was not found in queue")
	}
	svc.Matches = svc.Matches[:len(svc.Matches)-1]

	err := svc.db.DeleteTeam(teamID, time.Now())
	if err != nil {
		return fmt.Errorf("error deleting team err:%w", err)
	}
	err = svc.db.SetMatchesQueue(time.Now(), svc.Matches)
	if err != nil {
		return fmt.Errorf("error setting queue while deleting team err:%w", err)
	}

	return nil
}

func (svc *service) GetMyTeams(ownderID int64) ([]types.Team, error) {
	return svc.db.GetMyTeams(ownderID, time.Now())
}

func (svc *service) DeleteAllInformationToday() (err error) {
	tx, err := svc.db.StartTransaction()
	if err != nil {
		return fmt.Errorf("error creating transaction err: %w", err)
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error().Msgf("error in rollbalck err: %s", err.Error())
			}
		} else {
			err = tx.Commit()
		}
	}()

	currentTime := time.Now()
	err = svc.db.DeleteAllMatches(tx, currentTime)
	if err != nil {
		err = fmt.Errorf("error in deleting all matches, err: %w", err)
		return
	}
	err = svc.db.DeleteAllTeams(tx, currentTime)
	if err != nil {
		err = fmt.Errorf("error in deleting all teams, err: %w", err)
		return
	}
	svc.Matches = nil
	svc.Matches = make([]types.MatchQueue, 0)

	return nil

}

func (svc *service) matchPlayedLeftStays(playedMatch types.MatchQueue) {
	svc.Matches[len(svc.Matches)-1].Team2 = playedMatch.Team2
	svc.Matches[len(svc.Matches)-1].Team2ID = playedMatch.Team2ID

	newMatch := types.MatchQueue{
		Team1:       playedMatch.Team2,
		Team1ID:     playedMatch.Team2ID,
		Team2:       "",
		Current:     false,
		Played:      false,
		DateCreated: time.Now(),
	}
	svc.Matches = append(svc.Matches, newMatch)

	//Обновить в некст паре
	//TODO make better, currentMatchIndex Support Concurrency
	svc.Matches[svc.CurrentMatchIndex+1].Team1 = playedMatch.Team1
	svc.Matches[svc.CurrentMatchIndex+1].Team1ID = playedMatch.Team1ID

}

func (svc *service) matchPlayedRightStays(playedMatch types.MatchQueue) {
	svc.Matches[len(svc.Matches)-1].Team2 = playedMatch.Team1
	svc.Matches[len(svc.Matches)-1].Team2ID = playedMatch.Team1ID

	newMatch := types.MatchQueue{
		Team1:       playedMatch.Team1,
		Team1ID:     playedMatch.Team1ID,
		Team2:       "",
		Current:     false,
		Played:      false,
		DateCreated: time.Now(),
	}
	svc.Matches = append(svc.Matches, newMatch)
}
