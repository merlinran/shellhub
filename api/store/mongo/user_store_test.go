package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/pkg/fixtures"
	"github.com/shellhub-io/shellhub/api/pkg/guard"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserList(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	type Expected struct {
		users []models.User
		count int
		err   error
	}

	cases := []struct {
		description string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when users are found",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				users: []models.User{
					{
						ID:             "507f1f77bcf86cd799439011",
						CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						EmailMarketing: false,
						Confirmed:      false,
						UserData: models.UserData{
							Name:     "john doe",
							Username: "john_doe",
							Email:    "user@test.com",
						},
						MaxNamespaces: 0,
						UserPassword: models.UserPassword{
							Password: "secret123",
						},
					},
				},
				count: 1,
				err:   nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			users, count, err := mongostore.UserList(ctx, paginator.Query{Page: -1, PerPage: -1}, nil)
			assert.Equal(t, tc.expected, Expected{users: users, count: count, err: err})
		})
	}
}

func TestUserListWithFilter(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	type Expected struct {
		users []models.User
		count int
		err   error
	}

	cases := []struct {
		description string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when no users are found",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				users: []models.User{},
				count: 0,
				err:   nil,
			},
		},
	}

	filters := []models.Filter{
		{
			Type:   "property",
			Params: &models.PropertyParams{Name: "namespaces", Operator: "gt", Value: "1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			users, count, err := mongostore.UserList(ctx, paginator.Query{Page: -1, PerPage: -1}, filters)
			assert.Equal(t, tc.expected, Expected{users: users, count: count, err: err})
		})
	}
}

func TestUserCreate(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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
					Email:    "user@test.com",
				},
				UserPassword: models.UserPassword{
					Password: "secret123",
				},
			},
			fixtures: []string{},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserCreate(ctx, tc.user)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserGetByUsername(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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
			username:    "nonexistent",
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
					EmailMarketing: false,
					Confirmed:      false,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "user@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						Password: "secret123",
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			user, err := mongostore.UserGetByUsername(ctx, tc.username)
			assert.Equal(t, tc.expected, Expected{user: user, err: err})
		})
	}
}

func TestUserGetByEmail(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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
			email:       "nonexistent",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				user: nil,
				err:  store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when email is found",
			email:       "user@test.com",
			fixtures:    []string{fixtures.FixtureUsers},
			expected: Expected{
				user: &models.User{
					ID:             "507f1f77bcf86cd799439011",
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EmailMarketing: false,
					Confirmed:      false,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "user@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						Password: "secret123",
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			user, err := mongostore.UserGetByEmail(ctx, tc.email)
			assert.Equal(t, tc.expected, Expected{user: user, err: err})
		})
	}
}

func TestUserGetByID(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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
		{
			description: "fails when user is not found",
			id:          "507f1f77bcf86cd7994390bb",
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureNamespaces},
			expected: Expected{
				user: nil,
				ns:   0,
				err:  store.ErrNoDocuments,
			},
		},
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
					EmailMarketing: false,
					Confirmed:      false,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "user@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						Password: "secret123",
					},
				},
				ns:  0,
				err: nil,
			},
		},
		{
			description: "succeeds when user is found with ns equal true",
			id:          "507f1f77bcf86cd799439011",
			ns:          true,
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureNamespaces},
			expected: Expected{
				user: &models.User{
					ID:             "507f1f77bcf86cd799439011",
					CreatedAt:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:      time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EmailMarketing: false,
					Confirmed:      false,
					UserData: models.UserData{
						Name:     "john doe",
						Username: "john_doe",
						Email:    "user@test.com",
					},
					MaxNamespaces: 0,
					UserPassword: models.UserPassword{
						Password: "secret123",
					},
				},
				ns:  1,
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			user, ns, err := mongostore.UserGetByID(ctx, tc.id, tc.ns)
			assert.Equal(t, tc.expected, Expected{user: user, ns: ns, err: err})
		})
	}
}

func TestUserUpdateData(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserUpdateData(ctx, tc.id, tc.data)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserUpdatePassword(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	cases := []struct {
		description string
		id          string
		password    string
		fixtures    []string
		expected    error
	}{
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

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserUpdatePassword(ctx, tc.password, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserUpdateAccountStatus(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserUpdateAccountStatus(ctx, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserUpdateFromAdmin(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserUpdateFromAdmin(ctx, tc.name, tc.username, tc.email, tc.password, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserCreateToken(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserCreateToken(ctx, tc.token)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserTokenGet(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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
			id:          "000000000000000000000000",
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
					Token:     "token",
					User:      "507f1f77bcf86cd799439011",
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			token, err := mongostore.UserGetToken(ctx, tc.id)
			assert.Equal(t, tc.expected, Expected{token: token, err: err})
		})
	}
}

func TestUserDeleteTokens(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when user is not found",
			id:          "000000000000000000000000",
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureRecoveryTokens},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when user is found",
			id:          "507f1f77bcf86cd799439011",
			fixtures:    []string{fixtures.FixtureUsers, fixtures.FixtureRecoveryTokens},
			expected:    nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserDeleteTokens(ctx, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserDelete(t *testing.T) {
	ctx := context.TODO()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

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

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.UserDelete(ctx, tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestUserDetachInfo(t *testing.T) {
	db := dbtest.DBServer{}
	defer db.Stop()

	ctx := context.TODO()
	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	user := models.User{UserData: models.UserData{Name: "name", Username: "username", Email: "user@email.com"}, UserPassword: models.UserPassword{Password: "password"}, ID: "60af83d418d2dc3007cd445c"}

	objID, err := primitive.ObjectIDFromHex(user.ID)

	assert.NoError(t, err)

	_, _ = db.Client().Database("test").Collection("users").InsertOne(ctx, bson.M{
		"_id":      objID,
		"name":     user.Name,
		"username": user.Username,
		"password": user.Password,
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
		inserted, err := db.Client().Database("test").Collection("namespaces").InsertOne(ctx, n)
		t.Log(inserted.InsertedID)
		assert.NoError(t, err)
	}

	for _, n := range namespacesMember {
		inserted, err := db.Client().Database("test").Collection("namespaces").InsertOne(ctx, n)
		t.Log(inserted.InsertedID)
		assert.NoError(t, err)
	}

	u, err := mongostore.UserGetByUsername(ctx, "username")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, u.Username)

	namespacesMap, err := mongostore.UserDetachInfo(ctx, user.ID)

	assert.NoError(t, err)
	assert.Equal(t, namespacesMap["owner"], namespacesOwner)
	assert.Equal(t, namespacesMap["member"], namespacesMember)
}
