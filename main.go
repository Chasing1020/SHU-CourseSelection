package main

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const LoginUrl = "https://oauth.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MjE2MDUzMTQ1ODcxMTM5MDMsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6InlSUUxKZlVzeDMyNmZTZUtOVUN0b29LdyIsInNjb3BlIjoiIiwicmVkaXJlY3RVcmkiOiJodHRwOi8veGsuYXV0b2lzcC5zaHUuZWR1LmNuL3Bhc3Nwb3J0L3JldHVybiIsInN0YXRlIjoiIn0="
const TermSelectUrl = "http://xk.autoisp.shu.edu.cn/Home/TermSelect"
const CourseSelectionSaveUrl = "http://xk.autoisp.shu.edu.cn/CourseSelectionStudent/CourseSelectionSave"
const QueryCourseCheckUrl = "http://xk.autoisp.shu.edu.cn/CourseSelectionStudent/QueryCourseCheck"

const QuerySelector = "#tblcoursecheck > tbody > tr:nth-child(2) > td:nth-child(2)"

var count int64

func main() {
	start := time.Now()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-ch
		log.WithFields(log.Fields{
			"signal":    sig,
			"used_time": time.Now().Sub(start),
			"req_times": atomic.LoadInt64(&count),
		}).Infof("Program exit")
		log.Exit(-1)
	}()

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	Login(c)

	var wg sync.WaitGroup
	for _, course := range Conf.Courses {
		ctx, cancel := context.WithCancel(context.Background())
		for i := 0; i < runtime.NumCPU()<<1; i++ {
			cc := c.Clone() // will share the cookie jar
			wg.Add(1)
			go func(ctx context.Context, cancel context.CancelFunc, cc *colly.Collector, id int, course Course) {
				OnQueryCallbacks(ctx, cancel, cc, id, course)
				for {
					select {
					case <-ctx.Done():
						wg.Done()
						return
					default:
						QueryCourse(cc, id, course)
					}
				}
			}(ctx, cancel, cc, i+1, course)
		}
	}
	wg.Wait()

	log.WithFields(log.Fields{
		"used_time": time.Now().Sub(start),
		"req_times": atomic.LoadInt64(&count),
	}).Infof("All courses have been selected!")
}

// Login try to log in to the xk.autoisp.shu.edu.cn.
func Login(c *colly.Collector) {
	err := c.Post(LoginUrl, map[string]string{
		"username": Conf.Username,
		"password": EncryptPassword(Conf.Password),
	})
	if err != nil {
		log.Panic(err)
	}

	err = c.Post(TermSelectUrl, map[string]string{"termId": Conf.TermId})
	if err != nil {
		log.Panic(err)
	}

	log.WithFields(log.Fields{
		"username": Conf.Username,
		"password": Conf.Password,
	}).Info("Login successfully!")
}

// OnQueryCallbacks registers a function.
// It will save the course on every query if the course is not full.
func OnQueryCallbacks(ctx context.Context, cancel context.CancelFunc, c *colly.Collector, id int, course Course) {
	c.OnHTML(QuerySelector, func(e *colly.HTMLElement) {
		if !strings.Contains(e.DOM.Text(), course.CourseId) {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			err := c.Post(CourseSelectionSaveUrl, map[string]string{
				"cids": course.CourseId,
				"tnos": course.TeacherNo,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"id":  id,
					"err": err,
				}).Warn("Post CourseSelectionSaveUrl error")
				return
			}
			log.WithFields(
				CourseFieldsMap[course],
			).Infof("Goroutine %02d: Select successfully!", id)
			cancel()
			return
		}
	})
}

// QueryCourse will try to query every course status.
// If any course is able to save, it will be hooked by QueryCallbacks.
func QueryCourse(c *colly.Collector, id int, course Course) {
	err := c.Post(QueryCourseCheckUrl, map[string]string{
		"CID":            course.CourseId,
		"TeachNo":        course.TeacherNo,
		"CourseName":     course.CourseName,
		"TeachName":      course.TeacherName,
		"FunctionString": "LoadData",
		"IsNotFull":      "true",
		"PageIndex":      "1",
		"PageSize":       "10",
	})
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Warn("Post QueryCourseCheckUrl error")
	}
	atomic.AddInt64(&count, 1)

	log.WithFields(
		CourseFieldsMap[course],
	).Debugf("Goroutine %02d: The %dth time to query the course", id, atomic.LoadInt64(&count))
}
