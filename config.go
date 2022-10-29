package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strconv"
	"time"
)

var Conf Configuration
var CourseFieldsMap = make(map[Course]log.Fields)

type Configuration struct {
	TermId   string   `yaml:"termId"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Courses  []Course `yaml:"courses"`
}

type Course struct {
	CourseId    string `yaml:"courseId"`
	TeacherNo   string `yaml:"teacherNo"`
	CourseName  string `yaml:"courseName,omitempty"`
	TeacherName string `yaml:"teacherName,omitempty"`
}

func (c Course) ToLogFields() log.Fields {
	fields := log.Fields{}
	courseValue := reflect.ValueOf(c)
	courseType := reflect.TypeOf(c)
	for i := 0; i < courseValue.NumField(); i++ {
		k := courseType.Field(i).Name
		v := courseValue.Field(i).String()
		if v != "" {
			fields[k] = v
		}
	}
	return fields
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
