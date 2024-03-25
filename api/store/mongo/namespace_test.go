package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/pkg/fixtures"
	"github.com/shellhub-io/shellhub/api/pkg/guard"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/query"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNamespaceList(t *testing.T) {
	type Expected struct {
		ns    []models.Namespace
		count int
		err   error
	}

	cases := []struct {
		description string
		page        query.Paginator
		filters     query.Filters
		export      bool
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when namespaces list is not empty",
			page:        query.Paginator{Page: -1, PerPage: -1},
			filters:     query.Filters{},
			export:      false,
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns: []models.Namespace{
					{
						CreatedAt:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						Name:         "namespace-1",
						Owner:        "507f1f77bcf86cd799439011",
						TenantID:     "00000000-0000-4000-0000-000000000000",
						MaxDevices:   -1,
						Settings:     &models.NamespaceSettings{SessionRecord: true},
						DevicesCount: 0,
					},
					{
						CreatedAt:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						Name:         "namespace-2",
						Owner:        "6509e169ae6144b2f56bf288",
						TenantID:     "00000000-0000-4001-0000-000000000000",
						MaxDevices:   10,
						Settings:     &models.NamespaceSettings{SessionRecord: false},
						DevicesCount: 0,
					},
					{
						CreatedAt:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						Name:         "namespace-3",
						Owner:        "657b0e3bff780d625f74e49a",
						TenantID:     "00000000-0000-4002-0000-000000000000",
						MaxDevices:   3,
						Settings:     &models.NamespaceSettings{SessionRecord: true},
						DevicesCount: 0,
					},
					{
						CreatedAt:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						Name:         "namespace-4",
						Owner:        "6577267d8752d05270a4c07d",
						TenantID:     "00000000-0000-4003-0000-000000000000",
						MaxDevices:   -1,
						Settings:     &models.NamespaceSettings{SessionRecord: true},
						DevicesCount: 0,
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

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			for _, ns := range tc.expected.ns {
				doc := bson.M{
					"created_at":    ns.CreatedAt,
					"name":          ns.Name,
					"owner":         ns.Owner,
					"tenant_id":     ns.TenantID,
					"max_devices":   ns.MaxDevices,
					"settings":      ns.Settings,
					"devices_count": ns.DevicesCount,
				}
				testData = append(testData, doc)
			}

			if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			ns, count, err := mongostore.NamespaceList(context.TODO(), tc.page, tc.filters, tc.export)
			assert.Equal(t, tc.expected, Expected{ns: ns, count: count, err: err})
		})
	}
}

func TestNamespaceGet(t *testing.T) {
	type Expected struct {
		ns  *models.Namespace
		err error
	}

	cases := []struct {
		description string
		tenant      string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when tenant is not found",
			tenant:      "",
			fixtures:    []string{fixtures.FixtureNamespaces, fixtures.FixtureDevices},
			expected: Expected{
				ns:  nil,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when tenant is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			fixtures:    []string{fixtures.FixtureNamespaces, fixtures.FixtureDevices},
			expected: Expected{
				ns: &models.Namespace{
					Name:         "namespace-1",
					Owner:        "507f1f77bcf86cd799439011",
					TenantID:     "00000000-0000-4000-0000-000000000000",
					MaxDevices:   -1,
					Settings:     &models.NamespaceSettings{SessionRecord: true},
					DevicesCount: 0,
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

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.tenant != "" {
				doc := bson.M{
					"tenant_id":     tc.expected.ns.TenantID,
					"name":          tc.expected.ns.Name,
					"owner":         tc.expected.ns.Owner,
					"max_devices":   tc.expected.ns.MaxDevices,
					"settings":      tc.expected.ns.Settings,
					"devices_count": tc.expected.ns.DevicesCount,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			ns, err := mongostore.NamespaceGet(context.TODO(), tc.tenant)
			assert.Equal(t, tc.expected, Expected{ns: ns, err: err})
		})
	}
}

func TestNamespaceGetByName(t *testing.T) {
	type Expected struct {
		ns  *models.Namespace
		err error
	}

	cases := []struct {
		description string
		name        string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when namespace is not found",
			name:        "",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns:  nil,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when namespace is found",
			name:        "namespace-1",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns: &models.Namespace{
					Name:         "namespace-1",
					Owner:        "507f1f77bcf86cd799439011",
					TenantID:     "00000000-0000-4000-0000-000000000000",
					MaxDevices:   -1,
					Settings:     &models.NamespaceSettings{SessionRecord: true},
					DevicesCount: 0,
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

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.name != "" {
				doc := bson.M{
					"name":          tc.name,
					"owner":         tc.expected.ns.Owner,
					"tenant_id":     tc.expected.ns.TenantID,
					"max_devices":   tc.expected.ns.MaxDevices,
					"settings":      tc.expected.ns.Settings,
					"devices_count": tc.expected.ns.DevicesCount,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			ns, err := mongostore.NamespaceGetByName(context.TODO(), tc.name)
			assert.Equal(t, tc.expected, Expected{ns: ns, err: err})
		})
	}
}

func TestNamespaceGetFirst(t *testing.T) {
	type Expected struct {
		ns  *models.Namespace
		err error
	}

	cases := []struct {
		description string
		member      string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when member is not found",
			member:      "",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns:  nil,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when member is found",
			member:      "507f1f77bcf86cd799439011",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns: &models.Namespace{
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "namespace-1",
					Owner:     "507f1f77bcf86cd799439011",
					TenantID:  "00000000-0000-4000-0000-000000000000",
					Members: []models.Member{
						{
							ID:   "507f1f77bcf86cd799439011",
							Role: guard.RoleOwner,
						},
						{
							ID:   "6509e169ae6144b2f56bf288",
							Role: guard.RoleObserver,
						},
					},
					MaxDevices:   -1,
					Settings:     &models.NamespaceSettings{SessionRecord: true},
					DevicesCount: 0,
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

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.member != "" {
				ns := &models.Namespace{
					Name:         tc.expected.ns.Name,
					Owner:        tc.member,
					TenantID:     tc.expected.ns.TenantID,
					Members:      tc.expected.ns.Members,
					MaxDevices:   tc.expected.ns.MaxDevices,
					Settings:     tc.expected.ns.Settings,
					DevicesCount: tc.expected.ns.DevicesCount,
					CreatedAt:    tc.expected.ns.CreatedAt,
				}
				testData = append(testData, ns)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			ns, err := mongostore.NamespaceGetFirst(context.TODO(), tc.member)
			assert.Equal(t, tc.expected, Expected{ns: ns, err: err})
		})
	}
}

func TestNamespaceCreate(t *testing.T) {
	type Expected struct {
		ns  *models.Namespace
		err error
	}

	cases := []struct {
		description string
		ns          *models.Namespace
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when data is valid",
			ns: &models.Namespace{
				Name:     "namespace-1",
				Owner:    "507f1f77bcf86cd799439011",
				TenantID: "00000000-0000-4000-0000-000000000000",
				Members: []models.Member{
					{
						ID:   "507f1f77bcf86cd799439011",
						Role: guard.RoleOwner,
					},
				},
				MaxDevices: -1,
				Settings:   &models.NamespaceSettings{SessionRecord: true},
			},
			fixtures: []string{},
			expected: Expected{
				ns: &models.Namespace{
					Name:     "namespace-1",
					Owner:    "507f1f77bcf86cd799439011",
					TenantID: "00000000-0000-4000-0000-000000000000",
					Members: []models.Member{
						{
							ID:   "507f1f77bcf86cd799439011",
							Role: guard.RoleOwner,
						},
					},
					MaxDevices: -1,
					Settings:   &models.NamespaceSettings{SessionRecord: true},
				},
				err: nil,
			},
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

			ns, err := mongostore.NamespaceCreate(context.TODO(), tc.ns)
			assert.Equal(t, tc.expected, Expected{ns: ns, err: err})
		})
	}
}

func TestNamespaceEdit(t *testing.T) {
	cases := []struct {
		description string
		tenant      string
		changes     *models.NamespaceChanges
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when tenant is not found",
			tenant:      "",
			changes: &models.NamespaceChanges{
				Name: "edited-namespace",
			},
			fixtures: []string{fixtures.FixtureNamespaces},
			expected: store.ErrNoDocuments,
		},
		{
			description: "succeeds when tenant is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			changes: &models.NamespaceChanges{
				Name: "edited-namespace",
			},
			fixtures: []string{fixtures.FixtureNamespaces},
			expected: nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.tenant != "" {
				doc := bson.M{"tenant_id": tc.tenant, "name": "old"}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.NamespaceEdit(context.TODO(), tc.tenant, tc.changes)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestNamespaceUpdate(t *testing.T) {
	cases := []struct {
		description string
		tenant      string
		ns          *models.Namespace
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when tenant is not found",
			tenant:      "",
			ns: &models.Namespace{
				Name:       "edited-namespace",
				MaxDevices: 3,
				Settings:   &models.NamespaceSettings{SessionRecord: true},
			},
			fixtures: []string{fixtures.FixtureNamespaces},
			expected: store.ErrNoDocuments,
		},
		{
			description: "succeeds when tenant is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			ns: &models.Namespace{
				Name:       "edited-namespace",
				MaxDevices: 3,
				Settings:   &models.NamespaceSettings{SessionRecord: true},
			},
			fixtures: []string{fixtures.FixtureNamespaces},
			expected: nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.tenant != "" {
				doc := bson.M{
					"tenant_id":   tc.tenant,
					"name":        "old",
					"max_devices": 10,
					"settings":    bson.M{"session_record": false},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.NamespaceUpdate(context.TODO(), tc.tenant, tc.ns)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestNamespaceDelete(t *testing.T) {
	type Expected struct {
		err error
	}

	cases := []struct {
		description string
		tenant      string
		fixtures    []string
		expected    Expected
	}{
		// {
		// 	description: "fails when namespace is not found",
		// 	tenant:      "",
		// 	fixtures:    []string{fixtures.FixtureNamespaces},
		// 	expected: Expected{
		// 		err: store.ErrNoDocuments,
		// 	},
		// },
		{
			description: "succeeds when namespace is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				err: nil,
			},
		},
	}

	ctx := context.TODO()
	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.tenant != "" {
				mockNamespace := &models.Namespace{
					Owner:    "507f1f77bcf86cd799439011",
					TenantID: tc.tenant,
				}
				doc := bson.M{
					"tenant_id": mockNamespace.TenantID,
					"owner":     mockNamespace.Owner,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}

				for _, coll := range []string{"devices", "sessions", "connected_devices", "firewall_rules", "public_keys", "recorded_sessions"} {
					mockData := []interface{}{bson.M{"tenant_id": tc.tenant}}
					if err := dbtest.InsertMockData(ctx, mongostore.db.Collection(coll), mockData); err != nil {
						t.Fatalf("failed to insert documents: %v", err)
					}
				}

				mockUserUpdateData := bson.M{"_id": "507f1f77bcf86cd799439011", "namespaces": 2}
				if err := dbtest.InsertMockData(ctx, mongostore.db.Collection("users"), []interface{}{mockUserUpdateData}); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.NamespaceDelete(ctx, tc.tenant)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestNamespaceAddMember(t *testing.T) {
	type Expected struct {
		ns  *models.Namespace
		err error
	}

	cases := []struct {
		description string
		tenant      string
		member      string
		role        string
		fixtures    []string
		expected    Expected
	}{
		// {
		// 	description: "fails when tenant is not found",
		// 	tenant:      "",
		// 	member:      "6509de884238881ac1b2b289",
		// 	role:        guard.RoleObserver,
		// 	fixtures:    []string{fixtures.FixtureNamespaces},
		// 	expected: Expected{
		// 		ns:  nil,
		// 		err: store.ErrNoDocuments,
		// 	},
		// },
		// {
		// 	description: "fails when member has already been added",
		// 	tenant:      "00000000-0000-4000-0000-000000000000",
		// 	member:      "6509e169ae6144b2f56bf287",
		// 	role:        guard.RoleObserver,
		// 	fixtures:    []string{fixtures.FixtureNamespaces},
		// 	expected: Expected{
		// 		ns:  nil,
		// 		err: ErrNamespaceDuplicatedMember,
		// 	},
		// },
		{
			description: "succeeds when tenant is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			member:      "6509de884238881ac1b2b289",
			role:        guard.RoleObserver,
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns: &models.Namespace{
					TenantID: "00000000-0000-4000-0000-000000000000",
					Members: []models.Member{
						{
							ID:   "507f1f77bcf86cd79439011",
							Role: guard.RoleOwner,
						},
						{
							ID:   "6509e169ae6144b2f56bf288",
							Role: guard.RoleObserver,
						},
					},
					MaxDevices:   0,
					Settings:     &models.NamespaceSettings{SessionRecord: true},
					DevicesCount: 0,
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

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.tenant != "" && tc.member != "" {
				doc := bson.M{
					"tenant_id": tc.tenant,
					"members": []bson.M{
						{
							"id":   "6509e169ae6144b2f56bf287",
							"role": "old",
						},
					},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			ns, err := mongostore.NamespaceAddMember(context.TODO(), tc.tenant, tc.member, tc.role)
			// assert.Equal(t, tc.expected, Expected{err: err})
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, len(tc.expected.ns.Members), len(ns.Members))

		})
	}
}

func TestNamespaceEditMember(t *testing.T) {
	cases := []struct {
		description string
		tenant      string
		member      string
		role        string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user is not found",
			tenant:      "",
			member:      "000000000000000000000000",
			role:        guard.RoleObserver,
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected:    ErrUserNotFound,
		},
		{
			description: "succeeds when tenant and user is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			member:      "6509e169ae6144b2f56bf288",
			role:        guard.RoleOperator,
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown()

			var testData []interface{}
			if tc.tenant != "" && tc.member != "" {
				doc := bson.M{
					"tenant_id": tc.tenant,
					"members": []bson.M{
						{
							"id":   tc.member,
							"role": guard.RoleOwner,
						},
					},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.NamespaceEditMember(context.TODO(), tc.tenant, tc.member, tc.role)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestNamespaceRemoveMember(t *testing.T) {
	type Expected struct {
		ns  *models.Namespace
		err error
	}

	cases := []struct {
		description string
		tenant      string
		member      string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when tenant is not found",
			tenant:      "",
			member:      "6509de884238881ac1b2b289",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns:  nil,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "fails when member is not found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			member:      "",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns:  nil,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when tenant and user is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			member:      "6509e169ae6144b2f56bf288",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				ns: &models.Namespace{
					Name:     "namespace-1",
					Owner:    "507f1f77bcf86cd799439011",
					TenantID: "00000000-0000-4000-0000-000000000000",
					Settings: &models.NamespaceSettings{SessionRecord: true},
					Members:  []models.Member{},
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

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown()

			var testData []interface{}
			if tc.tenant != "" && tc.member != "" {
				doc := bson.M{
					"tenant_id": tc.tenant,
					"name":      tc.expected.ns.Name,
					"owner":     tc.expected.ns.Owner,
					"settings":  tc.expected.ns.Settings,
					"members": []bson.M{
						{
							"id":   tc.member,
							"role": guard.RoleOwner,
						},
					},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			ns, err := mongostore.NamespaceRemoveMember(context.TODO(), tc.tenant, tc.member)
			fmt.Println(ns)
			assert.Equal(t, tc.expected, Expected{ns: ns, err: err})
		})
	}
}

func TestNamespaceSetSessionRecord(t *testing.T) {
	cases := []struct {
		description string
		tenant      string
		sessionRec  bool
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when tenant is not found",
			tenant:      "",
			sessionRec:  true,
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when tenant is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			sessionRec:  true,
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown()

			var testData []interface{}
			if tc.tenant != "" {
				doc := bson.M{"tenant_id": tc.tenant, "settings": bson.M{"session_record": true}}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.NamespaceSetSessionRecord(context.TODO(), tc.sessionRec, tc.tenant)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestNamespaceGetSessionRecord(t *testing.T) {
	type Expected struct {
		set bool
		err error
	}

	cases := []struct {
		description string
		tenant      string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when tenant is not found",
			tenant:      "",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				set: false,
				err: store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when tenant is found",
			tenant:      "00000000-0000-4000-0000-000000000000",
			fixtures:    []string{fixtures.FixtureNamespaces},
			expected: Expected{
				set: true,
				err: nil,
			},
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown()

			var testData []interface{}
			if tc.tenant != "" {
				doc := bson.M{"tenant_id": tc.tenant, "settings": bson.M{"session_record": true}}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			set, err := mongostore.NamespaceGetSessionRecord(context.TODO(), tc.tenant)
			assert.Equal(t, tc.expected, Expected{set: set, err: err})
		})
	}
}
