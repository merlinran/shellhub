package mongo

import (
	"context"
	"fmt"
	"sort"
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserList(t *testing.T) {
	type Expected struct {
		users []models.User
		count int
		err   error
	}

	cases := []struct {
		description string
		page        query.Paginator
		filters     query.Filters
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when users are found",
			page:        query.Paginator{Page: -1, PerPage: -1},
			filters:     query.Filters{},
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				users: []models.User{
					{
						ID:             "507f1f77bcf86cd799439011",
						CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						EmailMarketing: true,
						Confirmed:      true,
						UserData: models.UserData{
							Name:     "john doe",
							Username: "john_doe",
							Email:    "john.doe@test.com",
						},
						MaxNamespaces: 0,
						UserPassword: models.UserPassword{
							HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
						},
					},
					{
						ID:             "608f32a2c7351f001f6475e0",
						CreatedAt:      time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						LastLogin:      time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						EmailMarketing: true,
						Confirmed:      true,
						UserData: models.UserData{
							Name:     "Jane Smith",
							Username: "jane_smith",
							Email:    "jane.smith@test.com",
						},
						MaxNamespaces: 3,
						UserPassword: models.UserPassword{
							HashedPassword: "a0b8c29f4c8d57e542f5e81d35ebe801fd27f569f116fe670e8962d798512a1d",
						},
					},
					{
						ID:             "709f45b5e812c1002f3a67e7",
						CreatedAt:      time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						LastLogin:      time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						EmailMarketing: true,
						Confirmed:      true,
						UserData: models.UserData{
							Name:     "Bob Johnson",
							Username: "bob_johnson",
							Email:    "bob.johnson@test.com",
						},
						MaxNamespaces: 10,
						UserPassword: models.UserPassword{
							HashedPassword: "5f3b3956a1a150b73e6b27e674f27d7aeb01ab1a40c179c3e1aa6026a36655a2",
						},
					},
					{
						ID:             "80fdcea1d7299c002f3a67e8",
						CreatedAt:      time.Date(2023, 1, 4, 12, 0, 0, 0, time.UTC),
						EmailMarketing: false,
						Confirmed:      false,
						UserData: models.UserData{
							Name:     "Alex Rodriguez",
							Username: "alex_rodriguez",
							Email:    "alex.rodriguez@test.com",
						},
						MaxNamespaces: 3,
						UserPassword: models.UserPassword{
							HashedPassword: "c5093eb98678c7a3324825b84c6b67c1127b93786482ddbbd356e67e29b2763f",
						},
					},
				},
				count: 4,
				err:   nil,
			},
		},
		{
			description: "succeeds with filters",
			page:        query.Paginator{Page: -1, PerPage: -1},
			filters: query.Filters{
				Data: []query.Filter{
					{
						Type: "property",
						Params: &query.FilterProperty{
							Name:     "max_namespaces",
							Operator: "gt",
							Value:    "3",
						},
					},
				},
			},
			fixtures: []string{fixtures.FixtureUsers},
			expected: Expected{
				users: []models.User{
					{
						ID:             "709f45b5e812c1002f3a67e7",
						CreatedAt:      time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						LastLogin:      time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
						EmailMarketing: true,
						Confirmed:      true,
						UserData: models.UserData{
							Name:     "Bob Johnson",
							Username: "bob_johnson",
							Email:    "bob.johnson@test.com",
						},
						MaxNamespaces: 10,
						UserPassword: models.UserPassword{
							HashedPassword: "5f3b3956a1a150b73e6b27e674f27d7aeb01ab1a40c179c3e1aa6026a36655a2",
						},
					},
				},
				count: 1,
				err:   nil,
			},
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	// Due to the non-deterministic order of applying fixtures when dealing with multiple datasets,
	// we ensure that both the expected and result arrays are correctly sorted.
	sort := func(users []models.User) {
		sort.Slice(users, func(i, j int) bool {
			return users[i].ID < users[j].ID
		})
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			users, count, err := mongostore.UserList(context.TODO(), tc.page, tc.filters)
			sort(tc.expected.users)
			sort(users)
			assert.Equal(t, tc.expected, Expected{users: users, count: count, err: err})
		})
	}
}

func TestUserCreate(t *testing.T) {
	cases := []struct {
		description string
		user        *models.User
		fixtures    []string
		expected    error
	}{
		{
			description: "succeeds when data is valid",
			user: &models.User{
				ID: "507f1f77bcf86cd799439011",
				UserData: models.UserData{
					Name:     "john doe",
					Username: "john_doe",
					Email:    "john.doe@test.com",
				},
				UserPassword: models.UserPassword{
					HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
				},
			},
			fixtures: []string{},
			expected: nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			doc := bson.M{
				"id":       tc.user.ID,
				"name":     tc.user.Name,
				"email":    tc.user.Email,
				"username": tc.user.Username,
			}
			testData = append(testData, doc)

			if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			err := mongostore.UserCreate(context.TODO(), tc.user)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserGetByUsername(t *testing.T) {
	type Expected struct {
		user *models.User
		err  error
	}

	cases := []struct {
		description string
		username    string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when user is not found",
			username:    "",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				user: nil,
				err:  store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when user is found",
			username:    "john_doe",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				user: &models.User{
					ID:             "507f1f77bcf86cd799439011",
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EmailMarketing: true,
					Confirmed:      true,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "john.doe@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
					},
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

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.username != "" {
				doc := &models.User{
					ID:             tc.expected.user.ID,
					CreatedAt:      tc.expected.user.CreatedAt,
					LastLogin:      tc.expected.user.LastLogin,
					EmailMarketing: tc.expected.user.EmailMarketing,
					Confirmed:      tc.expected.user.Confirmed,
					UserData: models.UserData{
						Name:     tc.expected.user.UserData.Name,
						Username: tc.expected.user.UserData.Username,
						Email:    tc.expected.user.UserData.Email,
					},
					MaxNamespaces: tc.expected.user.MaxNamespaces,
					UserPassword: models.UserPassword{
						HashedPassword: tc.expected.user.UserPassword.HashedPassword,
					},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			user, err := mongostore.UserGetByUsername(context.TODO(), tc.username)
			assert.Equal(t, tc.expected, Expected{user: user, err: err})
		})
	}
}

func TestUserGetByEmail(t *testing.T) {
	type Expected struct {
		user *models.User
		err  error
	}

	cases := []struct {
		description string
		email       string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when email is not found",
			email:       "",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				user: nil,
				err:  store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when email is found",
			email:       "john.doe@test.com",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				user: &models.User{
					ID:             "507f1f77bcf86cd799439011",
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EmailMarketing: true,
					Confirmed:      true,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "john.doe@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
					},
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

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.email != "" {
				doc := &models.User{
					ID:             tc.expected.user.ID,
					CreatedAt:      tc.expected.user.CreatedAt,
					LastLogin:      tc.expected.user.LastLogin,
					EmailMarketing: tc.expected.user.EmailMarketing,
					Confirmed:      tc.expected.user.Confirmed,
					UserData: models.UserData{
						Name:     tc.expected.user.UserData.Name,
						Username: tc.expected.user.UserData.Username,
						Email:    tc.expected.user.UserData.Email,
					},
					MaxNamespaces: tc.expected.user.MaxNamespaces,
					UserPassword: models.UserPassword{
						HashedPassword: tc.expected.user.UserPassword.HashedPassword,
					},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			user, err := mongostore.UserGetByEmail(context.TODO(), tc.email)
			assert.Equal(t, tc.expected, Expected{user: user, err: err})
		})
	}
}
func TestUserGetByID(t *testing.T) {
	type Expected struct {
		user *models.User
		ns   int
		err  error
	}

	cases := []struct {
		description string
		id          string
		ns          bool
		fixtures    []string
		expected    Expected
	}{
		// {
		// 	description: "fails when user is not found",
		// 	id:          "000000000000000000000000",
		// 	fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureNamespaces},
		// 	expected: Expected{
		// 		user: nil,
		// 		ns:   0,
		// 		err:  store.ErrNoDocuments,
		// 	},
		// },
		{
			description: "succeeds when user is found with ns equal false",
			id:          "507f1f77bcf86cd799439011",
			ns:          false,
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureNamespaces},
			expected: Expected{
				user: &models.User{
					ID:             "507f1f77bcf86cd799439011",
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EmailMarketing: true,
					Confirmed:      true,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "john.doe@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
					},
				},
				ns:  0,
				err: nil,
			},
		},
		{
			description: "succeeds when user is found with ns equal true",
			id:          "507f1f77bcf86cd799439012",
			ns:          true,
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureNamespaces},
			expected: Expected{
				user: &models.User{
					ID:             "507f1f77bcf86cd799439012",
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EmailMarketing: true,
					Confirmed:      true,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "john.doe@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
					},
				},
				ns:  1,
				err: nil,
			},
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collectionUsers := mongostore.db.Collection("users")
	collectionNamespaces := mongostore.db.Collection("namespaces")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "000000000000000000000000" {
				doc := &models.User{
					ID:             tc.expected.user.ID,
					CreatedAt:      tc.expected.user.CreatedAt,
					LastLogin:      tc.expected.user.LastLogin,
					EmailMarketing: tc.expected.user.EmailMarketing,
					Confirmed:      tc.expected.user.Confirmed,
					UserData: models.UserData{
						Name:     tc.expected.user.UserData.Name,
						Username: tc.expected.user.UserData.Username,
						Email:    tc.expected.user.UserData.Email,
					},
					MaxNamespaces: tc.expected.user.MaxNamespaces,
					UserPassword: models.UserPassword{
						HashedPassword: tc.expected.user.UserPassword.HashedPassword,
					},
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collectionUsers, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}

				if tc.expected.ns > 0 {
					var testDataNamespaces []interface{}
					for i := 0; i < tc.expected.ns; i++ {
						docNamespace := bson.M{
							"created_at":  time.Now(),
							"name":        fmt.Sprintf("namespace-%d", i+1),
							"owner":       tc.expected.user.ID,
							"tenant_id":   "00000000-0000-4000-0000-000000000000",
							"description": fmt.Sprintf("Description of namespace-%d", i+1),
						}
						testDataNamespaces = append(testDataNamespaces, docNamespace)
					}

					if err := dbtest.InsertMockData(ctx, collectionNamespaces, testDataNamespaces); err != nil {
						t.Fatalf("failed to insert namespace documents: %v", err)

					}
				}
			}

			user, ns, err := mongostore.UserGetByID(context.TODO(), tc.id, tc.ns)
			assert.Equal(t, tc.expected, Expected{user: user, ns: ns, err: err})
		})
	}
}

func TestUserUpdateData(t *testing.T) {
	cases := []struct {
		description string
		id          string
		data        models.User
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user is not found",
			id:          "000000000000000000000000",
			data: models.User{
				LastLogin: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UserData: models.UserData{
					Name:     "edited name",
					Username: "edited_name",
					Email:    "edited@test.com",
				},
			},
			fixtures: []string{fixtures.FixtureUsers},
			expected: store.ErrNoDocuments,
		},
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			fixtures:    []string{fixtures.FixtureUsers},
			data: models.User{
				LastLogin: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UserData: models.UserData{
					Name:     "edited name",
					Username: "edited_name",
					Email:    "edited@test.com",
				},
			},
			expected: nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "000000000000000000000000" {
				objID, err := primitive.ObjectIDFromHex(tc.id)
				if err != nil {
					t.Fatalf("failed to convert ID to ObjectID: %v", err)
				}

				doc := bson.M{
					"_id":        objID,
					"last_login": tc.data.LastLogin,
					"name":       "old",
					"username":   tc.data.UserData.Username,
					"email":      tc.data.UserData.Email,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.UserUpdateData(context.TODO(), tc.id, tc.data)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserUpdatePassword(t *testing.T) {
	cases := []struct {
		description string
		id          string
		password    string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user id is not valid",
			id:          "",
			password:    "other_password",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    store.ErrInvalidHex,
		},
		{
			description: "fails when user is not found",
			id:          "000000000000000000000000",
			password:    "other_password",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			password:    "other_password",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "000000000000000000000000" && tc.id != "" {
				objID, err := primitive.ObjectIDFromHex(tc.id)
				if err != nil {
					t.Fatalf("failed to convert ID to ObjectID: %v", err)
				}
				doc := bson.M{
					"_id":      objID,
					"password": "old",
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.UserUpdatePassword(context.TODO(), tc.password, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserUpdateAccountStatus(t *testing.T) {
	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user id is not valid",
			id:          "",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    store.ErrInvalidHex,
		},
		{
			description: "fails when user is not found",
			id:          "000000000000000000000000",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when user is found",
			id:          "80fdcea1d7299c002f3a67e8",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "000000000000000000000000" && tc.id != "" {
				objID, err := primitive.ObjectIDFromHex(tc.id)
				if err != nil {
					t.Fatalf("failed to convert ID to ObjectID: %v", err)
				}
				doc := bson.M{
					"_id":       objID,
					"confirmed": false,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.UserUpdateAccountStatus(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserUpdateFromAdmin(t *testing.T) {
	cases := []struct {
		description string
		id          string
		name        string
		username    string
		email       string
		password    string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user is not found",
			id:          "000000000000000000000000",
			name:        "other name",
			username:    "other_name",
			email:       "other.email@test.com",
			password:    "other_password",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			name:        "other name",
			username:    "other_name",
			email:       "other.email@test.com",
			password:    "other_password",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "000000000000000000000000" {
				objID, err := primitive.ObjectIDFromHex(tc.id)
				if err != nil {
					t.Fatalf("failed to convert ID to ObjectID: %v", err)
				}
				doc := bson.M{
					"_id":      objID,
					"name":     "old",
					"username": tc.username,
					"email":    tc.email,
					"password": tc.password,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.UserUpdateFromAdmin(context.TODO(), tc.name, tc.username, tc.email, tc.password, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserCreateToken(t *testing.T) {
	cases := []struct {
		description string
		token       *models.UserTokenRecover
		fixtures    []string
		expected    error
	}{
		{
			description: "succeeds when data is valid",
			token: &models.UserTokenRecover{
				Token: "token",
				User:  "507f1f77bcf86cd799439011",
			},
			fixtures: []string{},
			expected: nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("recovery_tokens")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}

			doc := &models.UserTokenRecover{
				Token: tc.token.Token,
			}
			testData = append(testData, doc)

			if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
				t.Fatalf("failed to insert documents: %v", err)
			}

			err := mongostore.UserCreateToken(context.TODO(), tc.token)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserTokenGet(t *testing.T) {
	type Expected struct {
		token *models.UserTokenRecover
		err   error
	}

	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when user is not found",
			id:          "",
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureRecoveryTokens},
			expected: Expected{
				token: nil,
				err:   store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureRecoveryTokens},
			expected: Expected{
				token: &models.UserTokenRecover{
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					User:      "507f1f77bcf86cd799439011",
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

	collection := mongostore.db.Collection("recovery_tokens")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "" {
				doc := bson.M{
					"uid":        tc.expected.token.Token,
					"created_at": tc.expected.token.CreatedAt,
					"user":       tc.id,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			token, err := mongostore.UserGetToken(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, Expected{token: token, err: err})
		})
	}
}

func TestUserDeleteTokens(t *testing.T) {
	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    error
	}{
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureRecoveryTokens},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("recovery_tokens")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "" {
				doc := bson.M{
					"user": tc.id,
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}

			err := mongostore.UserDeleteTokens(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserDelete(t *testing.T) {
	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user is not found",
			id:          "000000000000000000000000",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			fixtures:    []string{fixtures.FixtureUsers},
			expected:    nil,
		},
	}

	ctx := context.TODO()

	client, host, stopContainer := dbtest.StartTestContainer(ctx)
	defer stopContainer()

	mongostore := NewStore(client.Database("test"), cache.NewNullCache())
	fixtures.Init(host, "test")

	collection := mongostore.db.Collection("users")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			var testData []interface{}
			if tc.id != "000000000000000000000000" {
				objID, err := primitive.ObjectIDFromHex(tc.id)
				if err != nil {
					t.Fatalf("failed to convert ID to ObjectID: %v", err)
				}
				doc := bson.M{
					"_id":  objID,
					"name": "name",
				}
				testData = append(testData, doc)

				if err := dbtest.InsertMockData(ctx, collection, testData); err != nil {
					t.Fatalf("failed to insert documents: %v", err)
				}
			}
			err := mongostore.UserDelete(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserDetachInfo(t *testing.T) {
	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	user := models.User{
		ID: "60af83d418d2dc3007cd445c",
		UserData: models.UserData{
			Name:     "name",
			Username: "username",
			Email:    "user@email.com",
		},
		UserPassword: models.UserPassword{
			HashedPassword: "fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4",
		},
	}

	objID, err := primitive.ObjectIDFromHex(user.ID)

	assert.NoError(t, err)

	_, _ = db.Client().Database("test").Collection("users").InsertOne(context.TODO(), bson.M{
		"_id":      objID,
		"name":     user.Name,
		"username": user.Username,
		"password": user.HashedPassword,
		"email":    user.Email,
	})

	namespacesOwner := []*models.Namespace{
		{
			Owner: user.ID,
			Name:  "ns2",
			Members: []models.Member{
				{
					ID:   user.ID,
					Role: guard.RoleOwner,
				},
			},
		},
		{
			Owner: user.ID,
			Name:  "ns4",
			Members: []models.Member{
				{
					ID:   user.ID,
					Role: guard.RoleOwner,
				},
			},
		},
	}

	namespacesMember := []*models.Namespace{
		{
			Owner: "id2",
			Name:  "ns1",
			Members: []models.Member{
				{
					ID:   user.ID,
					Role: guard.RoleObserver,
				},
			},
		},
		{
			Owner: "id2",
			Name:  "ns3",
			Members: []models.Member{
				{
					ID:   user.ID,
					Role: guard.RoleObserver,
				},
			},
		},
		{
			Owner: "id2",
			Name:  "ns5",
			Members: []models.Member{
				{
					ID:   user.ID,
					Role: guard.RoleObserver,
				},
			},
		},
	}

	for _, n := range namespacesOwner {
		inserted, err := db.Client().Database("test").Collection("namespaces").InsertOne(context.TODO(), n)
		t.Log(inserted.InsertedID)
		assert.NoError(t, err)
	}

	for _, n := range namespacesMember {
		inserted, err := db.Client().Database("test").Collection("namespaces").InsertOne(context.TODO(), n)
		t.Log(inserted.InsertedID)
		assert.NoError(t, err)
	}

	u, err := mongostore.UserGetByUsername(context.TODO(), "username")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, u.Username)

	namespacesMap, err := mongostore.UserDetachInfo(context.TODO(), user.ID)

	assert.NoError(t, err)
	assert.Equal(t, namespacesMap["owner"], namespacesOwner)
	assert.Equal(t, namespacesMap["member"], namespacesMember)
}
