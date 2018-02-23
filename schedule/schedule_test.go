package schedule

import (
	"testing"
	"log"
	"bytes"
)

var conf *Conf  = &Conf{
	SourcePath: `testdata/right.csv`,
	CsvDelimiter: `,`,
	ScheduleDelimiter:`-`,
}
var buf = bytes.NewBuffer(make([]byte, 1024))
var logger = log.New(buf, `test-transport`, log.Flags())

var sentMessages []*message = make([]*message, 0)

var runFunc RunFunc = func(i interface{}) error {
	sentMessages = append(sentMessages, i.(*message))
}

func TestNewSchedule(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, conf)
	if schedule.conf.ScheduleDelimiter != `-` {
		t.Error(`Incorrect conf.ScheduleDelimiter on creation!`)
	}
}

func TestSchedule_ParseShedules(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, conf)
	schedule.ParseShedules()
	if schedule.header.email != 0 {
		t.Error(`Incorrect header parsing!`)
	}
}

func TestSchedule_Run(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, conf)
	schedule.ParseShedules()
	schedule.Run(runFunc)
	if len(sentMessages) != 4 {
		t.Error(`Incorrect schedule execution!`)
	}
}