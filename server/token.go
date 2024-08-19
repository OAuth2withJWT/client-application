package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserInfo struct {
	ID    int
	Name  string
	Email string
}

func (s *Server) GetUserInfoFromIDToken(idToken string) (*UserInfo, error) {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.IDPConfig.PublicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse ID token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid ID token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from ID token")
	}

	userIdStr, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("user ID (sub claim) is missing in ID token")
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user ID to integer: %w", err)
	}

	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)

	return &UserInfo{
		ID:    userId,
		Name:  name,
		Email: email,
	}, nil
}

func (s *Server) ValidateIDToken(idToken, accessToken string) error {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.IDPConfig.PublicKey, nil
	})

	if err != nil {
		return fmt.Errorf("failed to verify ID token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("failed to extract claims or token is invalid")
	}

	if iss, ok := claims["iss"].(string); !ok || iss != s.IDPConfig.IdentityProvider {
		return fmt.Errorf("iss claim does not match expected value")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != s.IDPConfig.ClientID {
		return fmt.Errorf("aud claim does not match expected value")
	}

	if iat, ok := claims["iat"].(float64); !ok || time.Unix(int64(iat), 0).After(time.Now()) {
		return fmt.Errorf("iat claim is invalid")
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return fmt.Errorf("exp claim is invalid or token has expired")
	}

	atHash, ok := claims["at_hash"].(string)
	if !ok {
		return fmt.Errorf("at_hash claim is missing in ID token")
	}

	hash := sha256.Sum256([]byte(accessToken))
	computedAtHash := base64.RawURLEncoding.EncodeToString(hash[:len(hash)/2])

	if computedAtHash != atHash {
		return fmt.Errorf("at_hash claim does not match")
	}

	return nil
}
