package todo

import (
	"fmt"
	"time"
)

var sample []RadioProgram = []RadioProgram{
	{
		name:      "よなよな水曜日",
		weekday:   time.Wednesday,
		startTime: "22:00",
		endTime:   "23:30",
	},
	{
		name:      "三四郎のオールナイトニッポン",
		weekday:   time.Friday,
		startTime: "22:00",
		endTime:   "22:30",
	},
}

type ProgramRegister struct {
	todo *TODO
}

func NewProgramRegister(todo *TODO) *ProgramRegister {
	return &ProgramRegister{
		todo: todo,
	}
}

func (r *ProgramRegister) RegisterAndRun() {
	rr := NewRadioRemainder(r.todo)

	sheetDatas, _ := ReadGoogleSheets()
	var programs []RadioProgram
	for _, row := range sheetDatas.Values {
		//fmt.Printf("%s, %s, %s, %s\n", row[0], row[1], row[2], row[3])
		programs = append(programs, RadioProgram{
			name:      fmt.Sprint(row[0]),
			weekday:   Weekday(fmt.Sprint(row[1])),
			startTime: fmt.Sprint(row[2]),
			endTime:   fmt.Sprint(row[3]),
		})
	}
	for _, p := range programs {
		rr.AddProgram(p)
	}
	rr.Do()
}

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
