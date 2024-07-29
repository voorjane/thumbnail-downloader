package test

import (
	"getThumb/internal"
	"testing"
)

func TestClient(t *testing.T) {
	urls := []string{"https://www.youtube.com/watch?v=wBxOKQpdESg", "https://www.youtube.com/watch?v=teSZhXIUjmw", "https://youtu.be/dQw4w9WgXcQ"}
	err := internal.RunClient(urls, true)
	if err != nil {
		t.Fatal("error not expected")
	}
	err = internal.RunClient(urls, false)
	if err != nil {
		t.Fatal("error not expected")
	}
}
