package webserver_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sgaunet/template-api/pkg/webserver"
	"github.com/stretchr/testify/assert"
)

func TestWebserverStart(t *testing.T) {
	// mockSvc := authors.NewService(nil)
	var wg sync.WaitGroup
	w, err := webserver.NewWebServer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wg.Add(1)
	go func() {
		err = w.Start()
		assert.Nil(t, err, "unexpected error: %v", err)
		wg.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	err = w.Shutdown(context.Background())
	assert.Nil(t, err, "unexpected error: %v", err)
	wg.Wait()
}

func TestWebserverStartTwiceOnSamePort(t *testing.T) {
	// mockSvc := authors.NewService(nil)
	var wg sync.WaitGroup
	w, err := webserver.NewWebServer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wg.Add(1)
	go func() {
		err = w.Start()
		assert.Nil(t, err, "unexpected error: %v", err)
		wg.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	wg.Add(1)
	go func() {
		err = w.Start()
		assert.NotNil(t, err, "expected error, got nil")
		wg.Done()
	}()
	time.Sleep(200 * time.Millisecond)
	err = w.Shutdown(context.Background())
	assert.Nil(t, err, "unexpected error: %v", err)
	wg.Wait()
}
