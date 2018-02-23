package schedule

import (
	"time"
	"log"
	"os"
	"fmt"
	"bufio"
	"errors"
	"strings"
)

type Conf struct{
	SourcePath string
	CsvDelimiter string
	ScheduleDelimiter string
}

type message struct {
	email string `json:"email"`
	text string `json:"text"`
}

type Task struct {
	delay   time.Duration
	Message *message
}

type scheduleHeader struct {
	text int
	email int
	schedule int
}
type Schedule struct {
	conf *Conf
	header *scheduleHeader
	tasks []*Task
	logger *log.Logger
}

type RunFunc func(interface{}) error

func NewSchedule(logger *log.Logger, conf *Conf) *Schedule {
	return &Schedule{
		conf: conf,
		header:&scheduleHeader{},
		tasks: make([]*Task, 0),
		logger: logger,
	}
}

func(sch *Schedule) ParseShedules() error {
	var (
		err error
		currentLineNumber int
	)
	file, err := os.Open(sch.conf.SourcePath)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf(`Could not open file %s!`, sch.conf.SourcePath))
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if currentLineNumber > 0 {
			sch.parseShedule(line)
			continue
		}
		if err = sch.parseHeader(line); err != nil {
			return err
		}
	}
	return nil
}

func(sch *Schedule) parseShedule(line string) error {
	chunks, err := sch.splitLine(line)
	if err != nil {
		return err
	}
	mess := &message{
		email:chunks[sch.header.email],
		text:chunks[sch.header.text],
	}
	durations, err := sch.splitShedule(chunks[sch.header.schedule])
	for _, dur := range durations {
		task := &Task{
			dur,
			mess,
		}
		sch.tasks = append(sch.tasks, task)
	}
	return nil
}

func(sch *Schedule) parseHeader(line string) error {
	chunks, err := sch.splitLine(line)
	if err != nil {
		return err
	}
	for i, chunk := range chunks {
		if chunk == `email`{
			sch.header.email = i
		}
		if chunk == `text`{
			sch.header.text = i
		}
		if chunk == `Schedule`{
			sch.header.schedule = i
		}
	}
	if sch.header.email + sch.header.text + sch.header.schedule < 3 {
		return errors.New(`Invalid file header format!`)
	}
	return nil
}

func(sch *Schedule) splitLine(line string) (chunks []string, err error) {
	chunks = strings.Split(string(line), sch.conf.CsvDelimiter)
	if len(chunks) < 3 {
		err = errors.New(`Invalid file format!`)
	}
	return chunks, err
}

func(sch *Schedule) splitShedule(line string) (durations []time.Duration, err error) {
	chunks := strings.Split(string(line), sch.conf.ScheduleDelimiter)
	if len(chunks) < 1 {
		err = errors.New(`Invalid Schedule format!`)
	}
	for _, chunk := range chunks {
		dur, err := time.ParseDuration(chunk)
		if err == nil {
			durations = append(durations, dur)
		}
	}
	return durations, err
}

//TODO(h.lazar) fix it!
func(sch *Schedule) Run(rf RunFunc) error {
	for _, t := range sch.tasks{
		err := rf(t.Message)
		if err != nil {
			sch.logger.Println(err.Error())
		}
	}
	return nil
}

