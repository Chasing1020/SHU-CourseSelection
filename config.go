package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
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

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		panic(err)
	}

	if Conf.TermId == "3" || Conf.TermId == "5" {
		Conf.TermId = strconv.Itoa(time.Now().Year()-1) + Conf.TermId
	} else {
		Conf.TermId = strconv.Itoa(time.Now().Year()) + Conf.TermId
	}
	log.Printf("Read confing.yaml successful: %+v", Conf)
}
