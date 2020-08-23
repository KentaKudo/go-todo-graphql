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
	mapping := map[model.Status]service.Todo_Status{
		model.StatusCreated:    service.TODO_STATUS_CREATED,
		model.StatusInProgress: service.TODO_STATUS_IN_PROGRESS,
		model.StatusDone:       service.TODO_STATUS_DONE,
	}

	req := &service.UpdateRequest{
		Todo: &service.Todo{
			Id:          id,
			Title:       *input.Title,
			Description: *input.Description,
			Status:      mapping[*input.Status],
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

func (r *todoResolver) Status(ctx context.Context, obj *service.Todo) (model.Status, error) {
	mapping := map[service.Todo_Status]model.Status{
		service.TODO_STATUS_UNKNOWN:     model.StatusCreated,
		service.TODO_STATUS_CREATED:     model.StatusCreated,
		service.TODO_STATUS_IN_PROGRESS: model.StatusInProgress,
		service.TODO_STATUS_DONE:        model.StatusDone,
	}

	return mapping[obj.GetStatus()], nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
