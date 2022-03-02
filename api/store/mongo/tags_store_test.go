package mongo

import (
	"context"
	"sort"
	"testing"

	"github.com/shellhub-io/shellhub/api/cache"
	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	db := dbtest.DBServer{}
	defer db.Stop()

	ctx := context.TODO()
	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	device1 := models.Device{
		UID:       "1",
		Namespace: "namespace1",
		TenantID:  "tenant1",
		Tags: []string{
			"device1",
			"device5",
			"device3",
		},
	}

	device2 := models.Device{
		UID:       "2",
		Namespace: "namespace2",
		TenantID:  "tenant2",
		Tags: []string{
			"device4",
			"device5",
			"device6",
		},
	}

	_, err := db.Client().Database("test").Collection("devices").InsertOne(ctx, &device1)
	assert.NoError(t, err)

	_, err = db.Client().Database("test").Collection("devices").InsertOne(ctx, &device2)
	assert.NoError(t, err)

	tags, count, err := mongostore.TagsGet(ctx, "tenant1")
	assert.NoError(t, err)
	assert.Equal(t, count, 3)

	sort.Strings(tags) // Guarantee the order for comparison.
	assert.Equal(t, []string{"device1", "device3", "device5"}, tags)
}

func TestRenameTag(t *testing.T) {
	db := dbtest.DBServer{}
	defer db.Stop()

	ctx := context.TODO()
	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	device1 := models.Device{
		UID:      "1",
		TenantID: "tenant1",
		Tags: []string{
			"device1",
			"device2",
			"device3",
		},
	}

	device2 := models.Device{
		UID:      "2",
		TenantID: "tenant2",
		Tags: []string{
			"device1",
			"device2",
			"device3",
		},
	}

	device3 := models.Device{
		UID:      "3",
		TenantID: "tenant1",
		Tags: []string{
			"device1",
			"device2",
			"device3",
		},
	}

	_, err := db.Client().Database("test").Collection("devices").InsertOne(ctx, &device1)
	assert.NoError(t, err)

	_, err = db.Client().Database("test").Collection("devices").InsertOne(ctx, &device2)
	assert.NoError(t, err)

	_, err = db.Client().Database("test").Collection("devices").InsertOne(ctx, &device3)
	assert.NoError(t, err)

	err = mongostore.TagRename(ctx, "tenant1", "device2", "device4")
	assert.NoError(t, err)

	d1, err := mongostore.DeviceGetByUID(ctx, models.UID(device1.UID), "tenant1")
	assert.NoError(t, err)
	assert.Equal(t, len(d1.Tags), 3)
	assert.Equal(t, d1.Tags[1], "device4")

	d2, err := mongostore.DeviceGetByUID(ctx, models.UID(device2.UID), "tenant2")
	assert.NoError(t, err)
	assert.Equal(t, len(d2.Tags), 3)
	assert.Equal(t, d2.Tags[1], "device2")

	d3, err := mongostore.DeviceGetByUID(ctx, models.UID(device3.UID), "tenant1")
	assert.NoError(t, err)
	assert.Equal(t, len(d3.Tags), 3)
	assert.Equal(t, d3.Tags[1], "device4")
}

func TestDeleteTag(t *testing.T) {
	db := dbtest.DBServer{}
	defer db.Stop()

	ctx := context.TODO()
	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	device1 := models.Device{
		UID:       "1",
		Namespace: "namespace1",
		TenantID:  "tenant1",
		Tags: []string{
			"device1",
			"device5",
			"device3",
		},
	}

	device2 := models.Device{
		UID:       "2",
		Namespace: "namespace1",
		TenantID:  "tenant1",
		Tags: []string{
			"device1",
			"device5",
			"device6",
		},
	}

	device3 := models.Device{
		UID:       "3",
		Namespace: "namespace2",
		TenantID:  "tenant2",
		Tags: []string{
			"device1",
			"device5",
			"device6",
		},
	}

	_, err := db.Client().Database("test").Collection("devices").InsertOne(ctx, &device1)
	assert.NoError(t, err)

	_, err = db.Client().Database("test").Collection("devices").InsertOne(ctx, &device2)
	assert.NoError(t, err)

	_, err = db.Client().Database("test").Collection("devices").InsertOne(ctx, &device3)
	assert.NoError(t, err)

	err = mongostore.TagDelete(ctx, "tenant1", "device1")
	assert.NoError(t, err)

	d1, err := mongostore.DeviceGetByUID(ctx, models.UID(device1.UID), "tenant1")
	assert.NoError(t, err)
	assert.Equal(t, len(d1.Tags), 2)
	assert.Equal(t, d1.Tags, []string{"device5", "device3"})

	d2, err := mongostore.DeviceGetByUID(ctx, models.UID(device2.UID), "tenant1")
	assert.NoError(t, err)
	assert.Equal(t, len(d2.Tags), 2)
	assert.Equal(t, d2.Tags, []string{"device5", "device6"})

	d3, err := mongostore.DeviceGetByUID(ctx, models.UID(device3.UID), "tenant2")
	assert.NoError(t, err)
	assert.Equal(t, len(d3.Tags), 3)
	assert.Equal(t, d3.Tags, device3.Tags)
}
