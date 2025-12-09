package dao

import (
	"darbelis.eu/stabas/db"
	"darbelis.eu/stabas/entities"
	"database/sql"
	"errors"
	"fmt"
)

type LiteParticipantsRepository struct {
	database *db.Database
}

func NewLiteParticipantsRepository(database *db.Database) (*LiteParticipantsRepository, error) {
	repo := &LiteParticipantsRepository{database: database}
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (repo *LiteParticipantsRepository) initSchema() error {
	db, err := repo.database.GetDB()
	if err != nil {
		return err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS participants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		deleted BOOLEAN DEFAULT 0,
		token VARCHAR(255),
		password VARCHAR(32)
	);
	`
	_, err = db.Exec(schema)
	return err
}

func (repo *LiteParticipantsRepository) Close() error {
	return repo.database.Close()
}

func (repo *LiteParticipantsRepository) GetParticipants() []*entities.Participant {
	db, err := repo.database.GetDB()
	if err != nil {
		return []*entities.Participant{}
	}

	query := `
		SELECT id, name, deleted, token, password
		FROM participants
		WHERE deleted = 0
	`

	rows, err := db.Query(query)
	if err != nil {
		return []*entities.Participant{}
	}
	defer rows.Close()

	var participants []*entities.Participant
	for rows.Next() {
		var participant entities.Participant
		err := rows.Scan(
			&participant.Id,
			&participant.Name,
			&participant.Deleted,
			&participant.Token,
			&participant.Password,
		)
		if err != nil {
			continue
		}
		participants = append(participants, &participant)
	}

	return participants
}

func (repo *LiteParticipantsRepository) FindParticipant(id int) (*entities.Participant, error) {
	db, err := repo.database.GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, deleted, token, password
		FROM participants
		WHERE id = ?
	`

	var participant entities.Participant
	err = db.QueryRow(query, id).Scan(
		&participant.Id,
		&participant.Name,
		&participant.Deleted,
		&participant.Token,
		&participant.Password,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(fmt.Sprintf("participant with id %d not found", id))
	}
	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (repo *LiteParticipantsRepository) AddParticipant(participant *entities.Participant) (*entities.Participant, error) {
	db, err := repo.database.GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO participants (name, deleted)
		VALUES (?, ?)
	`

	result, err := db.Exec(query, participant.Name, participant.Deleted)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	participant.Id = int(id)
	return participant, nil
}

func (repo *LiteParticipantsRepository) RemoveParticipant(id int) error {
	db, err := repo.database.GetDB()
	if err != nil {
		return err
	}

	query := "UPDATE participants SET deleted = 1 WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf("participant with id %d not found", id))
	}

	return nil
}

func (repo *LiteParticipantsRepository) UpdateParticipant(participant *entities.Participant) error {
	db, err := repo.database.GetDB()
	if err != nil {
		return err
	}

	// Check if participant exists
	_, err = repo.FindParticipant(participant.Id)
	if err != nil {
		return err
	}

	query := `
		UPDATE participants
		SET name = ?
		WHERE id = ?
	`

	_, err = db.Exec(query, participant.Name, participant.Id)
	return err
}

func (repo *LiteParticipantsRepository) FindParticipantByName(name string) *entities.Participant {
	db, err := repo.database.GetDB()
	if err != nil {
		return nil
	}

	query := `
		SELECT id, name, deleted, token, password
		FROM participants
		WHERE name = ? AND deleted = 0
	`

	var participant entities.Participant
	err = db.QueryRow(query, name).Scan(
		&participant.Id,
		&participant.Name,
		&participant.Deleted,
		&participant.Token,
		&participant.Password,
	)

	if err != nil {
		return nil
	}

	return &participant
}

func (repo *LiteParticipantsRepository) FindParticipantByToken(token string) *entities.Participant {
	db, err := repo.database.GetDB()
	if err != nil {
		return nil
	}

	query := `
		SELECT id, name, deleted, token, password
		FROM participants
		WHERE token = ? AND deleted = 0
	`

	var participant entities.Participant
	err = db.QueryRow(query, token).Scan(
		&participant.Id,
		&participant.Name,
		&participant.Deleted,
		&participant.Token,
		&participant.Password,
	)

	if err != nil {
		return nil
	}

	return &participant
}

func (repo *LiteParticipantsRepository) UpdateParticipantToken(id int, token string) error {
	db, err := repo.database.GetDB()
	if err != nil {
		return err
	}

	// Check if participant exists
	_, err = repo.FindParticipant(id)
	if err != nil {
		return err
	}

	query := "UPDATE participants SET token = ? WHERE id = ?"
	_, err = db.Exec(query, token, id)
	return err
}

func (repo *LiteParticipantsRepository) UpdateParticipantPassword(id int, password string) error {
	db, err := repo.database.GetDB()
	if err != nil {
		return err
	}

	// Check if participant exists
	_, err = repo.FindParticipant(id)
	if err != nil {
		return err
	}

	query := "UPDATE participants SET password = ? WHERE id = ?"
	_, err = db.Exec(query, password, id)
	return err
}
