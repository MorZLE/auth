package tests

import (
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"github.com/MorZLE/auth/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	key = "test"
)

func TestCreate_Admin_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	appID := appID
	login := gofakeit.Name()
	lvl := gofakeit.Int32()
	key := key

	createReq := authv1.CreateAdminRequest{
		Login: login,
		Lvl:   lvl,
		Key:   key,
		AppId: int32(appID),
	}

	respReg, err := st.AuthClient.CreateAdmin(ctx, &createReq)
	require.NoError(t, err)
	require.NotEmpty(t, respReg)
}

func TestCreate_AddApp_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	login := gofakeit.Name()
	key := key
	secret := gofakeit.Name()

	createReq := authv1.AddAppRequest{
		Name:   login,
		Key:    key,
		Secret: secret,
	}

	respReg, err := st.AuthClient.AddApp(ctx, &createReq)
	require.NoError(t, err)
	require.NotEmpty(t, respReg)
}

func TestCreate_IsAdmin_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	appID := appID
	login := gofakeit.Name()
	lvl := gofakeit.Int32()
	pass := RandomPassword()
	key := key

	reqReg := authv1.RegisterRequest{
		Login:    login,
		Password: pass,
		AppId:    int32(appID),
	}

	respReg, err := st.AuthClient.Register(ctx, &reqReg)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	createReq := authv1.CreateAdminRequest{
		Login: login,
		Lvl:   lvl,
		Key:   key,
		AppId: int32(appID),
	}

	respAdm, err := st.AuthClient.CreateAdmin(ctx, &createReq)
	require.NoError(t, err)
	require.NotEmpty(t, respAdm)

	checkIsAdmin := authv1.IsAdminRequest{
		UserId: int32(respReg.UserId),
		AppId:  int32(appID),
	}

	respCheck, err := st.AuthClient.IsAdmin(ctx, &checkIsAdmin)
	require.NoError(t, err)
	require.True(t, respCheck.IsAdmin)

}

func Test_DeleteAdmin(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	appID := appID
	login := gofakeit.Name()
	lvl := gofakeit.Int32()
	pass := RandomPassword()
	key := key

	reqReg := authv1.RegisterRequest{
		Login:    login,
		Password: pass,
		AppId:    int32(appID),
	}

	respReg, err := st.AuthClient.Register(ctx, &reqReg)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	createReq := authv1.CreateAdminRequest{
		Login: login,
		Lvl:   lvl,
		Key:   key,
		AppId: int32(appID),
	}

	respAdm, err := st.AuthClient.CreateAdmin(ctx, &createReq)
	require.NoError(t, err)
	require.NotEmpty(t, respAdm)

	checkIsAdmin := authv1.IsAdminRequest{
		UserId: int32(respReg.UserId),
		AppId:  int32(appID),
	}

	respCheck, err := st.AuthClient.IsAdmin(ctx, &checkIsAdmin)
	require.NoError(t, err)
	assert.Equal(t, lvl, respCheck.Lvl)
	require.True(t, respCheck.IsAdmin)

	deleteReq := authv1.DeleteAdminRequest{
		Login: login,
		Key:   key,
	}
	respDelete, err := st.AuthClient.DeleteAdmin(ctx, &deleteReq)
	require.NoError(t, err)
	require.True(t, respDelete.Result)

	respCheck2, err := st.AuthClient.IsAdmin(ctx, &checkIsAdmin)
	assert.ErrorContains(t, err, "user not admin")
	require.Empty(t, respCheck2)

}
