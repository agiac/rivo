package rivo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func TestRunErrorSync(t *testing.T) {
	t.Run("should handle a batch of errors", func(t *testing.T) {
		// Given
		var receivedErrors []error
		handler := func(ctx context.Context, errs <-chan error) {
			for err := range errs {
				receivedErrors = append(receivedErrors, err)
			}
		}

		// When
		errsCh, wait := rivo.RunErrorSync(context.Background(), handler)

		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		errsCh <- err1
		errsCh <- err2

		wait()

		// Then
		assert.Equal(t, []error{err1, err2}, receivedErrors)
	})
}

func TestRunErrorSyncFunc(t *testing.T) {
	t.Run("should handle errors one by one", func(t *testing.T) {
		// Given
		var receivedErrors []error
		handler := func(ctx context.Context, err error) {
			receivedErrors = append(receivedErrors, err)
		}

		// When
		errsCh, wait := rivo.RunErrorSyncFunc(context.Background(), handler)

		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		errsCh <- err1
		errsCh <- err2

		wait()

		// Then
		assert.ElementsMatch(t, []error{err1, err2}, receivedErrors)
	})
}
