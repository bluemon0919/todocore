package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todotool/entity"
)

// Response is server response
type Response struct {
	Items []Item `json:"Item"`
}

// AddRequest is a request parameter of Add function
type AddRequest struct {
	Title  string
	Detail string
}

// UpdateRequest is a request parameter of Update function
type UpdateRequest struct {
	ID     int
	Status int
}

// Server provides http server
type Server struct {
	addr string
	td   *TODO
}

// NewServer creates new Server
func NewServer(addr string, ent entity.Entity) *Server {
	return &Server{
		addr: addr,
		td:   NewTODO(ent),
	}
}

// StartService starts http server.
func (srv *Server) StartService() error {
	http.HandleFunc("/", srv.handler)
	return http.ListenAndServe(srv.addr, nil)
}

// getHandler
func (srv *Server) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		srv.get(w, r)
	case http.MethodDelete:
		srv.delete(w, r)
	case http.MethodPut:
		srv.put(w, r)
	default:
		err := fmt.Sprintf("does not support method %s", r.Method)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func (srv *Server) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)

	kind := r.FormValue("kind")
	var items []Item
	var err error

	switch kind {
	case "active":
		items, err = srv.td.GetActive()
	case "complete":
		items, err = srv.td.GetComplete()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp Response
	for _, item := range items {
		resp.Items = append(resp.Items, item)
	}
	if err := enc.Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv *Server) put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	kind := r.FormValue("kind")

	switch kind {
	case "add":
		var reqs []AddRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, req := range reqs {
			if err := srv.td.Add(req.Title, req.Detail); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	case "update":
		var reqs []UpdateRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, req := range reqs {
			if err := srv.td.ChangeStatus(req.ID, req.Status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	default:
		err := fmt.Sprintf("does not support param %s", kind)
		http.Error(w, err, http.StatusInternalServerError)
	}
}

func (srv *Server) delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var ids []int
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&ids); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, id := range ids {
		if err := srv.td.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
