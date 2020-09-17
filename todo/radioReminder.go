package todo

import (
	"fmt"
	"time"

	"github.com/bluemon0919/timeext"
	"google.golang.org/api/sheets/v4"
)

const checkInterval = 10

// RadioRemainder ラジオリマインダー機能
type RadioRemainder struct {
	programs     []RadioProgram
	todo         *TODO
	sheetDatas   sheets.ValueRange // googleスプレッドシートから取得したデータ
	flgReadSheet bool
}

// RadioProgram 番組情報
type RadioProgram struct {
	name      string       // 番組名
	weekday   time.Weekday // 曜日
	startTime string       // 開始時間
	endTime   string       // 終了時間
	stationID string       // ラジオ局
}

// RadioLayout ラジオ番組の時刻レイアウト
const RadioLayout = "15:04"

// NewRadioRemainder 新しいラジオリマインダーを作成する
func NewRadioRemainder(todo *TODO) *RadioRemainder {
	return &RadioRemainder{
		todo: todo,
	}
}

func (r *RadioRemainder) Update() error {
	// GoogleSheetsからリストを取得する
	datas, err := ReadGoogleSheets()
	r.sheetDatas = *datas
	r.flgReadSheet = true
	return err
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

	// GoogleSheetsから取得したリストを処理する
	if !r.flgReadSheet {
		r.Update()
	}
	var programs []RadioProgram
	for _, row := range r.sheetDatas.Values {
		flg := true

		p := RadioProgram{
			weekday:   Weekday(fmt.Sprint(row[1])),
			startTime: fmt.Sprint(row[2]),
			endTime:   fmt.Sprint(row[3]),
			stationID: fmt.Sprint(row[4]),
		}
		_, nextEnd, _ := NextDate(&p)

		for _, a := range actives {
			if a.Title == fmt.Sprint(row[0]) &&
				a.Deadline.Format(Layout) == nextEnd {
				flg = false
				break
			}
		}
		if !flg {
			continue
		}
		for _, c := range completes {
			if c.Title == fmt.Sprint(row[0]) &&
				c.Deadline.Format(Layout) == nextEnd {
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
			stationID: fmt.Sprint(row[4]),
		})
	}

	// プログラムを登録する
	for _, p := range programs {
		startTime, endTime, err := NextDate(&p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		r.todo.Add(p.name, "", startTime, endTime, p.stationID)
	}
}

// NextDate 次の日程を設定する
func NextDate(program *RadioProgram) (string, string, error) {
	// ロケーション情報を取得
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", "", err
	}

	// 指定ロケーション(Asia/Tokyo)の現在時刻を取得
	now := time.Now().In(loc)

	// 指定ロケーション(Asia/Tokyo)の該当時刻を作成
	start, _ := timeext.Parse(RadioLayout, program.startTime)
	end, _ := timeext.Parse(RadioLayout, program.endTime)
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), start.Hour(), start.Minute(), 0, 0, loc)
	nextEnd := time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, loc)

	// 24:00以降の時刻の場合は調整する
	isext, err := timeext.IsExt(RadioLayout, program.startTime)
	if err != nil {
		return "", "", err
	}
	wd := program.weekday - nextStart.Weekday()
	if wd < 0 {
		wd += 7
	}
	if isext {
		wd++
		if wd >= 7 {
			wd -= 7
		}
	}
	nextStart = nextStart.AddDate(0, 0, int(wd))

	isext, err = timeext.IsExt(RadioLayout, program.endTime)
	if err != nil {
		return "", "", err
	}
	wd = program.weekday - nextEnd.Weekday()
	if wd < 0 {
		wd += 7
	}
	if isext {
		wd++
		if wd >= 7 {
			wd -= 7
		}
	}
	nextEnd = nextEnd.AddDate(0, 0, int(wd))
	return nextStart.Format(Layout), nextEnd.Format(Layout), nil
}

// Run will get triggered automatically.
func (r RadioRemainder) Run() {
	// 期限切れのプログラムを削除する
	items, err := r.todo.GetDeadline(DeadlineExpired)
	if err != nil {
		fmt.Println(err)
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
