package main

import (
	"flag"
	"log"
	"os"
	"github.com/morrah77/messagesender/schedule"
	"github.com/morrah77/messagesender/transport"
)

var (
	scheduleConf *schedule.Conf
	transportConf *transport.Conf
	logger *log.Logger
)

func init() {
	setupLog()
	fillConf()
}

func main() {
	var transp *transport.Transport = transport.NewTransport(logger, transportConf)
	var sched *schedule.Schedule = schedule.NewSchedule(logger, scheduleConf)
	sched.ParseShedules()
	sched.Run(transp.Send)
}

func setupLog() {
	logger = log.New(os.Stdout, `Message sender`, log.Flags())
}

func fillConf() {
	scheduleConf = &schedule.Conf{}
	transportConf = &transport.Conf{}
	flag.StringVar(&(scheduleConf.SourcePath), `file`, `docs/customers.csv`, `Path to CSV file containing shedules`)
	flag.StringVar(&(scheduleConf.CsvDelimiter), `csv-delimiter`, `,`, `CSV fields delimiter`)
	flag.StringVar(&(scheduleConf.ScheduleDelimiter), `schedule-delimiter`, `-`, `schedule field delimiter`)
	flag.StringVar(&(transportConf.SendUrl), `url`, `localhost:9090/messages/`, `url to send messages`)
	flag.Parse()
}

