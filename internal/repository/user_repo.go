package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sumbul/music-player-backend/internal/models"
	"google.golang.org/api/iterator"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	UpdateLastSeen(id string) error
}

type userRepo struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) UserRepository {
	return &userRepo{client}
}

func (r *userRepo) Create(user *models.User) error {
	ctx := context.Background()
	_, _, err := r.client.Collection("users").Add(ctx, user)
	return err
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	iter := r.client.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	user.ID = doc.Ref.ID
	return &user, nil
}

func (r *userRepo) GetByID(id string) (*models.User, error) {
	ctx := context.Background()
	doc, err := r.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	user.ID = doc.Ref.ID
	return &user, nil
}

func (r *userRepo) UpdateLastSeen(id string) error {
	ctx := context.Background()
	_, err := r.client.Collection("users").Doc(id).Update(ctx, []firestore.Update{
		{Path: "last_seen", Value: time.Now()},
	})
	return err
}
