package tests

import (
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"github.com/MorZLE/auth/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test"

	passDefaultLen = 10
)

func TestRegister_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	login := gofakeit.Name()
	pass := RandomPassword()

	reqReg := authv1.RegisterRequest{
		Login:    login,
		Password: pass,
		AppId:    appID,
	}

	respReg, err := st.AuthClient.Register(ctx, &reqReg)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	reqLogin := authv1.LoginRequest{
		Login:    login,
		Password: pass,
		AppId:    appID,
	}

	respLog, err := st.AuthClient.Login(ctx, &reqLogin)
	require.NoError(t, err)
	loginTime := time.Now()

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	tokenParset, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParset.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))

	const deltaSec = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.GRPC.Timeout).Unix(), claims["exp"].(float64), deltaSec)
}

func TestDuplicateRegister_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	login := gofakeit.Name()
	pass := RandomPassword()

	reqReg := authv1.RegisterRequest{
		Login:    login,
		Password: pass,
		AppId:    appID,
	}

	respReg, err := st.AuthClient.Register(ctx, &reqReg)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respDuble, err := st.AuthClient.Register(ctx, &reqReg)
	require.Error(t, err)
	assert.Empty(t, respDuble.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

	reqLogin := authv1.LoginRequest{
		Login:    login,
		Password: pass,
		AppId:    appID,
	}

	respLog, err := st.AuthClient.Login(ctx, &reqLogin)
	require.NoError(t, err)
	loginTime := time.Now()

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	tokenParset, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParset.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))

	const deltaSec = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.GRPC.Timeout).Unix(), claims["exp"].(float64), deltaSec)
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		name    string
		login   string
		pass    string
		wantErr string
	}{
		{
			name:    "register empty login",
			login:   "",
			pass:    RandomPassword(),
			wantErr: "data not exist",
		},
		{
			name:    "register empty pass",
			login:   gofakeit.Name(),
			pass:    "",
			wantErr: "data not exist",
		},
		{
			name:    "register empty",
			login:   "",
			pass:    "",
			wantErr: "data not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqReg := authv1.RegisterRequest{
				Login:    tt.login,
				Password: tt.pass,
				AppId:    appID,
			}

			_, err := st.AuthClient.Register(ctx, &reqReg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		name    string
		login   string
		pass    string
		wantErr string
		appid   int32
	}{
		{
			name:    "login empty login",
			login:   "",
			pass:    RandomPassword(),
			wantErr: "data not exist",
			appid:   appID,
		},
		{
			name:    "login empty pass",
			login:   gofakeit.Name(),
			pass:    "",
			wantErr: "data not exist",
			appid:   appID,
		},
		{
			name:    "login empty",
			login:   "",
			pass:    "",
			wantErr: "data not exist",
			appid:   appID,
		},
		{
			name:    "appid empty",
			login:   gofakeit.Name(),
			pass:    RandomPassword(),
			wantErr: "data not exist",
			appid:   emptyAppID,
		},
		{
			name:    "bad appid empty",
			login:   gofakeit.Name(),
			pass:    RandomPassword(),
			wantErr: "internal cerror",
			appid:   124,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqReg := authv1.LoginRequest{
				Login:    tt.login,
				Password: tt.pass,
				AppId:    tt.appid,
			}

			_, err := st.AuthClient.Login(ctx, &reqReg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func RandomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
