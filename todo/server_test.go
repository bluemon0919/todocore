package todo

import (
	"log"
	"net/http/httptest"
	"testing"
	"todocore/entity"
)

func TestServerGetActive(t *testing.T) {
	ent := entity.NewMap()
	srv := NewServer(":8080", ent)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?kind=active", nil)
	srv.handler(w, r)
	resp := w.Result()

	defer resp.Body.Close()
	ctype := resp.Header.Get("Content-Type")
	if "application/json" != ctype {
		log.Fatalf("Content-Type does not match. %s\n", ctype)
	}
}

func TestServerGetComplete(t *testing.T) {
	ent := entity.NewMap()
	srv := NewServer(":8080", ent)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?kind=complete", nil)
	srv.handler(w, r)
	resp := w.Result()

	defer resp.Body.Close()
	ctype := resp.Header.Get("Content-Type")
	if "application/json" != ctype {
		log.Fatalf("Content-Type does not match. %s\n", ctype)
	}
}

func TestServerDelete(t *testing.T) {
	ent := entity.NewMap()
	srv := NewServer(":8080", ent)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/", nil)
	srv.handler(w, r)
	resp := w.Result()

	defer resp.Body.Close()

	// expect
}

func TestServerAdd(t *testing.T) {
	ent := entity.NewMap()
	srv := NewServer(":8080", ent)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/?kind=add", nil)
	srv.handler(w, r)
	resp := w.Result()

	defer resp.Body.Close()
	// expect
}

func TestServerUpdate(t *testing.T) {
	ent := entity.NewMap()
	srv := NewServer(":8080", ent)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/?kind=update", nil)
	srv.handler(w, r)
	resp := w.Result()

	defer resp.Body.Close()
	// expect
}
