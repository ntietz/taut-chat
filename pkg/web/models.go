package web

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/typesense/typesense-go/v2/typesense"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

type User struct {
	ID      string `json:"id"`
	Handle  string `json:"handle"`
	Credits int64  `json:"credits"`
}

type Message struct {
	ID string `json:"id"`

	Sender    string `json:"from_id"`
	Recipient string `json:"to_id"`

	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

func CreateUser(ts *typesense.Client, handle string) (bool, error) {
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
		id := handle
		user := User{
			ID:      id,
			Handle:  handle,
			Credits: 100,
		}
		fmt.Println("Created new user")
		_, err = ts.Collection("users").Documents().Create(ctx, user)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, nil
	}
}

func CreateMessage(ts *typesense.Client, from string, to string, content string) error {
	ctx := context.Background()

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	message := Message{
		ID:        id.String(),
		Sender:    from,
		Recipient: to,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	_, err = ts.Collection("messages").Documents().Create(ctx, message)
	return err
}

func SearchMessages(ts *typesense.Client, currentUser string, query string) ([]Message, error) {
	ctx := context.Background()

	filter := fmt.Sprintf("from_id:=%s || to_id:=%s", currentUser, currentUser)
	qparams := api.SearchCollectionParams{
		Q:        pointer.String(query),
		QueryBy:  pointer.String("content"),
		FilterBy: pointer.String(filter),
		SortBy:   pointer.String("timestamp:desc"),
        HighlightStartTag: pointer.String("<b>"),
        HighlightEndTag: pointer.String("</b>"),
        HighlightFullFields: pointer.String("content"),

	}

	messageRecords, err := ts.Collection("messages").Documents().Search(ctx, &qparams)

	if err != nil {
		fmt.Println("err?", err)
		return nil, err
	}
	fmt.Println("search. err?", err, "found?", (*(*messageRecords).Found))

	messages := make([]Message, 0)

	for _, messageRecord := range *(*messageRecords).Hits {
        content := fmt.Sprintf("%v", (*(*messageRecord.Highlights)[0].Value))
        fmt.Println("messageRecord. highlights?", len(*messageRecord.Highlights))
		message := Message{
			ID:        (*messageRecord.Document)["id"].(string),
			Sender:    (*messageRecord.Document)["from_id"].(string),
			Recipient: (*messageRecord.Document)["to_id"].(string),
			Content:   content,
			Timestamp: int64((*messageRecord.Document)["timestamp"].(float64)),
		}
		messages = append(messages, message)
	}

	return messages, nil

}

func ListUserHandles(ts *typesense.Client) ([]string, error) {
	ctx := context.Background()
	query := api.SearchCollectionParams{
		Q:       pointer.String("*"),
		QueryBy: pointer.String("handle"),
	}
	userRecords, err := ts.Collection("users").Documents().Search(ctx, &query)
	fmt.Println("handles. err?", err, "found?", (*(*userRecords).Found))
	if err != nil {
		return nil, err
	}

	handles := make([]string, 0)

	for _, userRecord := range *(*userRecords).Hits {
		handle := (*userRecord.Document)["handle"].(string)
		handles = append(handles, handle)
	}

	return handles, nil
}

func ListMessages(ts *typesense.Client, from string, to string) ([]Message, error) {
	msgs, err := ListMessagesOneWay(ts, from, to)
	if err != nil {
		return nil, err
	}

	if from != to {
		msgsSwapped, err := ListMessagesOneWay(ts, to, from)
		if err != nil {
			return nil, err
		}

		msgs = slices.Concat(msgs, msgsSwapped)
	}

	slices.SortFunc(msgs, func(a Message, b Message) int {
		return int(a.Timestamp - b.Timestamp)
	})

	return msgs, nil
}

func ListMessagesOneWay(ts *typesense.Client, from string, to string) ([]Message, error) {
	ctx := context.Background()

	filter := fmt.Sprintf("from_id:=%s", from)

	query := api.SearchCollectionParams{
		Q:        pointer.String(to),
		QueryBy:  pointer.String("to_id"),
		FilterBy: pointer.String(filter),
		SortBy:   pointer.String("timestamp:desc"),
	}

	messageRecords, err := ts.Collection("messages").Documents().Search(ctx, &query)

	if err != nil {
		fmt.Println("err?", err)
		return nil, err
	}
	fmt.Println("list messages. err?", err, "found?", (*(*messageRecords).Found))

	messages := make([]Message, 0)

	for _, messageRecord := range *(*messageRecords).Hits {
		message := Message{
			ID:        (*messageRecord.Document)["id"].(string),
			Sender:    (*messageRecord.Document)["from_id"].(string),
			Recipient: (*messageRecord.Document)["to_id"].(string),
			Content:   (*messageRecord.Document)["content"].(string),
			Timestamp: int64((*messageRecord.Document)["timestamp"].(float64)),
		}
		messages = append(messages, message)
	}

	return messages, nil

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
	ctx := context.Background()

	userSchema := &api.CollectionSchema{
		Name: "users",
		Fields: []api.Field{
			{
				Name: "handle",
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
            {
                Name: "content",
                Type: "string",
            },
			{
				Name: "timestamp",
				Type: "int64",
			},
		},
	}

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
