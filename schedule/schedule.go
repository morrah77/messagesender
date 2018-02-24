package schedule

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Conf struct {
	SourcePath        string
	CsvDelimiter      string
	ScheduleDelimiter string
}

type message struct {
	Email string `json:"email"`
	Text  string `json:"text"`
	Paid  bool   `json:"paid, omitempty"`
}

type Task struct {
	delays  []time.Duration
	Message *message
}

type scheduleHeader struct {
	text     int
	email    int
	schedule int
}
type Schedule struct {
	conf        *Conf
	header      *scheduleHeader
	tasks       []*Task
	logger      *log.Logger
	stopChannel chan struct{}
}

type RunFunc func(interface{}, interface{}) error

func NewSchedule(logger *log.Logger, conf *Conf) *Schedule {
	return &Schedule{
		conf:   conf,
		header: &scheduleHeader{},
		tasks:  make([]*Task, 0),
		logger: logger,
	}
}

func (sch *Schedule) ParseShedules() error {
	var (
		err               error
		csvReader         *csv.Reader
		currentLineNumber int
		currentFieldset   []string
	)
	file, err := os.Open(sch.conf.SourcePath)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf(`Could not open file %s!`, sch.conf.SourcePath))
	}

	csvReader = csv.NewReader(file)
	csvReader.Comma = []rune(sch.conf.CsvDelimiter)[0]

	for {
		currentFieldset, err = csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			if currentLineNumber == 0 {
				return err
			}
			sch.logger.Printf(`CSV parse parse error at line [%v]: %#v`, currentLineNumber, err)
		}
		if currentLineNumber > 0 {
			err = sch.parseShedule(currentFieldset)
			if err != nil {
				sch.logger.Printf(`Shedule parse error: %#v`, err)
			}
			currentLineNumber++
			continue
		}
		if err = sch.parseHeader(currentFieldset); err == nil {
			csvReader.FieldsPerRecord = len(currentFieldset)
		} else {
			return err
		}
		currentLineNumber++
	}
	return nil
}

func (sch *Schedule) parseShedule(chunks []string) error {
	mess := &message{
		Email: chunks[sch.header.email],
		Text:  chunks[sch.header.text],
	}
	durations, err := sch.splitShedule(chunks[sch.header.schedule])
	if err != nil {
		return err
	}
	task := &Task{
		make([]time.Duration, 0),
		mess,
	}
	for _, dur := range durations {
		task.delays = append(task.delays, dur)
	}
	sch.tasks = append(sch.tasks, task)
	return nil
}

func (sch *Schedule) parseHeader(chunks []string) error {
	for i, chunk := range chunks {
		if chunk == `email` {
			sch.header.email = i
		}
		if chunk == `text` {
			sch.header.text = i
		}
		if chunk == `schedule` {
			sch.header.schedule = i
		}
	}
	if sch.header.email+sch.header.text+sch.header.schedule < 3 {
		return errors.New(`Invalid file header format!`)
	}
	return nil
}

func (sch *Schedule) splitShedule(line string) ([]time.Duration, error) {
	var err error
	chunks := strings.Split(string(line), sch.conf.ScheduleDelimiter)
	if len(chunks) < 1 {
		err = errors.New(`Invalid Schedule format!`)
		return nil, err
	}
	durations := make([]time.Duration, 0)
	for _, chunk := range chunks {
		dur, err := time.ParseDuration(chunk)
		if err != nil {
			sch.logger.Println(err.Error())
			continue
		}
		durations = append(durations, dur)
	}
	return durations, err
}

//TODO(h.lazar) fix it!
func (sch *Schedule) Run(rf RunFunc) error {
	if sch.stopChannel != nil {
		return errors.New(`Shedule is already runned!`)
	}
	sch.stopChannel = make(chan struct{}, len(sch.tasks))
	finishedTasksCounter := 0
	for _, t := range sch.tasks {
		go func(t *Task) {
			for _, duration := range t.delays {
				time.Sleep(duration)
				var r *message = &message{}
				err := rf(t.Message, r)
				if err != nil {
					sch.logger.Println(err.Error())
				}
				if r.Paid {
					sch.logger.Printf(`%#v is paid already!`, t.Message.Email)
					break
				}
			}
			sch.stopChannel <- struct{}{}
		}(t)
	}
	for {
		_ = <-sch.stopChannel
		finishedTasksCounter++
		if finishedTasksCounter >= len(sch.tasks) {
			break
		}

	}
	return nil
}
