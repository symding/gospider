package gospider

import (
	"fmt"
	"log"
	"time"
)

func newSpiderStat() spiderStat {
	s := spiderStat{}
	s.Start()
	return s
}

type spiderStat struct {
	StartTime       time.Time
	EndTime         time.Time
	RequestCount    int
	ResponseCount   int
	StatusCodeCount map[int]int
	ErrorCount      map[error]int
}

// Start spider
func (s *spiderStat) Start() {
	s.StartTime = time.Now()
	s.StatusCodeCount = map[int]int{}
	s.ErrorCount = map[error]int{}
}

// Stop spider
func (s *spiderStat) Stop() {
	s.EndTime = time.Now()
	statusStatStr := ""
	for k, v := range s.StatusCodeCount {
		statusStatStr = fmt.Sprintf("%s\n\t\t%d :\t%d", statusStatStr, k, v)
	}
	errorStatStr := ""
	for k, v := range s.ErrorCount {
		errorStatStr = fmt.Sprintf("%s\n\t\t%v :\t%d", errorStatStr, k, v)
	}
	log.Println(
		"\n\tStartAt: ", s.StartTime,
		"\n\tEndAt: ", time.Now(),
		"\n\tTotal Request Count:", s.RequestCount,
		"\n\tTotal Response Count: ", s.ResponseCount,
		"\n\tStatus Code Stat: ", statusStatStr,
		"\n\tError Stat: ", errorStatStr,
	)
	log.Println("Good Bye.")
}

func (s *spiderStat) RequestIncr() {
	s.RequestCount++
}

func (s *spiderStat) ResponseIncr(response Response) {
	if response.Error == nil {
		s.ResponseCount++
		if _, ok := s.StatusCodeCount[response.StatusCode]; !ok {
			s.StatusCodeCount[response.StatusCode] = 0
		}
		s.StatusCodeCount[response.StatusCode]++
	} else {
		if _, ok := s.ErrorCount[response.Error]; !ok {
			s.ErrorCount[response.Error] = 0
		}
		s.ErrorCount[response.Error]++
	}
}
