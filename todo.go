package todo

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/KentaKudo/go-todo-graphql/internal/pb/service"
	log "github.com/sirupsen/logrus"
)

func MarshalTodoStatus(s service.Todo_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if _, err := io.WriteString(w, strconv.Quote(s.String())); err != nil {
			log.Fatal(err)
		}
	})
}

func UnmarshalTodoStatus(v interface{}) (service.Todo_Status, error) {
	switch i := v.(type) {
	case string:
		v, ok := service.Todo_Status_value[fmt.Sprintf("TODO_STATUS_%s", i)]
		if !ok {
			return service.TODO_STATUS_UNKNOWN, fmt.Errorf("invalid Todo Status")
		}
		return service.Todo_Status(v), nil
	case int32:
		return service.Todo_Status(i), nil
	}

	return service.TODO_STATUS_UNKNOWN, fmt.Errorf("invalid type")
}
