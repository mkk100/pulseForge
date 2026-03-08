package httpapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultJWTSecret = "pulseforge-dev-secret"

type contextKey string

const userIDContextKey contextKey = "userID"

type tokenManager struct {
	secret []byte
}

type tokenClaims struct {
	UserID int64 `json:"uid"`
	Exp    int64 `json:"exp"`
}

func newTokenManager(secret string) *tokenManager {
	if secret == "" {
		secret = defaultJWTSecret
	}
	return &tokenManager{secret: []byte(secret)}
}

func (m *tokenManager) issueToken(userID int64) (string, error) {
	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}
	payLoadJSON, err := json.Marshal(tokenClaims{
		UserID: userID,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		return "", err
	}
	
	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadPart := base64.RawURLEncoding.EncodeToString(payLoadJSON)
	signingInput := headerPart + "." + payloadPart

	signature := m.sign(signingInput)
	return signingInput + "." + signature, nil
}

func (m *tokenManager) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := m.parseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(withUserID(r.Context(), userID)))
	})
}

func userIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	return userID, ok
}

func withUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func (m *tokenManager) parseToken(token string) (int64, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return 0, errors.New("token must have 3 parts")
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := m.sign(signingInput)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return 0, errors.New("invalid signature")
	}
	
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, fmt.Errorf("decode payload: %w", err)
	}
	
	var claims tokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return 0, fmt.Errorf("decode claims: %w", err)
	}
	if claims.UserID <= 0 {
		return 0, errors.New("missing user id")
	}
	if time.Now().Unix() > claims.Exp {
		return 0, errors.New("token expired")
	}
	return claims.UserID, nil
}

func (m *tokenManager) sign(input string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}