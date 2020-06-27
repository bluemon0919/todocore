package todo

import (
	"fmt"
	"time"

	"github.com/bamzi/jobrunner"
)

const checkInterval = 10

// RadioRemainder ラジオリマインダー機能
type RadioRemainder struct {
	programs []RadioProgram
	todo     *TODO
}

// RadioProgram 番組情報
type RadioProgram struct {
	name      string       // 番組名
	weekday   time.Weekday // 曜日
	startTime string       // 開始時間
	endTime   string       // 終了時間
}

// RadioLayout ラジオ番組の時刻レイアウト
const RadioLayout = "15:04"

// NewRadioRemainder 新しいラジオリマインダーを作成する
func NewRadioRemainder(todo *TODO) *RadioRemainder {
	return &RadioRemainder{
		todo: todo,
	}
}

// NewRadioProgram 新しいラジオプログラムを作成する
func NewRadioProgram(name string, weekday time.Weekday, start, end string) *RadioProgram {
	return &RadioProgram{
		name:      name,
		weekday:   weekday,
		startTime: start,
		endTime:   end,
	}
}

// AddProgram 番組を登録する
func (r *RadioRemainder) AddProgram(program RadioProgram) {
	r.programs = append(r.programs, program)
}

// Do ラジオリマインダーを実行する
func (r *RadioRemainder) Do() {
	// 番組をTODOリストに登録する
	for _, p := range r.programs {
		deadline, err := NextDate(&p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		r.todo.Add(p.name, "", deadline)
	}

	// チェッカーを起動する
	r.runnerSetup()
}

// NextDate 次の日程を設定する
func NextDate(program *RadioProgram) (string, error) {
	// ロケーション情報を取得
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", err
	}

	// 指定ロケーション(Asia/Tokyo)の現在時刻を取得
	now := time.Now().In(loc)

	// 指定ロケーション(Asia/Tokyo)の該当時刻を作成
	end, _ := time.Parse(RadioLayout, program.endTime)
	next := time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, loc)

	// 該当の曜日まで進める
	d := program.weekday - next.Weekday()
	if d <= 0 {
		d += 7
	}
	next = next.AddDate(0, 0, int(d))
	return next.Format(Layout), nil
}

// Setup job runner
func (r *RadioRemainder) runnerSetup() {
	spec := fmt.Sprintf("@every %dm", checkInterval)
	jobrunner.Start()
	jobrunner.Schedule(spec, r)
}

// Run will get triggered automatically.
func (r RadioRemainder) Run() {
	items, err := r.todo.GetDeadline(DeadlineExpired)
	//items, err := r.todo.GetComplete() // デバッグ
	if err != nil {
		return
	}
	for _, item := range items {
		// 同じプログラムを登録する
		for _, p := range r.programs {
			if p.name == item.Title {
				deadline, err := NextDate(&p)
				if err != nil {
					fmt.Println(err)
					continue
				}
				r.todo.Delete(item.ID)
				r.todo.Add(p.name, "", deadline)
			}
		}
	}
}
