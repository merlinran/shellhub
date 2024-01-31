package mongo

import (
	"context"
	"testing"

	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/pkg/fixtures"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetStats(t *testing.T) {
	type Expected struct {
		stats *models.Stats
		err   error
	}

	cases := []struct {
		description string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds",
			fixtures: []string{
				fixtures.FixtureUsers,
				fixtures.FixtureNamespaces,
				fixtures.FixtureSessions,
				fixtures.FixtureActiveSessions,
				fixtures.FixtureDevices,
				fixtures.FixtureConnectedDevices,
			},
			expected: Expected{
				stats: &models.Stats{},
				err:   nil,
			},
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown()

			stats, err := mongostore.GetStats(context.TODO())
			assert.Equal(t, tc.expected, Expected{stats: stats, err: err})
		})
	}
}
