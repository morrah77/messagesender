package schedule

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
)

var conf *Conf = &Conf{
	SourcePath:        `testdata/right.csv`,
	CsvDelimiter:      `,`,
	ScheduleDelimiter: `-`,
}

var wrongHeaderonf *Conf = &Conf{
	SourcePath:        `testdata/wrong_header.csv`,
	CsvDelimiter:      `,`,
	ScheduleDelimiter: `-`,
}

var wrongFieldsonf *Conf = &Conf{
	SourcePath:        `testdata/wrong_fields.csv`,
	CsvDelimiter:      `,`,
	ScheduleDelimiter: `-`,
}

var buf = bytes.NewBuffer(make([]byte, 1024))

var logger = log.New(buf, `test-schedule`, log.Flags())

type messagesCounter struct {
	counters     map[string]int32
	sentMessages []*message
	sync.Mutex
}

var sentMessagesCounter messagesCounter = messagesCounter{
	counters:     make(map[string]int32, 2),
	sentMessages: make([]*message, 0),
}

var runFunc RunFunc = func(i interface{}, j interface{}) error {
	msg, ok := i.(*message)
	if !ok {
		return errors.New(`Incorrect message!`)
	}

	sentMessagesCounter.Lock()
	sentMessagesCounter.sentMessages = append(sentMessagesCounter.sentMessages, i.(*message))
	sentMessagesCounter.counters[msg.Email]++
	currentCounter := sentMessagesCounter.counters[msg.Email]
	sentMessagesCounter.Unlock()

	j.(*message).Email = msg.Email
	j.(*message).Text = msg.Text

	j.(*message).Paid = false
	if msg.Email == `test2@test.com` && currentCounter > 0 {
		j.(*message).Paid = true
	}

	return nil
}

func TestNewSchedule(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, conf)
	if schedule.conf.CsvDelimiter != `,` {
		t.Error(`Incorrect conf.CsvDelimiter on creation!`)
	}
	if schedule.conf.ScheduleDelimiter != `-` {
		t.Error(`Incorrect conf.ScheduleDelimiter on creation!`)
	}
	if schedule.conf.SourcePath != `testdata/right.csv` {
		t.Error(`Incorrect conf.SourcePath on creation!`)
	}
}

func TestSchedule_ParseShedules(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, conf)
	var err = schedule.ParseShedules()
	if err != nil {
		t.Error(err)
	}
	if schedule.header.email != 0 {
		t.Error(`Incorrect email header parsing!`)
	}
	if schedule.header.text != 1 {
		t.Error(`Incorrect text header parsing!`)
	}
	if schedule.header.schedule != 2 {
		t.Error(`Incorrect schedule header parsing!`)
	}
}

func TestSchedule_ParseShedules2(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, wrongHeaderonf)
	var err = schedule.ParseShedules()
	if err == nil {
		t.Error(`Schedules parsing should fail on wrong header!`)
	}
}

func TestSchedule_ParseShedules3(t *testing.T) {
	var schedule *Schedule = NewSchedule(logger, wrongFieldsonf)
	var err = schedule.ParseShedules()
	if err != nil {
		t.Error(fmt.Sprintf(`Should skip lines containing wrong fields during schedules parsing! Returned error: %#v`, err.Error()))
	}
}

func TestSchedule_Run(t *testing.T) {
	defer sentMessagesCounter.Unlock()
	var schedule *Schedule = NewSchedule(logger, conf)
	var err = schedule.ParseShedules()
	if err != nil {
		t.Error(err)
	}
	schedule.Run(runFunc)
	sentMessagesCounter.Lock()
	if len(sentMessagesCounter.sentMessages) != 3 {
		t.Error(`Incorrect schedule execution!`)
	}

	if sentMessagesCounter.counters[`test1@test.com`] != 2 {
		t.Error(`Incorrect schedule execution for non-paid loan!`)
	}
	if sentMessagesCounter.counters[`test2@test.com`] != 1 {
		t.Error(`Incorrect schedule execution for paid loan!`)
	}
}
