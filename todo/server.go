package todo

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
	"todocore/entity"

	"github.com/bluemon0919/timeext"
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
	addr      string
	td        *TODO
	remainder *RadioRemainder
}

type ListItem struct {
	Title    string
	Weekday  string
	Deadline string
}

type By func(p1, p2 *ListItem) bool

// planetSorter joins a By function and a slice of Planets to be sorted.
type itemSorter struct {
	items []ListItem
	by    func(p1, p2 *ListItem) bool // Closure used in the Less method.
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(items []ListItem) {
	s := &itemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(s)
}

// Len is part of sort.Interface.
func (s *itemSorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *itemSorter) Less(i, j int) bool {
	return s.by(&s.items[i], &s.items[j])
}

// NewServer creates new Server
func NewServer(addr string, ent entity.Entity) *Server {
	td := NewTODO(ent)
	r := NewRadioRemainder(td)
	srv := &Server{
		addr:      addr,
		td:        td,
		remainder: r,
	}
	return srv
}

// GetList アクティブな番組一覧を取得します
func (srv *Server) GetList() ([]ListItem, error) {
	srv.remainder.Run()
	items, _ := srv.td.GetActive()

	var lis []ListItem
	for _, item := range items {
		t30 := timeext.TimeExt(item.Deadline)
		li := ListItem{
			Title:    item.Title,
			Weekday:  item.Deadline.Weekday().String(),
			Deadline: t30.Format(Layout),
		}
		lis = append(lis, li)
	}

	// Keyでソートする
	keysort := func(p1, p2 *ListItem) bool {
		return p1.Deadline < p2.Deadline
	}
	By(keysort).Sort(lis)
	return lis, nil
}

// PostProgram IDで指定した番組の聴取状態を"完了"に更新します
func (srv *Server) PostProgram(id string) error {
	iid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	items, _ := srv.td.GetActive()
	isExist := false
	for _, item := range items {
		if item.ID == iid {
			isExist = true
			break
		}
	}
	if !isExist {
		return fmt.Errorf("No item was found with the specified ID[%d]", iid)
	}

	return srv.td.ChangeStatus(iid, COMPLETE)
}

// StartService starts http server.
func (srv *Server) StartService() error {
	http.HandleFunc("/list", srv.listHandler)
	http.HandleFunc("/post", srv.postHandler)
	http.HandleFunc("/play", srv.playHandler)
	return http.ListenAndServe(srv.addr, nil)
}

// listHandler writes list to http.ResponseWriter
// display the list in html
func (srv *Server) listHandler(w http.ResponseWriter, r *http.Request) {
	lis, _ := srv.GetList()
	tpl := template.Must(template.ParseFiles("static/list.html"))
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
		t30 := timeext.TimeExt(item.Deadline)
		value := r.FormValue(item.Title + t30.Format(Layout))
		if "" != value {
			srv.td.ChangeStatus(item.ID, COMPLETE)
		}
	}
	http.Redirect(w, r, "/list", http.StatusFound)
}

// playHandler reflects the input result from html in the list
func (srv *Server) playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	items, _ := srv.td.GetActive()

	var url string
	for _, item := range items {
		t30 := timeext.TimeExt(item.Deadline)
		value := r.FormValue(item.Title + t30.Format(Layout))
		if value != "" {
			// radikoのURLを取得する
			// 一週間前の番組を取得。終了時間が登録されているので、終了１分前に調整。
			m, _ := time.ParseDuration("-1m")
			t := item.Deadline.AddDate(0, 0, -6).Add(m)
			GetTimeshiftURL := func(stationID string, start time.Time) string {
				const defaultEndpoint = "http://radiko.jp"
				location, _ := time.LoadLocation("Asia/Tokyo")
				localTime := start.In(location)
				endpoint := "share/?sid=" + stationID + "&t=" + localTime.Format("20060102150405")
				return defaultEndpoint + "/" + endpoint
			}
			url = GetTimeshiftURL(item.StationID, t)
			if url != "" {
				break
			}
		}
	}
	http.Redirect(w, r, url, http.StatusFound)
}
