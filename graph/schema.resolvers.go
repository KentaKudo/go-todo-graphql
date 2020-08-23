package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/KentaKudo/go-todo-graphql/graph/generated"
	"github.com/KentaKudo/go-todo-graphql/graph/model"
	"github.com/KentaKudo/go-todo-graphql/internal/pb/service"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, title string, description string) (string, error) {
	req := &service.CreateRequest{
		Title:       title,
		Description: description,
	}

	resp, err := r.TodoClient.Create(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetId(), nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input model.NewTodo) (bool, error) {
	var (
		title, description string
		status             service.Todo_Status
	)

	if t := input.Title; t != nil {
		title = *t
	}

	if d := input.Description; d != nil {
		description = *d
	}

	if s := input.Status; s != nil {
		status = *s
	}

	req := &service.UpdateRequest{
		Todo: &service.Todo{
			Id:          id,
			Title:       title,
			Description: description,
			Status:      status,
		},
	}

	if _, err := r.TodoClient.Update(ctx, req); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	req := &service.DeleteRequest{
		Id: id,
	}

	if _, err := r.TodoClient.Delete(ctx, req); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Todo(ctx context.Context, id string) (*service.Todo, error) {
	req := &service.GetRequest{
		Id: id,
	}

	resp, err := r.TodoClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetTodo(), nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*service.Todo, error) {
	req := &service.ListRequest{}

	resp, err := r.TodoClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetTodos(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
