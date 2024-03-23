package services

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"

	"github.com/shellhub-io/shellhub/pkg/api/internalclient/mocks"
	"github.com/shellhub-io/shellhub/pkg/clock"
	clockmocks "github.com/shellhub-io/shellhub/pkg/clock/mocks"
	"github.com/shellhub-io/shellhub/pkg/envs"
	env_mocks "github.com/shellhub-io/shellhub/pkg/envs/mocks"
	"github.com/shellhub-io/shellhub/pkg/password"
	passwordmock "github.com/shellhub-io/shellhub/pkg/password/mocks"
	gomock "github.com/stretchr/testify/mock"
)

var (
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	clientMock   *mocks.Client
	envMock      *env_mocks.Backend
	clockMock    *clockmocks.Clock
	passwordMock *passwordmock.Password
	now          time.Time
)

func TestMain(m *testing.M) {
	privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	publicKey = &privateKey.PublicKey
	clientMock = &mocks.Client{}
	clockMock = &clockmocks.Clock{}
	clock.DefaultBackend = clockMock

	envMock = &env_mocks.Backend{}
	envs.DefaultBackend = envMock
	// Mock ParseWithPrefix in NewService globally
	cfg := new(config)
	envMock.On("Process", "API_", cfg).Return(nil).Run(func(args gomock.Arguments) {
		cfg := args.Get(1).(*config)
		cfg.MaximumLoginTimeout = 0
	})

	passwordMock = &passwordmock.Password{}
	password.Backend = passwordMock
	now = time.Now()
	code := m.Run()
	os.Exit(code)
}
