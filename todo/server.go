package todo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"todocore/entity"
)

// Response is server response
type Response struct {
	Items []Item `json:"Item"`
}

// AddRequest is a request parameter of Add function
type AddRequest struct {
	Title    string
	Detail   string
	Deadline string
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
	srv := &Server{
		addr: addr,
		td:   NewTODO(ent),
	}
	r := NewProgramRegister(srv.td)
	r.RegisterAndRun()
	return srv
}

// StartService starts http server.
func (srv *Server) StartService() error {
	http.HandleFunc("/", srv.handler)
	http.HandleFunc("/list", srv.listHandler)
	http.HandleFunc("/post", srv.postHandler)
	return http.ListenAndServe(srv.addr, nil)
}

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
	deadline := r.FormValue("deadline")
	var items []Item
	var err error

	switch kind {
	case "active":
		if deadline == "today" {
			items, err = srv.td.GetDeadline(DeadlineToday)
		} else {
			items, err = srv.td.GetActive()
		}
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
			//var t time.Time
			if err := srv.td.Add(req.Title, req.Detail, req.Deadline); err != nil {
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

// listHandler writes list to http.ResponseWriter
// display the list in html
func (srv *Server) listHandler(w http.ResponseWriter, r *http.Request) {

	items, _ := srv.td.GetActive()

	type ListItem struct {
		Title    string
		Weekday  string
		Deadline string
	}

	var lis []ListItem
	for _, item := range items {
		li := ListItem{
			Title:    item.Title,
			Weekday:  item.Deadline.Weekday().String(),
			Deadline: item.Deadline.Format(Layout),
		}
		lis = append(lis, li)
	}

	tpl := template.Must(template.ParseFiles("static/list2.html"))
	tpl.Execute(w, lis)
}

// postHandler reflects the input result from html in the list
func (srv *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	items, _ := srv.td.GetActive()
	for _, item := range items {
		value := r.FormValue(item.Title + item.Deadline.Format(Layout))
		if "" != value {
			srv.td.ChangeStatus(item.ID, COMPLETE)
		}
	}
	http.Redirect(w, r, "/list", http.StatusFound)
}
