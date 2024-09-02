package app

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Session struct {
	Id          int
	SessionId   string
	AccessToken string
	IdToken     string
	ExpiresAt   time.Time
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

type SessionService struct {
	repository SessionRepository
}

func NewSessionService(sr SessionRepository) *SessionService {
	return &SessionService{
		repository: sr,
	}
}

type SessionRepository interface {
	CreateSession(sessionID string, accessToken string, idToken string, expiresAt time.Time) (string, error)
	GetSessionByID(sessionID string) (Session, error)
	UpdateStatus(sessionID string) error
	DeleteSession(sessionID string) error
}

func (s *SessionService) UpdateStatus(sessionID string) error {
	err := s.repository.UpdateStatus(sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionService) ValidateSession(sessionID string) (Session, error) {
	session, err := s.GetSessionByID(sessionID)
	if err != nil || session.SessionId == "" {
		return Session{}, err
	}

	nowUnix := time.Now().Unix()
	expiresAtUnix := session.ExpiresAt.Add(-2 * time.Hour).Unix()

	if expiresAtUnix <= nowUnix {
		s.repository.DeleteSession(sessionID)
		return Session{}, err
	}

	return session, nil
}

func (s *SessionService) generateSessionID() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func (s *SessionService) CreateSession(accessToken string, idToken string, expiresAt time.Time) (string, error) {
	sessionID, err := s.generateSessionID()
	if err != nil {
		return "", err
	}

	sessionID, err = s.repository.CreateSession(sessionID, accessToken, idToken, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *SessionService) GetSessionByID(sessionID string) (Session, error) {
	session, err := s.repository.GetSessionByID(sessionID)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}
