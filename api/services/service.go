package services

import (
	"crypto/rsa"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/envs"
	"github.com/shellhub-io/shellhub/pkg/geoip"
	"github.com/shellhub-io/shellhub/pkg/validator"
	log "github.com/sirupsen/logrus"
)

type APIService struct {
	*service
}

var _ Service = (*APIService)(nil)

type service struct {
	store     store.Store
	privKey   *rsa.PrivateKey
	pubKey    *rsa.PublicKey
	cache     cache.Cache
	client    interface{}
	locator   geoip.Locator
	validator *validator.Validator
	cfg       *config
}

type config struct {
	// Specifies the maximum duration in minutes for which a user can be blocked from login attempts.
	// The default value is 32768, equivalent to 15 days.
	MaximumLoginTimeout int `env:"MAXIMUM_LOGIN_TIMEOUT,default=0"`
}

//go:generate mockery --name Service --filename services.go
type Service interface {
	BillingInterface
	TagsService
	DeviceService
	DeviceTags
	UserService
	SSHKeysService
	SSHKeysTagsService
	SessionService
	NamespaceService
	AuthService
	StatsService
	SetupService
	SystemService
	APIKeyService
}

func NewService(store store.Store, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, cache cache.Cache, c interface{}, l geoip.Locator) *APIService {
	if privKey == nil || pubKey == nil {
		var err error
		privKey, pubKey, err = LoadKeys()
		if err != nil {
			panic(err)
		}
	}

	cfg, err := envs.ParseWithPrefix[config]("API_")
	if err != nil {
		log.WithError(err).Fatal("Failed to load environment variables")
	}

	return &APIService{service: &service{
		client:    c,
		locator:   l,
		store:     store,
		privKey:   privKey,
		pubKey:    pubKey,
		cache:     cache,
		validator: validator.New(),
		cfg:       cfg,
	}}
}
