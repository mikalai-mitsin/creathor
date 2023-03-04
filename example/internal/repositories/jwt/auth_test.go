package jwt

import (
	"context"
	"crypto"
	"crypto/rsa"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/example/internal/configs"
	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/pkg/clock"
	mock_clock "github.com/018bf/example/pkg/clock/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
)

func TestAuthRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Now()).AnyTimes()
	config := configs.NewMockConfig(t)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(config.Auth.PublicKey))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Auth.PrivateKey))
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		in0  context.Context
		user *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  1000 * time.Hour,
				refreshTTL: 1000 * time.Hour,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
			args: args{
				in0:  nil,
				user: mock_models.NewUser(t),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			got, err := r.Create(tt.args.in0, tt.args.user)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := r.validate(got.Access); err != nil {
				t.Errorf("validate() error = %v, wantErr %v", err, nil)
			}
		})
	}
}

func TestAuthRepository_GetSubject(t *testing.T) {
	privatePEM := []byte(
		"-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
	)
	publicPEM := []byte(
		"-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
	)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		in0   context.Context
		token models.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  0,
				refreshTTL: 0,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				in0:   context.Background(),
				token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsiYWNjZXNzIl0sImV4cCI6MTc0NzY1MjQzOCwibmJmIjoxNjYxMjUyNDM4LCJpYXQiOjE2NjEyNTI0MzgsImp0aSI6ImViY2M0MDUwLTU3YzMtNGVlMy1hNjMzLWY1NzgyOTc0MjRjYSJ9.N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo",
			},
			want:    "364b8b0d-4d85-4b23-9d77-054f0690e498",
			wantErr: nil,
		},
		{
			name: "invalid token",
			fields: fields{
				accessTTL:  0,
				refreshTTL: 0,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				in0:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhY2Nlc3MiLCJqdGkiOiJwMzE3OWQwZC0zM2VkLTQ0ZjYtOGEyMS05MzExZGExMWU3OWYiLCJzdWIiOiJjYzRhMzU1YS05NDdmLTRmYTQtYWY0Ny00ZmZmMWVhZTA2YTEifQ.6-U_fSeDCLfXtYlJ6aivTw-4O3-EuVVw2KYWNNay5zU",
			},
			want: "",
			wantErr: &errs.Error{
				Code:    16,
				Message: "Invalid token.",
				Params:  map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			got, err := r.GetSubject(tt.args.in0, tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetSubject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSubject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthRepository_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Date(2022, 2, 22, 0, 0, 0, 0, time.UTC)).AnyTimes()
	privatePEM := []byte(
		"-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
	)
	publicPEM := []byte(
		"-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
	)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		in0   context.Context
		token models.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  time.Second * 172800000,
				refreshTTL: time.Second * 172800000,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
			args: args{
				in0:   context.Background(),
				token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsicmVmcmVzaCJdLCJleHAiOjE4MzQwNTI0MzgsIm5iZiI6MTY2MTI1MjQzOCwiaWF0IjoxNjYxMjUyNDM4LCJqdGkiOiI2MTg3YTI1Ni0zZjg3LTQwZDktOGFhZi01MzlmMzEzNDFkNTYifQ.AWQ9-QPvCoqRIq1I8_xGfCyv8KDgDaQBiGrtqmnJeW78_4HKcjTqMG6E30O6R2nV7XuVhz-ZRmLHgng2pdCem2L6BC2maPSI4U9YZ9Lh_Bp-K_niS-cgqFxm6akjHpFAxzE2Y-7ZqE30Wp_U2B3ts9-6EL7tyx1FMQZ_lgIGeto",
			},
			wantErr: nil,
		},
		{
			name: "not refresh",
			fields: fields{
				accessTTL:  time.Second * 999999999,
				refreshTTL: time.Second * 999999999,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
			args: args{
				in0:   context.Background(),
				token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsiYWNjZXNzIl0sImV4cCI6MTc0NzY1MjQzOCwibmJmIjoxNjYxMjUyNDM4LCJpYXQiOjE2NjEyNTI0MzgsImp0aSI6ImViY2M0MDUwLTU3YzMtNGVlMy1hNjMzLWY1NzgyOTc0MjRjYSJ9.N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo",
			},
			wantErr: errs.NewBadToken(),
		},
		{
			name: "invalid",
			fields: fields{
				accessTTL:  time.Second * 999999999,
				refreshTTL: time.Second * 999999999,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
			args: args{
				in0:   context.Background(),
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJyZWZyZXNoIiwiZXhwIjoxNj1Q1NzE4NjcwLCJqdGkiOiI4NGVkN2Y0My01YTZiLTQzNjQtODk2Mi05ZWFiZDU1YTVlNjciLCJpYXQiOjE2NDU1NDU4NzAsIm5iZiI6MTY0NTU0NTg3MCwic3ViIjoiY2M0YTM1NWEtOTQ3Zi00ZmE0LWFmNDctNGZmZjFlYWUwNmExIn0.QduyrjDwcIyaywheuJiaKG59q0mVo8qjjCmWnAzuy8Q",
			},
			wantErr: errs.NewBadToken(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			_, err := r.RefreshToken(tt.args.in0, tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthRepository_Validate(t *testing.T) {
	privatePEM := []byte(
		"-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
	)
	publicPEM := []byte(
		"-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
	)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		in0   context.Context
		token models.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  0,
				refreshTTL: 0,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsiYWNjZXNzIl0sImV4cCI6MTc0NzY1MjQzOCwibmJmIjoxNjYxMjUyNDM4LCJpYXQiOjE2NjEyNTI0MzgsImp0aSI6ImViY2M0MDUwLTU3YzMtNGVlMy1hNjMzLWY1NzgyOTc0MjRjYSJ9.N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo",
			},
			wantErr: nil,
		},
		{
			name: "refresh",
			fields: fields{
				accessTTL:  time.Second * 172800000,
				refreshTTL: time.Second * 172800000,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsicmVmcmVzaCJdLCJleHAiOjE4MzQwNTI0MzgsIm5iZiI6MTY2MTI1MjQzOCwiaWF0IjoxNjYxMjUyNDM4LCJqdGkiOiI2MTg3YTI1Ni0zZjg3LTQwZDktOGFhZi01MzlmMzEzNDFkNTYifQ.AWQ9-QPvCoqRIq1I8_xGfCyv8KDgDaQBiGrtqmnJeW78_4HKcjTqMG6E30O6R2nV7XuVhz-ZRmLHgng2pdCem2L6BC2maPSI4U9YZ9Lh_Bp-K_niS-cgqFxm6akjHpFAxzE2Y-7ZqE30Wp_U2B3ts9-6EL7tyx1FMQZ_lgIGeto",
			},
			wantErr: errs.NewBadToken(),
		},
		{
			name: "invalid",
			fields: fields{
				accessTTL:  0,
				refreshTTL: 0,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiO1JhY2Nlc3MiLCJqdGkiOiJjMzE3OWQwZC0zM2VkLTQ0ZjYtOGEyMS05MzExZGExMWU3OWYiLCJzdWIiOiJjYzRhMzU1YS05NDdmLTRmYTQtYWY0Ny00ZmZmMWVhZTA2YTEifQ.6-U_fSeDCLfXtYlJ6aivTw-4O3-EuVVw2KYWNNay5zU",
			},
			wantErr: errs.NewBadToken(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			if err := r.Validate(tt.args.in0, tt.args.token); !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthRepository_createPair(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Date(2022, 2, 22, 0, 0, 0, 0, time.UTC)).AnyTimes()
	privatePEM := []byte(
		"-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
	)
	publicPEM := []byte(
		"-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
	)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		subject string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  time.Second * 999999999,
				refreshTTL: time.Second * 999999999,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
			args: args{
				subject: "asd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			got := r.createPair(tt.args.subject)
			if _, err := r.validate(got.Access); err != nil {
				t.Errorf("createPair() error = %v, wantErr %v", err, nil)
			}
		})
	}
}

func TestAuthRepository_keyFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	privatePEM := []byte(
		"-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
	)
	publicPEM := []byte(
		"-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
	)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		in0 *jwt.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  123,
				refreshTTL: 53245,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
			args: args{
				in0: &jwt.Token{
					Raw:       "",
					Method:    nil,
					Header:    nil,
					Claims:    nil,
					Signature: "",
					Valid:     false,
				},
			},
			want:    publicKey,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			got, err := r.keyFunc(tt.args.in0)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("keyFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyFunc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthRepository_validate(t *testing.T) {
	privatePEM := []byte(
		"-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
	)
	publicPEM := []byte(
		"-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
	)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	type fields struct {
		accessTTL  time.Duration
		refreshTTL time.Duration
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		clock      clock.Clock
		logger     log.Logger
	}
	type args struct {
		token models.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *jwt.Token
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessTTL:  0,
				refreshTTL: 0,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsiYWNjZXNzIl0sImV4cCI6MTc0NzY1MjQzOCwibmJmIjoxNjYxMjUyNDM4LCJpYXQiOjE2NjEyNTI0MzgsImp0aSI6ImViY2M0MDUwLTU3YzMtNGVlMy1hNjMzLWY1NzgyOTc0MjRjYSJ9.N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo",
			},
			want: &jwt.Token{
				Raw: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsiYWNjZXNzIl0sImV4cCI6MTc0NzY1MjQzOCwibmJmIjoxNjYxMjUyNDM4LCJpYXQiOjE2NjEyNTI0MzgsImp0aSI6ImViY2M0MDUwLTU3YzMtNGVlMy1hNjMzLWY1NzgyOTc0MjRjYSJ9.N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo",
				Method: &jwt.SigningMethodRSA{
					Name: "RS512",
					Hash: crypto.SHA512,
				},
				Header: map[string]interface{}{
					"alg": "RS512",
					"typ": "JWT",
				},
				Claims: jwt.MapClaims{
					"aud": []any{"access"},
					"exp": float64(1747652438),
					"nbf": float64(1661252438),
					"iat": float64(1661252438),
					"jti": "ebcc4050-57c3-4ee3-a633-f578297424ca",
					"sub": "364b8b0d-4d85-4b23-9d77-054f0690e498",
				},
				Signature: "N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo",
				Valid:     true,
			},
			wantErr: nil,
		},
		{
			name: "invalid",
			fields: fields{
				accessTTL:  0,
				refreshTTL: 0,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      nil,
				logger:     nil,
			},
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiO1JhY2Nlc3MiLCJqdGkiOiJjMzE3OWQwZC0zM2VkLTQ0ZjYtOGEyMS05MzExZGExMWU3OWYiLCJzdWIiOiJjYzRhMzU1YS05NDdmLTRmYTQtYWY0Ny00ZmZmMWVhZTA2YTEifQ.6-U_fSeDCLfXtYlJ6aivTw-4O3-EuVVw2KYWNNay5zU",
			},
			want:    nil,
			wantErr: errs.NewBadToken(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepository{
				accessTTL:  tt.fields.accessTTL,
				refreshTTL: tt.fields.refreshTTL,
				publicKey:  tt.fields.publicKey,
				privateKey: tt.fields.privateKey,
				clock:      tt.fields.clock,
				logger:     tt.fields.logger,
			}
			got, err := r.validate(tt.args.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAuthRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	config := configs.NewMockConfig(t)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(config.Auth.PublicKey))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Auth.PrivateKey))
	type args struct {
		config *configs.Config
		clock  clock.Clock
		logger log.Logger
	}
	tests := []struct {
		name string
		args args
		want repositories.AuthRepository
	}{
		{
			name: "ok",
			args: args{
				config: config,
				clock:  mockClock,
				logger: logger,
			},
			want: &AuthRepository{
				accessTTL:  86400 * time.Second,
				refreshTTL: 172800 * time.Second,
				publicKey:  publicKey,
				privateKey: privateKey,
				clock:      mockClock,
				logger:     logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthRepository(tt.args.config, tt.args.clock, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
