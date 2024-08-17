package web

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/typesense/typesense-go/v2/typesense"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"github.com/typesense/typesense-go/v2/typesense/api/pointer"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id"`
	Handle       string `json:"handle"`
	PasswordHash string `json:"password_hash"`
	Credits      int64  `json:"credits"`
}

type Message struct {
	ID string `json:"id"`

	Sender    string `json:"from_id"`
	Recipient string `json:"to_id"`

	Content string `json:"content"`
}

func CheckLogin(ts *typesense.Client, handle string, password string) (bool, error) {
	ctx := context.Background()
	query := api.SearchCollectionParams{
		Q:       pointer.String(handle),
		QueryBy: pointer.String("handle"),
	}
	matchingUsers, err := ts.Collection("users").Documents().Search(ctx, &query)
	if err != nil {
		return false, err
	}

	count := *matchingUsers.Found

	if count == 0 {
		id, err := uuid.NewV7()
		if err != nil {
			return false, err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return false, err
		}
		user := User{
			ID:           id.String(),
			Handle:       handle,
			PasswordHash: string(hash),
			Credits:      100,
		}
        fmt.Println("Created new user")
		_, err = ts.Collection("users").Documents().Create(ctx, user)
		if err != nil {
			return false, err
		}

        return true, nil
	} else {
		// There's a bug here where if there are multiple users with the same
		// handle somehow, there will be multiple records. YOLO!

		hash, ok := (*(*matchingUsers.Hits)[0].Document)["password_hash"].(string)
        if !ok {
            return false, errors.New("invalid password hash")
        }
        matched := (bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil)
        fmt.Println("password matched?", matched, ", password?", password)
        return matched, nil
	}
}

func DropCollections(ts *typesense.Client) error {
	collections, err := ts.Collections().Retrieve(context.Background())
	if err != nil {
		return err
	}

	for _, collection := range collections {
		_, err = ts.Collection(collection.Name).Delete(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateCollections(ts *typesense.Client) error {
	userSchema := &api.CollectionSchema{
		Name: "users",
		Fields: []api.Field{
			{
				Name: "handle",
				Type: "string",
			},
			{
				Name: "password_hash",
				Type: "string",
			},
			{
				Name: "credits",
				Type: "int64",
			},
		},
	}

	messageSchema := &api.CollectionSchema{
		Name: "messages",
		Fields: []api.Field{
			{
				Name: "from_id",
				Type: "string",
			},
			{
				Name: "to_id",
				Type: "string",
			},
		},
	}

	ctx := context.Background()

	_, err := ts.Collections().Create(ctx, userSchema)
	if err != nil {
		return err
	}

	_, err = ts.Collections().Create(ctx, messageSchema)
	if err != nil {
		return err
	}

	return nil
}
