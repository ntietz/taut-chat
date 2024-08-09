package views

import "github.com/typesense/typesense-go/v2/typesense"

type Handler struct {
	ts *typesense.Client
}

func NewHandler() *Handler {
	ts := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("1667b96f-da3c-40f9-a3b5-8b461a78ed68"))
	return &Handler{
		ts,
	}
}
