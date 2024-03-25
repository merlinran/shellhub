package mongo

import (
	"context"
	"fmt"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/api/store/mongo/queries"
	"github.com/shellhub-io/shellhub/pkg/api/query"
	"github.com/shellhub-io/shellhub/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Store) AnnouncementList(ctx context.Context, paginator query.Paginator, sorter query.Sorter) ([]models.AnnouncementShort, int, error) {
	query := []bson.M{}
	fmt.Println("store annoucmentes :", paginator, sorter)
	fmt.Println("store annoucmentes :", paginator, sorter)
	fmt.Println("store annoucmentes :", paginator, sorter)
	fmt.Println("store annoucmentes :", paginator, sorter)
	queryCount := append(query, bson.M{"$count": "count"})
	count, err := AggregateCount(ctx, s.db.Collection("announcements"), queryCount)
	if err != nil {
		return nil, 0, FromMongoError(err)
	}

	fmt.Println("store annoucmentes count and err,", count, err)
	fmt.Println("store annoucmentes count and err,", count, err)
	fmt.Println("store annoucmentes count and err,", count, err)
	sorter.By = "date"
	query = append(query, queries.FromSorter(&sorter)...)
	query = append(query, queries.FromPaginator(&paginator)...)

	cursor, err := s.db.Collection("announcements").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, FromMongoError(err)
	}

	var announcements []models.AnnouncementShort
	if err := cursor.All(ctx, &announcements); err != nil {
		return nil, 0, FromMongoError(err)
	}

	return announcements, count, nil
}

func (s *Store) AnnouncementGet(ctx context.Context, uuid string) (*models.Announcement, error) {
	ann := new(models.Announcement)

	err := s.db.Collection("announcements").FindOne(ctx, bson.M{"uuid": uuid}).Decode(&ann)
	if err != nil {
		return nil, FromMongoError(err)
	}

	return ann, nil
}

func (s *Store) AnnouncementCreate(ctx context.Context, announcement *models.Announcement) error {
	fmt.Println("crete;", s.db.Name())
	fmt.Println("crete;", s.db.Client())
	if _, err := s.db.Collection("announcements").InsertOne(ctx, announcement); err != nil {
		return FromMongoError(err)
	}

	return nil
}

func (s *Store) AnnouncementUpdate(ctx context.Context, announcement *models.Announcement) error {
	result, err := s.db.Collection("announcements").UpdateOne(ctx, bson.M{"uuid": announcement.UUID}, bson.M{"$set": bson.M{"title": announcement.Title, "content": announcement.Content}})
	if err != nil {
		return FromMongoError(err)
	}

	if result.MatchedCount < 1 {
		return store.ErrNoDocuments
	}

	return nil
}

func (s *Store) AnnouncementDelete(ctx context.Context, uuid string) error {
	result, err := s.db.Collection("announcements").DeleteOne(ctx, bson.M{"uuid": uuid})
	if err != nil {
		return FromMongoError(err)
	}
	fmt.Println("result", result)
	fmt.Println("result", result)
	fmt.Println("result", result)
	if result.DeletedCount < 1 {
		return store.ErrNoDocuments
	}

	return nil
}
