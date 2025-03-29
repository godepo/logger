package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/godepo/groat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type State struct {
	Message string
	Attrs   []slog.Attr
}

type Deps struct {
	MockLogger *MockLogger
}

func newTestCase(t *testing.T) *groat.Case[Deps, State, context.Context] {
	tc := groat.New[Deps, State, context.Context](t,
		func(t *testing.T, deps Deps) context.Context {
			return context.WithValue(context.Background(), logInstance, deps.MockLogger)
		},
		func(t *testing.T, deps Deps) Deps {
			deps.MockLogger = NewMockLogger(t)
			return deps
		},
	)
	tc.Go()
	return tc
}

func TestInfo(t *testing.T) {
	tc := newTestCase(t)
	tc.Given(
		ArrangeMessageString,
		ArrangeStringAttr,
	).When(
		ActLogInfo,
	)

	Info(tc.SUT, tc.State.Message, tc.State.Attrs...)
}

func TestDebug(t *testing.T) {
	tc := newTestCase(t)
	tc.Given(
		ArrangeMessageString,
		ArrangeStringAttr,
	).When(
		ActLogDebug,
	)

	Debug(tc.SUT, tc.State.Message, tc.State.Attrs...)
}

func TestError(t *testing.T) {
	tc := newTestCase(t)
	tc.Given(
		ArrangeMessageString,
		ArrangeStringAttr,
	).When(
		ActLogError,
	)

	Error(tc.SUT, tc.State.Message, tc.State.Attrs...)
}

func TestWarn(t *testing.T) {
	tc := newTestCase(t)
	tc.Given(
		ArrangeMessageString,
		ArrangeStringAttr,
	).When(
		ActLogWarn,
	)

	Warn(tc.SUT, tc.State.Message, tc.State.Attrs...)
}

func TestWith(t *testing.T) {
	t.Run("should be able to be able", func(t *testing.T) {
		t.Run("when call empty context for debug", func(t *testing.T) {
			ctx := With(context.Background(), slog.String("key", "value"))
			Debug(ctx, "test")
		})
		t.Run("when call empty context for error", func(t *testing.T) {
			ctx := With(context.Background(), slog.String("key", "value"))
			Error(ctx, "test")
		})
		t.Run("when call empty context for warn", func(t *testing.T) {
			ctx := With(context.Background(), slog.String("key", "value"))
			Warn(ctx, "test")
		})
		t.Run("when call empty context for info", func(t *testing.T) {
			ctx := With(context.Background(), slog.String("key", "value"))
			Info(ctx, "test")
		})
	})

}

func ActLogDebug(_ *testing.T, deps Deps, state State) State {
	deps.MockLogger.EXPECT().Debug(mock.Anything, state.Message, sliceToAny(state.Attrs)...)
	return state
}

func ActLogWarn(_ *testing.T, deps Deps, state State) State {
	deps.MockLogger.EXPECT().Warn(mock.Anything, state.Message, sliceToAny(state.Attrs)...)
	return state
}

func ActLogError(_ *testing.T, deps Deps, state State) State {
	deps.MockLogger.EXPECT().Error(mock.Anything, state.Message, sliceToAny(state.Attrs)...)
	return state
}

func sliceToAny[T any](a []T) []any {
	res := make([]any, len(a))
	for i, v := range a {
		res[i] = v
	}
	return res
}

func ActLogInfo(_ *testing.T, deps Deps, state State) State {
	deps.MockLogger.EXPECT().Info(mock.Anything, state.Message, sliceToAny(state.Attrs)...)
	return state
}

func ArrangeStringAttr(_ *testing.T, state State) State {
	state.Attrs = append(state.Attrs, slog.String(uuid.NewString(), uuid.NewString()))
	return state
}

func ArrangeMessageString(_ *testing.T, state State) State {
	state.Message = uuid.NewString()
	return state
}
