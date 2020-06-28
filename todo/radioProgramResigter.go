package todo

import (
	"time"
)

var programs []RadioProgram = []RadioProgram{
	{
		// name:      "よなよな",
		name:      "yonayona-suiyo-bi",
		weekday:   time.Wednesday,
		startTime: "22:00",
		endTime:   "23:30",
	},
	{
		name:      "sanshiro-no-all night japan",
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
	for _, p := range programs {
		rr.AddProgram(p)
	}
	rr.Do()
}
