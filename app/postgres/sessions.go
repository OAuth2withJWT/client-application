package postgres

import (
	"database/sql"
	"time"

	"github.com/OAuth2withJWT/client-application/app"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) CreateSession(sessionID string, accessToken string, idToken string, expiresAt time.Time) (string, error) {
	err := sr.db.QueryRow(
		"INSERT INTO sessions (session_id, access_token, id_token, expires_at) VALUES ($1, $2, $3, $4) RETURNING session_id",
		sessionID, accessToken, idToken, expiresAt,
	).Scan(&sessionID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sr *SessionRepository) UpdateStatus(sessionID string) error {
	query := `UPDATE sessions SET status = 'inactive' WHERE session_id = $1`
	_, err := sr.db.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionRepository) GetSessionByID(sessionID string) (app.Session, error) {
	var session app.Session
	err := sr.db.QueryRow(
		"SELECT id, session_id, access_token, id_token, expires_at FROM sessions WHERE session_id = $1",
		sessionID,
	).Scan(&session.Id, &session.SessionId, &session.AccessToken, &session.IdToken, &session.ExpiresAt)
	if err != nil {
		return app.Session{}, err
	}
	return session, nil
}

func (sr *SessionRepository) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE session_id = $1`
	_, err := sr.db.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}
