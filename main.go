package main

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const LoginUrl = "https://oauth.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MjE2MDUzMTQ1ODcxMTM5MDMsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6InlSUUxKZlVzeDMyNmZTZUtOVUN0b29LdyIsInNjb3BlIjoiIiwicmVkaXJlY3RVcmkiOiJodHRwOi8veGsuYXV0b2lzcC5zaHUuZWR1LmNuL3Bhc3Nwb3J0L3JldHVybiIsInN0YXRlIjoiIn0="
const TermSelectUrl = "http://xk.autoisp.shu.edu.cn/Home/TermSelect"
const CourseSelectionSaveUrl = "http://xk.autoisp.shu.edu.cn/CourseSelectionStudent/CourseSelectionSave"
const QueryCourseCheckUrl = "http://xk.autoisp.shu.edu.cn/CourseSelectionStudent/QueryCourseCheck"

const QuerySelector = "#tblcoursecheck > tbody > tr:nth-child(2) > td:nth-child(2)"

var count int64
var selected = make(map[string]bool)
var rw sync.RWMutex

func main() {
	start := time.Now()
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	Login(c)

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU() << 1)
	for i := 0; i < runtime.NumCPU()<<1; i++ {
		cc := c.Clone() // will share the cookie jar

		go func(cc *colly.Collector, id int) {
			OnQueryCallbacks(cc, id)
			for {
				rw.RLock()
				if len(selected) == len(Conf.Courses) {
					rw.RUnlock()
					wg.Done()
					return
				}
				rw.RUnlock()

				QueryCourse(cc, id)
			}
		}(cc, i+1)
	}
	wg.Wait()

	log.Printf("*** All courses have been selected! running time of %v ***", time.Now().Sub(start))
}

// Login try to log in to the xk.autoisp.shu.edu.cn.
func Login(c *colly.Collector) {
	err := c.Post(LoginUrl, map[string]string{
		"username": Conf.Username,
		"password": EncryptPassword(Conf.Password),
	})
	if err != nil {
		panic(err)
	}

	err = c.Post(TermSelectUrl, map[string]string{"termId": Conf.TermId})
	if err != nil {
		panic(err)
	}
}

// OnQueryCallbacks registers a function.
// It will save the course on every query if the course is not full.
func OnQueryCallbacks(c *colly.Collector, id int) {
	c.OnHTML(QuerySelector, func(e *colly.HTMLElement) {
		defer func() {
			if info := recover(); info != nil {
				log.Printf("Goroutine %2d recovered: %v", id, info)
			}
		}()
		rw.RLock()
		if len(selected) == len(Conf.Courses) {
			rw.RUnlock()
			return
		}
		rw.RUnlock()

		for _, course := range Conf.Courses {
			rw.RLock()
			if !strings.Contains(e.DOM.Text(), course.CourseId) || selected[course.CourseId] {
				rw.RUnlock()
				continue
			}
			rw.RUnlock()

			err := c.Post(CourseSelectionSaveUrl, map[string]string{
				"cids": course.CourseId,
				"tnos": course.TeacherNo,
			})
			if err != nil {
				panic(err)
			}
			log.Printf("=== %v selection successful!!! ===", course)

			rw.Lock()
			selected[course.CourseId] = true
			rw.Unlock()
		}
	})
}

// QueryCourse will try to query every course status.
// If any course is able to save, it will be hooked by QueryCallbacks.
func QueryCourse(c *colly.Collector, id int) {
	defer func() {
		if info := recover(); info != nil {
			log.Printf("Goroutine %2d recovered: %v", id, info)
		}
	}()
	for _, course := range Conf.Courses {
		err := c.Post(QueryCourseCheckUrl, map[string]string{
			"CID":            course.CourseId,
			"TeachNo":        course.TeacherNo,
			"FunctionString": "LoadData",
			"IsNotFull":      "true",
			"PageIndex":      "1",
			"PageSize":       "10",
		})
		if err != nil {
			panic(err)
		}
		atomic.AddInt64(&count, 1)
		log.Printf("Goroutine %2d: The %dth attempt to query the course %v", id, atomic.LoadInt64(&count), course)
	}
}
