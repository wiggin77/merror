package merror

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	NumThreads = 25
	NumErrors  = 5
)

func TestRace(t *testing.T) {
	merr := New()
	wg := &sync.WaitGroup{}
	wg.Add(NumThreads)

	for i := 0; i < NumThreads; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < NumErrors; j++ {
				merr.Append(fmt.Errorf("Goroutine: %d   ErrorNum: %d", id, j))
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, NumThreads*NumErrors, len(merr.Errors()))
}

func TestCapacity(t *testing.T) {
	merr := NewWithCap(10)

	for i := 0; i < 15; i++ {
		merr.Append(errors.New(fmt.Sprintf("error %d", i)))
	}

	assert.Equal(t, 10, merr.Len())
	assert.Equal(t, 5, merr.Overflow())
}

func TestIs(t *testing.T) {
	err1 := ErrOne{"err1"}
	err2 := &ErrTwo{"err2"}
	err3 := ErrThree{"err3"}
	err4 := ErrFour{"err4"}
	err5 := fmt.Errorf("err5 : %w", err4)
	errNope := ErrNope{"nope"}

	merr := New()
	merr.Append(err1, err2, err3, err5)

	assert.True(t, merr.Is(err1))
	assert.True(t, merr.Is(err2)) // pointer
	assert.True(t, merr.Is(err3))
	assert.True(t, merr.Is(err4)) // wrapped
	assert.True(t, merr.Is(err5))

	assert.False(t, merr.Is(errNope))
}

func TestAs(t *testing.T) {
	err1 := ErrOne{"err1"}
	err2 := &ErrTwo{"err2"}
	err3 := ErrThree{"err3"}
	err4 := ErrFour{"err4"}
	err5 := fmt.Errorf("err5 : %w", err4)

	merr := New()
	merr.Append(err1, err2, err3, err5)

	var one ErrOne
	ok := errors.As(merr, &one)
	require.True(t, ok)
	assert.IsType(t, ErrOne{}, one)

	var two *ErrTwo
	ok = errors.As(merr, &two) // pointer
	require.True(t, ok)
	assert.IsType(t, &ErrTwo{}, two)

	var four ErrFour
	ok = errors.As(merr, &four) // wrapped
	require.True(t, ok)
	assert.IsType(t, ErrFour{}, four)

	var nope ErrNope
	ok = errors.As(merr, &nope)
	assert.False(t, ok)
}

type ErrOne struct {
	msg string
}

func (e ErrOne) Error() string {
	return e.msg
}

type ErrTwo struct {
	msg string
}

func (e ErrTwo) Error() string {
	return e.msg
}

type ErrThree struct {
	msg string
}

func (e ErrThree) Error() string {
	return e.msg
}

type ErrFour struct {
	msg string
}

func (e ErrFour) Error() string {
	return e.msg
}

type ErrNope struct {
	msg string
}

func (e ErrNope) Error() string {
	return e.msg
}
