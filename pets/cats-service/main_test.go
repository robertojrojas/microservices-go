package main


import (
	"testing"
	"os"
	"strings"
	"sync"
)

type dummySignal struct {
	Message string
}

func (s *dummySignal) String() string {
	return s.Message
}
func (s *dummySignal) Signal() {}


func Test_ShutdownHook_OK(t *testing.T) {

	errChan := make(chan error, 1)
	signalChan := make(chan os.Signal, 1)

	var expectedErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func(){
			wg.Done()
		}()
		expectedErr = shutdownHook(signalChan, errChan)

	}()

	signalChan <- &dummySignal{
		Message: "testSignal",
	}

	wg.Wait()
	errStr := expectedErr.Error()

	if !strings.Contains(errStr, "testSignal") {
		t.Error("Expected error to contain the string testSignal")
	}
}

