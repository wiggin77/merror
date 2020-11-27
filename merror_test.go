package merror

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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
