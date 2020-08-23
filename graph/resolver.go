//go:generate go run github.com/99designs/gqlgen

package graph

import "github.com/KentaKudo/go-todo-graphql/internal/pb/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoClient service.TodoAPIClient
}
