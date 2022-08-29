package types

import (
	"time"
)

type Team struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	OwnerID     int64     `db:"owner_id"`
	OnwerTag    string    `db:"owner_tag"`
	DateCreated time.Time `db:"date_created"`
	Deleted     bool      `db:"deleted"`
}
type DatabaseMatchesQueue struct {
	DateCreated        time.Time `db:"date_created"`
	BinaryMatchesQueue []byte    `db:"matches_queue"`
}

type MatchQueue struct {
	Id          int64     `db:"id" json:"id"`
	Team1       string    `db:"team_1" json:"team_1"`
	Team2       string    `db:"team_2" json:"team_2"`
	Team1ID     int64     `db:"team_1_id" json:"team_1_id"`
	Team2ID     int64     `db:"team_2_id" json:"team_2_id"`
	Score       string    `db:"score" json:"score"`
	Current     bool      `db:"current" json:"current"`
	Played      bool      `db:"played" json:"played"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}
