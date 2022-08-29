package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/whyslove/game-order-bot/internal/types"
)

const dateLayout = "2006-01-02"

func (gop *gameOrderPostgres) GetTeamsForDay(day time.Time) ([]types.Team, error) {
	query := fmt.Sprintf("SELECT * from %s WHERE date_created = $1", teamsTable)
	log.Debug().Msgf("query to database query: %s", query)
	teams := []types.Team{}

	err := gop.db.Select(&teams, query, day.Format(dateLayout))
	// rows, err := gop.db.Queryx(query, day.Format(dateLayout))
	if err != nil {
		return nil, fmt.Errorf("error while gettig all teams for today, err:%w", err)
	}
	return teams, nil
}
func (gop *gameOrderPostgres) CreateTeam(name string, ownderID int64, ownerTag string, dateCreated time.Time, deleted bool) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (name, owner_id, owner_tag, date_created, deleted) VALUES
		($1, $2, $3, $4, $5) RETURNING id;`, teamsTable)
	log.Debug().Msgf("database query: %s", query)

	var lastInsertedID int64
	err := gop.db.QueryRow(query, name, ownderID, ownerTag, dateCreated.Format(dateLayout), deleted).Scan(&lastInsertedID)
	if err != nil {
		return 0, fmt.Errorf("error inserting new team, err: %w", err)
	}
	return lastInsertedID, err
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

func (gop *gameOrderPostgres) GetTeam(tx *sqlx.Tx, teamID int64, date time.Time) (types.Team, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE team_id = $1 AND date_created = $2", teamsTable)
	log.Debug().Msgf("query to get team %s", query)

	team := types.Team{}
	err := tx.Get(&team, query)
	if err != nil {
		return types.Team{}, fmt.Errorf("error while getting team with id: %d, err: %w", teamID, err)
	}
	return team, nil
}

func (gop *gameOrderPostgres) DeleteTeam(teamID int64, date time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET deleted = $1 WHERE id = $2 AND date_created = $3", teamsTable)
	log.Debug().Msgf("query to delete team %s", query)

	_, err := gop.db.Exec(query, true, teamID, date.Format(dateLayout))
	if err != nil {
		return fmt.Errorf("error deleteing team with id teamID %d, err: %w", teamID, err)
	}
	return nil
}

func (gop *gameOrderPostgres) SetMatchesQueue(date time.Time, matches []types.MatchQueue) error {
	log.Info().Msgf("setting matches queue into database %v", matches)

	bts, err := json.Marshal(matches)
	if err != nil {
		return fmt.Errorf("error while marshaling struct, err: %w", err)
	}
	query := fmt.Sprintf("INSERT INTO %s (date_created, matches_queue) VALUES ($1, $2) ON CONFLICT ON CONSTRAINT matches_pkey DO UPDATE SET matches_queue = $2 WHERE %s.date_created = $1", matchesTable, matchesTable)
	_, err = gop.db.Exec(query, date.Format(dateLayout), bts)
	if err != nil {
		return fmt.Errorf("error upserting into matches err: %w", err)
	}
	return nil
}

func (gop *gameOrderPostgres) GetMatchesQueue(date time.Time) ([]types.MatchQueue, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE date_created = $1", matchesTable)

	mq := types.DatabaseMatchesQueue{}
	err := gop.db.Get(&mq, query, date.Format(dateLayout))
	if err != nil {
		return nil, fmt.Errorf("error in getting matches queue, err: %w", err)
	}

	matchesQueue := []types.MatchQueue{}
	err = json.Unmarshal(mq.BinaryMatchesQueue, &matchesQueue)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes into queue, err: %w", err)
	}
	log.Debug().Msgf("matchesQueue: %v", matchesQueue)
	return matchesQueue, nil
}

func (gop *gameOrderPostgres) GetMyTeams(ownderID int64, day time.Time) ([]types.Team, error) {
	query := fmt.Sprintf("SELECT * from %s WHERE date_created = $1 AND owner_id = $2 AND deleted = $3", teamsTable)
	log.Debug().Msgf("query to database query: %s", query)
	teams := []types.Team{}

	err := gop.db.Select(&teams, query, day.Format(dateLayout), ownderID, false)
	if err != nil {
		return nil, fmt.Errorf("error while gettig all teams for today for user: %d, err: %w", ownderID, err)
	}
	return teams, nil
}

func (gop *gameOrderPostgres) DeleteAllTeams(tx *sqlx.Tx, day time.Time) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE date_created = $1", teamsTable)
	_, err := tx.Exec(query, day.Format(dateLayout))
	return err
}

func (gop *gameOrderPostgres) DeleteAllMatches(tx *sqlx.Tx, day time.Time) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE date_created = $1", matchesTable)
	_, err := tx.Exec(query, day.Format(dateLayout))
	return err
}
