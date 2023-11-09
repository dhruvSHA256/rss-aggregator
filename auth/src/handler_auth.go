package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var secretKey = []byte("your-secret-key")

func (apiCfg apiConfig) handleAuth(w http.ResponseWriter, r *http.Request) {
	userId, err := decodeJWT(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Invalid token: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByID(r.Context(), userId)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
		return
	}
	respondWithJson(w, 200, user)
}

// Claims represents the JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
	jwt.StandardClaims
}

func decodeJWT(headers http.Header) (uuid.UUID, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return uuid.UUID{}, errors.New("auth header missing")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return uuid.UUID{}, errors.New("invalid auth header")
	}
	token := vals[1]
	validatedJWT, err := validateJWT(token)
	if err != nil {
		return uuid.UUID{}, err
	}
	return validatedJWT.UserID, nil
}

func GenerateJWT(userID uuid.UUID, password string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 1 day
	claims := &Claims{
		UserID:   userID,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error generating JWT: %v", err)
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("error extracting JWT claims")
	}

	return claims, nil
}
