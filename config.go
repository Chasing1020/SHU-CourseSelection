package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"time"
)

var Conf Configuration

type Configuration struct {
	TermId   string   `yaml:"termId"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Courses  []Course `yaml:"courses"`
}

type Course struct {
	CourseId  string `yaml:"courseId"`
	TeacherNo string `yaml:"teacherNo"`
}

func (c Course) String() string {
	return fmt.Sprintf("<%s, %s>", c.CourseId, c.TeacherNo)
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
		DisableQuote:    true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Panicf("Func `os.ReadFile` failed, details: %v", err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Panicf("Func `yaml.Unmarshal` failed, details: %v", err)
	}

	if Conf.TermId == "3" || Conf.TermId == "5" {
		Conf.TermId = strconv.Itoa(time.Now().Year()-1) + Conf.TermId
	} else {
		Conf.TermId = strconv.Itoa(time.Now().Year()) + Conf.TermId
	}
	log.WithField("termId", Conf.TermId).Info("Read file `config.yaml` successfully")
	log.WithFields(log.Fields{
		"courses": Conf.Courses,
	}).Infof("There %d are courses to be selected", len(Conf.Courses))
}
