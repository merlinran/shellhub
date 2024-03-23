package services

import (
	"context"
	"crypto/rsa"
	"math"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/shellhub-io/shellhub/pkg/clock"
)

func LoadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signBytes, err := os.ReadFile(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, nil, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, nil, err
	}

	verifyBytes, err := os.ReadFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

// TODO: put it inside s.cache.ID()
func redisID(kind, ip, id string) string {
	return kind + ":" + ip + "-" + id
}

// isBlocked reports whether the sourceIP is currently blocked from attempting to
// log in to a user with the specified id. It returns the absolute Unix timestamp
// in seconds representing the end of the block, or 0 if no block was found.
func (s *service) isBlocked(ctx context.Context, sourceIP, id string) (int64, error) {
	x := int64(0)
	if err := s.cache.Get(ctx, redisID("login-timeout", sourceIP, id), &x); err != nil {
		return 0, err
	}

	return x, nil
}

// storeAttempt stores a login attempt from sourceIP to the user with the specified id.
// If the attempt number equals or exceeds 3, it sets a timeout for future login attempts.
//
// The timeout is calculated based on the number of attempts using the following equation:
//
//	F(x) = c + min(4^(a - 3), M)
//
// Attempts TTL increases exponentially along with the timeout, being at least 1.5 times greater
// than the timeout and at most M. Attempts can never have a duration less than 2 minutes. This means that
// an attempt must last for:
//
//	F(y) = c + max(min((4^(a - 3) * 1.5), M), 2)
//
// Where:
//
//	x is the timeout duration in minutes.
//	y is the attempt duration in minutes.
//	c is the current Unix timestamp in seconds.
//	a is the attempt number.
//	M is the maximum timeout value, specified by the "SHELLHUB_MAXIMUM_LOGIN_TIMEOUT" environment variable.
//
// Examples for M equal to 32768 (15 days):
//
//	         timeout(x) | attemptTTL(y)
//	F(3)  = [ c + 1,      c + 2        ]
//	F(4)  = [ c + 4,      c + 6        ]
//	F(5)  = [ c + 16,     c + 24       ]
//	F(8)  = [ c + 1024,   c + 1536     ]
//	F(11) = [ c + 32768,  c + 32768    ]
func (s *service) storeAttempt(ctx context.Context, sourceIP, id string) (int64, error) {
	a := 0
	if err := s.cache.Get(ctx, redisID("login-attempt", sourceIP, id), &a); err != nil {
		return 0, err
	}

	M := s.cfg.MaximumLoginTimeout
	a += 1

	ttl := math.Max(math.Min(math.Pow(4, float64(a-3)), float64(M)), 2)

	y := float64(ttl) * 1.5
	if err := s.cache.Set(ctx, redisID("login-attempt", sourceIP, id), a, time.Duration(y)*time.Minute); err != nil {
		return 0, err
	}

	if a <= 2 {
		return 0, nil
	}

	c := clock.Now()
	dur := time.Duration(ttl) * time.Minute
	x := c.Add(dur).Unix()

	if err := s.cache.Set(ctx, redisID("login-timeout", sourceIP, id), x, dur); err != nil {
		return 0, err
	}

	return x, nil
}

// resetAttempts resets the login attempts and associated timeout from the sourceIP to
// the user with the specified ID.
func (s *service) resetAttempts(ctx context.Context, ip, id string) error {
	if err := s.cache.Delete(ctx, redisID("login-attempt", ip, id)); err != nil {
		return err
	}

	if err := s.cache.Delete(ctx, redisID("login-timeout", ip, id)); err != nil {
		return err
	}

	return nil
}
