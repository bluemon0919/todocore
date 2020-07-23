package todo

import (
	"fmt"
	"time"

	"github.com/bluemon0919/timeext"
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

// Do はラジオリマインダーを実行する
func (r *RadioRemainder) Do() {
	// 登録済みのプログラムリストを取得する
	actives, err := r.todo.GetActive()
	if err != nil {
		return
	}
	completes, err := r.todo.GetComplete()
	if err != nil {
		return
	}

	// GoogleSheetsからリストを取得する
	sheetDatas, _ := ReadGoogleSheets()
	var programs []RadioProgram
	for _, row := range sheetDatas.Values {
		flg := true

		p := RadioProgram{
			weekday: Weekday(fmt.Sprint(row[1])),
			endTime: fmt.Sprint(row[3]),
		}
		next, _ := NextDate(&p)

		for _, a := range actives {
			if a.Title == fmt.Sprint(row[0]) &&
				a.Deadline.Format(Layout) == next {
				flg = false
				break
			}
		}
		if !flg {
			continue
		}
		for _, c := range completes {
			if c.Title == fmt.Sprint(row[0]) &&
				c.Deadline.Format(Layout) == next {
				flg = false
				break
			}
		}
		if !flg {
			continue
		}
		programs = append(programs, RadioProgram{
			name:      fmt.Sprint(row[0]),
			weekday:   Weekday(fmt.Sprint(row[1])),
			startTime: fmt.Sprint(row[2]),
			endTime:   fmt.Sprint(row[3]),
		})
	}

	// プログラムを登録する
	for _, p := range programs {
		deadline, err := NextDate(&p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		r.todo.Add(p.name, "", deadline)
	}
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
	end, _ := timeext.Parse(RadioLayout, program.endTime)
	next := time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, loc)

	// 該当の曜日まで進める
	d := program.weekday - next.Weekday()
	if d <= 0 {
		d += 7
	}
	next = next.AddDate(0, 0, int(d))
	return next.Format(Layout), nil
}

// Run will get triggered automatically.
func (r RadioRemainder) Run() {
	// 期限切れのプログラムを削除する
	items, err := r.todo.GetDeadline(DeadlineExpired)
	if err != nil {
		return
	}
	for _, item := range items {
		r.todo.Delete(item.ID)
	}

	// 次週のプログラムを追加する
	r.Do()
}

// Weekday は曜日の文字列を与えて、それに対応するtimeパッケージのWeekdayの値を返します
func Weekday(weekday string) time.Weekday {
	switch weekday {
	case "Sunday":
		return time.Sunday
	case "Monday":
		return time.Monday
	case "Tuesday":
		return time.Tuesday
	case "Wednesday":
		return time.Wednesday
	case "Thursday":
		return time.Thursday
	case "Friday":
		return time.Friday
	case "Saturday":
		return time.Saturday
	}
	return time.Sunday
}
