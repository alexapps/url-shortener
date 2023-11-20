package random_test

import (
	"testing"

	"github.com/alexapps/url-shortener/internal/lib/random"
)

func TestNewRandomString(t *testing.T) {
	if len(random.NewRandomString(7)) != 7 {
		t.Fatalf("Generated string length has wrong len")
	}
}
