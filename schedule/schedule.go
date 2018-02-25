package schedule

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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

type sortableDelays []time.Duration

func (sdItem sortableDelays) Len() int {
	return len(sdItem)
}

func (sdItem sortableDelays) Less(i, j int) bool {
	return sdItem[i] < sdItem[j]
}

func (sdItem sortableDelays) Swap(i, j int) {
	sdItem[i], sdItem[j] = sdItem[j], sdItem[i]
}

func newSortableDelays(length int) sortableDelays {
	return make([]time.Duration, length)
}

type task struct {
	start   time.Time
	delays  []time.Duration
	Message *message
}

func (t *task) calculateSleepTime(i int) (time.Duration, error) {
	//TODO(h.lazar) consider to remove index validation due to this method is not purposed to be used out of delays range
	if i < 0 || i+1 > len(t.delays) {
		return 0, errors.New(`Invalid task delays index`)
	}
	if i == 0 {
		return t.delays[i], nil
	} else {
		return (t.delays[i] - t.delays[i-1]), nil
	}
}

type scheduleHeader struct {
	text     int
	email    int
	schedule int
}

type RunFunc func(interface{}, interface{}) error

type Schedule struct {
	conf        *Conf
	header      *scheduleHeader
	tasks       []*task
	logger      *log.Logger
	stopChannel chan struct{}
}

func NewSchedule(logger *log.Logger, conf *Conf) *Schedule {
	return &Schedule{
		conf:   conf,
		header: &scheduleHeader{},
		tasks:  make([]*task, 0),
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
	if len(chunks) < sch.header.schedule+1 {
		return errors.New(`Too few fields in given line!`)
	}
	durations, err := sch.splitShedule(chunks[sch.header.schedule])
	if err != nil {
		return err
	}
	task := &task{
		Message: mess,
	}
	task.delays = durations
	sch.tasks = append(sch.tasks, task)
	return nil
}

//TODO(h.lazar) let's don't use any reflection here to do not decrease performance
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

func (sch *Schedule) splitShedule(line string) (sortableDelays, error) {
	var err error
	chunks := strings.Split(string(line), sch.conf.ScheduleDelimiter)
	if len(chunks) < 1 {
		err = errors.New(`Invalid Schedule format!`)
		return nil, err
	}
	durations := newSortableDelays(0)
	for _, chunk := range chunks {
		dur, err := time.ParseDuration(chunk)
		if err != nil {
			sch.logger.Println(err.Error())
			continue
		}
		durations = append(durations, dur)
	}
	sort.Sort(durations)
	return durations, err
}

func (sch *Schedule) Run(rf RunFunc) error {
	if sch.stopChannel != nil {
		return errors.New(`Shedule is already runned!`)
	}
	sch.stopChannel = make(chan struct{}, len(sch.tasks))
	finishedTasksCounter := 0
	for _, t := range sch.tasks {
		go func(t *task) {
			t.start = time.Now()
			//TODO(h.lazar) consider to check wheather it's time to send each n ms (improve precision but decrease performance)
			for i, _ := range t.delays {
				//TODO(h.lazar) consider to calculate sleep time in-place (allows to exclude error checking but decrease readability)
				timeToSleep, err := t.calculateSleepTime(i)
				if err != nil {
					sch.logger.Println(err.Error())
					break
				}
				time.Sleep(timeToSleep)
				var r *message = &message{}
				err = rf(t.Message, r)
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
