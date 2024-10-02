package pow

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// ChallengeStore defines methods for challenge persistence
type ChallengeStore interface {
	SaveChallenge(challenge *Challenge) error
	GetChallenge(id string) (*Challenge, error)
	DeleteChallenge(id string) error
}

// SQLiteChallengeStore implements ChallengeStore using SQLite
type SQLiteChallengeStore struct {
	DB *sql.DB
}

func NewSQLiteChallengeStore(dbPath string) (*SQLiteChallengeStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	store := &SQLiteChallengeStore{DB: db}
	err = store.initialize()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *SQLiteChallengeStore) initialize() error {
	query := `
    CREATE TABLE IF NOT EXISTS challenges (
        id TEXT PRIMARY KEY,
        algorithm TEXT,
        difficulty INTEGER,
        data TEXT,
        created_at DATETIME
    );
    `
	_, err := s.DB.Exec(query)
	return err
}

func (s *SQLiteChallengeStore) SaveChallenge(challenge *Challenge) error {
	query := `
    INSERT INTO challenges (id, algorithm, difficulty, data, created_at)
    VALUES (?, ?, ?, ?, ?);
    `
	_, err := s.DB.Exec(query, challenge.ID, challenge.Algorithm, challenge.Difficulty, challenge.Data, challenge.CreatedAt)
	return err
}

func (s *SQLiteChallengeStore) GetChallenge(id string) (*Challenge, error) {
	query := `
    SELECT id, algorithm, difficulty, data, created_at
    FROM challenges
    WHERE id = ?;
    `
	row := s.DB.QueryRow(query, id)
	challenge := &Challenge{}
	err := row.Scan(&challenge.ID, &challenge.Algorithm, &challenge.Difficulty, &challenge.Data, &challenge.CreatedAt)
	if err != nil {
		return nil, err
	}
	return challenge, nil
}

func (s *SQLiteChallengeStore) DeleteChallenge(id string) error {
	query := `DELETE FROM challenges WHERE id = ?;`
	_, err := s.DB.Exec(query, id)
	return err
}
