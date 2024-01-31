package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/pkg/fixtures"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/query"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestSessionList(t *testing.T) {
	type Expected struct {
		s     []models.Session
		count int
		err   error
	}

	cases := []struct {
		description string
		paginator   query.Paginator
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when sessions are found",
			paginator:   query.Paginator{Page: -1, PerPage: -1},
			fixtures: []string{
				fixtures.FixtureNamespaces,
				fixtures.FixtureDevices,
				fixtures.FixtureConnectedDevices,
				fixtures.FixtureSessions,
				fixtures.FixtureActiveSessions,
			},
			expected: Expected{
				s: []models.Session{
					{
						StartedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						LastSeen:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						UID:       "a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68",
						DeviceUID: "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
						TenantID:  "00000000-0000-4000-0000-000000000000",
						Username:  "john_doe",
						IPAddress: "0.0.0.0",
						Device: &models.Device{
							CreatedAt:        time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							StatusUpdatedAt:  time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							LastSeen:         time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							UID:              "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
							Name:             "device-3",
							Identity:         &models.DeviceIdentity{MAC: "mac-3"},
							Info:             nil,
							PublicKey:        "",
							TenantID:         "00000000-0000-4000-0000-000000000000",
							Online:           true,
							Namespace:        "namespace-1",
							Status:           "accepted",
							RemoteAddr:       "",
							Position:         nil,
							Tags:             []string{"tag-1"},
							PublicURL:        false,
							PublicURLAddress: "",
							Acceptable:       false,
						},
						Active:        true,
						Closed:        true,
						Authenticated: true,
						Recorded:      false,
						Type:          "shell",
						Term:          "xterm",
						Position:      models.SessionPosition{Longitude: 0, Latitude: 0},
					},
					{
						StartedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						LastSeen:  time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						UID:       "e7f3a56d8b9e1dc4c285c98c8ea9c33032a17bda5b6c6b05a6213c2a02f97824",
						DeviceUID: "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
						TenantID:  "00000000-0000-4000-0000-000000000000",
						Username:  "john_doe",
						IPAddress: "0.0.0.0",
						Device: &models.Device{
							CreatedAt:        time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							StatusUpdatedAt:  time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							LastSeen:         time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							UID:              "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
							Name:             "device-3",
							Identity:         &models.DeviceIdentity{MAC: "mac-3"},
							Info:             nil,
							PublicKey:        "",
							TenantID:         "00000000-0000-4000-0000-000000000000",
							Online:           true,
							Namespace:        "namespace-1",
							Status:           "accepted",
							RemoteAddr:       "",
							Position:         nil,
							Tags:             []string{"tag-1"},
							PublicURL:        false,
							PublicURLAddress: "",
							Acceptable:       false,
						},
						Active:        false,
						Closed:        true,
						Authenticated: true,
						Recorded:      true,
						Type:          "shell",
						Term:          "xterm",
						Position:      models.SessionPosition{Longitude: 45.6789, Latitude: -12.3456},
					},
					{
						StartedAt: time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						LastSeen:  time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						UID:       "fc2e1493d8b6a4c17bf6a2f7f9e55629e384b2d3a21e0c3d90f6e35b0c946178a",
						DeviceUID: "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
						TenantID:  "00000000-0000-4000-0000-000000000000",
						Username:  "john_doe",
						IPAddress: "0.0.0.0",
						Device: &models.Device{
							CreatedAt:        time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							StatusUpdatedAt:  time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							LastSeen:         time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							UID:              "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
							Name:             "device-3",
							Identity:         &models.DeviceIdentity{MAC: "mac-3"},
							Info:             nil,
							PublicKey:        "",
							TenantID:         "00000000-0000-4000-0000-000000000000",
							Online:           true,
							Namespace:        "namespace-1",
							Status:           "accepted",
							RemoteAddr:       "",
							Position:         nil,
							Tags:             []string{"tag-1"},
							PublicURL:        false,
							PublicURLAddress: "",
							Acceptable:       false,
						},
						Active:        false,
						Closed:        true,
						Authenticated: true,
						Recorded:      false,
						Type:          "exec",
						Term:          "",
						Position:      models.SessionPosition{Longitude: -78.9012, Latitude: 23.4567},
					},
					{
						StartedAt: time.Date(2023, 1, 4, 12, 0, 0, 0, time.UTC),
						LastSeen:  time.Date(2023, 1, 4, 12, 0, 0, 0, time.UTC),
						UID:       "bc3d75821a29cfe70bf7986f9ee5629e384b2d3a21e0c3d90f6e35b0c946178a",
						DeviceUID: "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
						TenantID:  "00000000-0000-4000-0000-000000000000",
						Username:  "john_doe",
						IPAddress: "0.0.0.0",
						Device: &models.Device{
							CreatedAt:        time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							StatusUpdatedAt:  time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							LastSeen:         time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
							UID:              "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
							Name:             "device-3",
							Identity:         &models.DeviceIdentity{MAC: "mac-3"},
							Info:             nil,
							PublicKey:        "",
							TenantID:         "00000000-0000-4000-0000-000000000000",
							Online:           true,
							Namespace:        "namespace-1",
							Status:           "accepted",
							RemoteAddr:       "",
							Position:         nil,
							Tags:             []string{"tag-1"},
							PublicURL:        false,
							PublicURLAddress: "",
							Acceptable:       false,
						},
						Active:        false,
						Closed:        true,
						Authenticated: true,
						Recorded:      true,
						Type:          "shell",
						Term:          "xterm",
						Position:      models.SessionPosition{Longitude: -56.7890, Latitude: 34.5678},
					},
				},
				count: 4,
				err:   nil,
			},
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("sessions")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			for _, session := range tc.expected.s {
				doc := bson.M{
					"started_at":    session.StartedAt,
					"last_seen":     session.LastSeen,
					"uid":           session.UID,
					"device_uid":    session.DeviceUID,
					"tenant_id":     session.TenantID,
					"username":      session.Username,
					"ip_address":    session.IPAddress,
					"active":        session.Active,
					"closed":        session.Closed,
					"authenticated": session.Authenticated,
					"recorded":      session.Recorded,
					"type":          session.Type,
					"term":          session.Term,
					"position": bson.M{
						"longitude": session.Position.Longitude,
						"latitude":  session.Position.Latitude,
					},
				}
				testData = append(testData, doc)
			}

			if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			s, count, err := mongostore.SessionList(context.TODO(), tc.paginator)
			assert.Equal(t, tc.expected, Expected{s: s, count: count, err: err})
		})
	}
}

func TestSessionGet(t *testing.T) {
	type Expected struct {
		s   *models.Session
		err error
	}

	cases := []struct {
		description string
		UID         models.UID
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when session is not found",
			UID:         models.UID("nonexistent"),
			fixtures: []string{
				fixtures.FixtureNamespaces,
				fixtures.FixtureDevices,
				fixtures.FixtureConnectedDevices,
				fixtures.FixtureSessions,
				fixtures.FixtureActiveSessions,
			},
			expected: Expected{
				s:   nil,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when session is found",
			UID:         models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
			fixtures: []string{
				fixtures.FixtureNamespaces,
				fixtures.FixtureDevices,
				fixtures.FixtureConnectedDevices,
				fixtures.FixtureSessions,
				fixtures.FixtureActiveSessions,
			},
			expected: Expected{
				s: &models.Session{
					StartedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastSeen:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UID:       "a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68",
					DeviceUID: "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
					TenantID:  "00000000-0000-4000-0000-000000000000",
					Username:  "john_doe",
					IPAddress: "0.0.0.0",
					Device: &models.Device{
						CreatedAt:        time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						StatusUpdatedAt:  time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						LastSeen:         time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						UID:              "2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c",
						Name:             "device-3",
						Identity:         &models.DeviceIdentity{MAC: "mac-3"},
						Info:             nil,
						PublicKey:        "",
						TenantID:         "00000000-0000-4000-0000-000000000000",
						Online:           true,
						Namespace:        "namespace-1",
						Status:           "accepted",
						RemoteAddr:       "",
						Position:         nil,
						Tags:             []string{"tag-1"},
						PublicURL:        false,
						PublicURLAddress: "",
						Acceptable:       false,
					},
					Active:        true,
					Closed:        true,
					Authenticated: true,
					Recorded:      false,
					Type:          "shell",
					Term:          "xterm",
					Position:      models.SessionPosition{Longitude: 0, Latitude: 0},
				},
				err: nil,
			},
		},
	}
	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("sessions")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			doc := bson.M{
				"started_at":    tc.expected.s.StartedAt,
				"last_seen":     tc.expected.s.LastSeen,
				"uid":           tc.expected.s.UID,
				"device_uid":    tc.expected.s.DeviceUID,
				"tenant_id":     tc.expected.s.TenantID,
				"username":      tc.expected.s.Username,
				"ip_address":    tc.expected.s.IPAddress,
				"active":        tc.expected.s.Active,
				"closed":        tc.expected.s.Closed,
				"authenticated": tc.expected.s.Authenticated,
				"recorded":      tc.expected.s.Recorded,
				"type":          tc.expected.s.Type,
				"term":          tc.expected.s.Term,
				"position": bson.M{
					"longitude": tc.expected.s.Position.Longitude,
					"latitude":  tc.expected.s.Position.Latitude,
				},
				"device": bson.M{
					"created_at":         tc.expected.s.Device.CreatedAt,
					"status_updated_at":  tc.expected.s.Device.StatusUpdatedAt,
					"last_seen":          tc.expected.s.Device.LastSeen,
					"uid":                tc.expected.s.Device.UID,
					"name":               tc.expected.s.Device.Name,
					"identity":           tc.expected.s.Device.Identity,
					"info":               tc.expected.s.Device.Info,
					"public_key":         tc.expected.s.Device.PublicKey,
					"tenant_id":          tc.expected.s.Device.TenantID,
					"online":             tc.expected.s.Device.Online,
					"namespace":          tc.expected.s.Device.Namespace,
					"status":             tc.expected.s.Device.Status,
					"remote_addr":        tc.expected.s.Device.RemoteAddr,
					"position":           tc.expected.s.Device.Position,
					"tags":               tc.expected.s.Device.Tags,
					"public_url":         tc.expected.s.Device.PublicURL,
					"public_url_address": tc.expected.s.Device.PublicURLAddress,
					"acceptable":         tc.expected.s.Device.Acceptable,
				},
			}
			testData = append(testData, doc)

			if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			s, err := mongostore.SessionGet(context.TODO(), tc.UID)
			assert.Equal(t, tc.expected, Expected{s: s, err: err})
		})
	}
}

func TestSessionCreate(t *testing.T) {
	cases := []struct {
		description string
		fixtures    []string
		session     models.Session
		expected    error
	}{
		{
			description: "",
			fixtures:    []string{fixtures.FixtureDevices, fixtures.FixtureNamespaces},
			session: models.Session{
				Username:      "username",
				UID:           "uid",
				TenantID:      "00000000-0000-4000-0000-000000000000",
				DeviceUID:     models.UID("2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c"),
				IPAddress:     "0.0.0.0",
				Authenticated: true,
			},
			expected: nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			session, err := mongostore.SessionCreate(context.TODO(), tc.session)
			assert.Equal(t, tc.expected, err)
			assert.NotEmpty(t, session)
		})
	}
}

func TestSessionUpdateDeviceUID(t *testing.T) {
	cases := []struct {
		description string
		oldUID      models.UID
		newUID      models.UID
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when device is not found",
			oldUID:      models.UID("nonexistent"),
			newUID:      models.UID("uid"),
			fixtures:    []string{fixtures.FixtureSessions},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when device is found",
			oldUID:      models.UID("2300230e3ca2f637636b4d025d2235269014865db5204b6d115386cbee89809c"),
			newUID:      models.UID("uid"),
			fixtures:    []string{fixtures.FixtureSessions},
			expected:    nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.SessionUpdateDeviceUID(context.TODO(), tc.oldUID, tc.newUID)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionSetAuthenticated(t *testing.T) {
	cases := []struct {
		description  string
		UID          models.UID
		authenticate bool
		fixtures     []string
		expected     error
	}{
		{
			description:  "fails when session is not found",
			UID:          models.UID("nonexistent"),
			authenticate: false,
			fixtures:     []string{fixtures.FixtureSessions},
			expected:     store.ErrNoDocuments,
		},
		{
			description:  "succeeds when session is found",
			UID:          models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
			authenticate: false,
			fixtures:     []string{fixtures.FixtureSessions},
			expected:     nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.SessionSetAuthenticated(context.TODO(), tc.UID, tc.authenticate)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionSetRecorded(t *testing.T) {
	cases := []struct {
		description  string
		UID          models.UID
		authenticate bool
		fixtures     []string
		expected     error
	}{
		{
			description:  "fails when session is not found",
			UID:          models.UID("nonexistent"),
			authenticate: false,
			fixtures:     []string{fixtures.FixtureSessions},
			expected:     store.ErrNoDocuments,
		},
		{
			description:  "succeeds when session is found",
			UID:          models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
			authenticate: false,
			fixtures:     []string{fixtures.FixtureSessions},
			expected:     nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.SessionSetAuthenticated(context.TODO(), tc.UID, tc.authenticate)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionSetLastSeen(t *testing.T) {
	cases := []struct {
		description string
		UID         models.UID
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when session is not found",
			UID:         models.UID("nonexistent"),
			fixtures:    []string{fixtures.FixtureSessions},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when session is found",
			UID:         models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
			fixtures:    []string{fixtures.FixtureSessions},
			expected:    nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.SessionSetLastSeen(context.TODO(), tc.UID)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionDeleteActives(t *testing.T) {
	cases := []struct {
		description string
		UID         models.UID
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when session is not found",
			UID:         models.UID("nonexistent"),
			fixtures:    []string{fixtures.FixtureSessions},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when session is found",
			UID:         models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
			fixtures:    []string{fixtures.FixtureSessions},
			expected:    nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.SessionDeleteActives(context.TODO(), tc.UID)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionGetRecordFrame(t *testing.T) {
	type Expected struct {
		r     []models.RecordedSession
		count int
		err   error
	}

	cases := []struct {
		description string
		UID         models.UID
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds",
			UID:         models.UID("e7f3a56d8b9e1dc4c285c98c8ea9c33032a17bda5b6c6b05a6213c2a02f97824"),
			fixtures:    []string{fixtures.FixtureSessions, fixtures.FixtureRecordedSessions},
			expected: Expected{
				r: []models.RecordedSession{
					{
						Time:     time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						UID:      "e7f3a56d8b9e1dc4c285c98c8ea9c33032a17bda5b6c6b05a6213c2a02f97824",
						Message:  "message",
						TenantID: "00000000-0000-4000-0000-000000000000",
						Width:    0,
						Height:   0,
					},
				},
				count: 1,
				err:   nil,
			},
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("recorded_sessions")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			doc := bson.M{
				"time":      tc.expected.r[0].Time,
				"uid":       tc.expected.r[0].UID,
				"message":   tc.expected.r[0].Message,
				"tenant_id": tc.expected.r[0].TenantID,
				"width":     tc.expected.r[0].Width,
				"height":    tc.expected.r[0].Height,
			}
			testData = append(testData, doc)

			if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			r, count, err := mongostore.SessionGetRecordFrame(context.TODO(), tc.UID)
			assert.Equal(t, tc.expected, Expected{r: r, count: count, err: err})
		})
	}
}

func TestSessionCreateRecordFrame(t *testing.T) {
	cases := []struct {
		description string
		UID         models.UID
		record      *models.RecordedSession
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when session is not found",
			UID:         models.UID(""),
			record: &models.RecordedSession{
				UID:      models.UID(""),
				Message:  "message",
				TenantID: "00000000-0000-4000-0000-000000000000",
				Time:     time.Now(),
				Width:    0,
				Height:   0,
			},
			fixtures: []string{fixtures.FixtureSessions},
			expected: store.ErrNoDocuments,
		},
		{
			description: "succeeds when session is found",
			UID:         models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
			record: &models.RecordedSession{
				UID:      models.UID("a3b0431f5df6a7827945d2e34872a5c781452bc36de42f8b1297fd9ecb012f68"),
				Message:  "message",
				TenantID: "00000000-0000-4000-0000-000000000000",
				Time:     time.Now(),
				Width:    0,
				Height:   0,
			},
			fixtures: []string{fixtures.FixtureSessions},
			expected: nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	// collection := mongostore.db.Collection("recorded_sessions")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck
			collection := mongostore.db.Collection("recorded_sessions")
			if tc.UID != "" {
				doc := bson.M{
					"uid":       tc.record.UID,
					"message":   tc.record.Message,
					"tenant_id": tc.record.TenantID,
					"time":      tc.record.Time,
					"width":     tc.record.Width,
					"height":    1,
				}

				if err := dbtest.InsertMockData(ctx, collection, []interface{}{doc}); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}

				collection := mongostore.db.Collection("sessions")
				mockData := bson.M{
					"uid":      tc.record.UID,
					"recorded": false,
				}
				if err := dbtest.InsertMockData(ctx, collection, []interface{}{mockData}); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}

			}

			err := mongostore.SessionCreateRecordFrame(context.TODO(), tc.UID, tc.record)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionDeleteRecordFrame(t *testing.T) {
	cases := []struct {
		description string
		UID         models.UID
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when record frame is not found",
			UID:         models.UID(""),
			fixtures:    []string{fixtures.FixtureSessions, fixtures.FixtureRecordedSessions},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when record frame is found",
			UID:         models.UID("e7f3a56d8b9e1dc4c285c98c8ea9c33032a17bda5b6c6b05a6213c2a02f97824"),
			fixtures:    []string{fixtures.FixtureSessions, fixtures.FixtureRecordedSessions},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")
	collection := mongostore.db.Collection("recorded_sessions")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			if tc.UID != "" {
				doc := bson.M{"uid": tc.UID}
				if err := dbtest.InsertMockData(ctx, collection, []interface{}{doc}); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.SessionDeleteRecordFrame(context.TODO(), tc.UID)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSessionDeleteRecordFrameByDate(t *testing.T) {
	type Expected struct {
		deletedCount int64
		updatedCount int64
		err          error
	}

	cases := []struct {
		description string
		lte         time.Time
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when there are no sessions to update or delete",
			lte:         time.Date(2023, time.January, 30, 12, 00, 0, 0, time.UTC),
			fixtures:    []string{},
			expected: Expected{
				deletedCount: 0,
				updatedCount: 0,
				err:          nil,
			},
		},
		{
			description: "succeeds to delete and update recorded sessions before specified date",
			lte:         time.Date(2023, time.January, 30, 12, 00, 0, 0, time.UTC),
			fixtures: []string{
				fixtures.FixtureSessions,
				fixtures.FixtureRecordedSessions,
			},
			expected: Expected{
				deletedCount: 2,
				updatedCount: 2,
				err:          nil,
			},
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")
	collection := mongostore.db.Collection("recorded_sessions")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint:errcheck

			doc := bson.M{"lte": tc.lte}
			if err := dbtest.InsertMockData(ctx, collection, []interface{}{doc}); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			deletedCount, updatedCount, err := mongostore.SessionDeleteRecordFrameByDate(ctx, tc.lte)
			assert.Equal(t, tc.expected, Expected{deletedCount, updatedCount, err})
		})
	}
}
