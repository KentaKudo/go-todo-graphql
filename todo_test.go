package todo

import (
	"bytes"
	"testing"

	"github.com/KentaKudo/go-todo-graphql/internal/pb/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalTodoStatus(t *testing.T) {
	var buf bytes.Buffer
	MarshalTodoStatus(service.TODO_STATUS_CREATED).MarshalGQL(&buf)

	assert.Equal(t, `"TODO_STATUS_CREATED"`, buf.String())
}

func TestUnmarshalTodoStatus(t *testing.T) {
	tts := map[string]struct {
		input interface{}
		want  service.Todo_Status
	}{
		"created": {
			input: "CREATED",
			want:  service.TODO_STATUS_CREATED,
		},
		"in progress": {
			input: "IN_PROGRESS",
			want:  service.TODO_STATUS_IN_PROGRESS,
		},
		"done": {
			input: "DONE",
			want:  service.TODO_STATUS_DONE,
		},
	}

	for name, tt := range tts {
		t.Run(name, func(t *testing.T) {
			v, err := UnmarshalTodoStatus(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, v)
		})
	}
}
