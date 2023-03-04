package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/018bf/example/internal/configs"
	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
	"github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
)

const refreshAudience = "refresh"
const accessAudience = "access"

type AuthRepository struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	clock      clock.Clock
	logger     log.Logger
}

func NewAuthRepository(
	config *configs.Config,
	clock clock.Clock,
	logger log.Logger,
) repositories.AuthRepository {
	private, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Auth.PrivateKey))
	if err != nil {
		panic(err)
	}
	public, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.Auth.PublicKey))
	if err != nil {
		panic(err)
	}
	return &AuthRepository{
		accessTTL:  time.Duration(config.Auth.AccessTTL) * time.Second,
		refreshTTL: time.Duration(config.Auth.RefreshTTL) * time.Second,
		publicKey:  public,
		privateKey: private,
		clock:      clock,
		logger:     logger,
	}
}

func (r *AuthRepository) Create(_ context.Context, user *models.User) (*models.TokenPair, error) {
	pair := r.createPair(string(user.ID))
	return pair, nil
}

func (r *AuthRepository) createPair(subject string) *models.TokenPair {
	now := r.clock.Now().UTC()
	accessClaims := jwt.RegisteredClaims{
		Audience:  []string{accessAudience},
		ExpiresAt: jwt.NewNumericDate(now.Add(r.accessTTL)),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "",
		NotBefore: jwt.NewNumericDate(now),
		Subject:   subject,
	}
	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), accessClaims)
	accessTokenString, err := accessToken.SignedString(r.privateKey)
	if err != nil {
		return nil
	}

	refreshClaims := jwt.RegisteredClaims{
		Audience:  []string{refreshAudience},
		ExpiresAt: jwt.NewNumericDate(now.Add(r.refreshTTL)),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "",
		NotBefore: jwt.NewNumericDate(now),
		Subject:   subject,
	}
	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(r.privateKey)
	if err != nil {
		return nil
	}
	return &models.TokenPair{
		Access:  models.Token(accessTokenString),
		Refresh: models.Token(refreshTokenString),
	}
}

func (r *AuthRepository) Validate(_ context.Context, token models.Token) error {
	jwtToken, err := r.validate(token)
	if err != nil {
		return err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if !claims.VerifyAudience(accessAudience, true) {
		return errs.NewBadToken()
	}
	return nil
}

func (r *AuthRepository) RefreshToken(
	_ context.Context,
	token models.Token,
) (*models.TokenPair, error) {
	jwtToken, err := r.validate(token)
	if err != nil {
		return nil, err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if !claims.VerifyAudience(refreshAudience, true) {
		return nil, errs.NewBadToken()
	}
	pair := r.createPair(fmt.Sprint(claims["sub"]))
	return pair, nil
}

func (r *AuthRepository) validate(token models.Token) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token.String(), r.keyFunc)
	if err != nil {
		e := errs.NewBadToken()
		return nil, e
	}
	return jwtToken, nil
}

func (r *AuthRepository) GetSubject(_ context.Context, token models.Token) (string, error) {
	jwtToken, err := jwt.Parse(token.String(), r.keyFunc)
	if err != nil {
		e := errs.NewError(errs.ErrorCodeUnauthenticated, "Invalid token.")
		return "", e
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	return fmt.Sprint(claims["sub"]), nil
}

func (r *AuthRepository) keyFunc(_ *jwt.Token) (interface{}, error) {
	return r.publicKey, nil
}
