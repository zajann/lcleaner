package target

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"

	log "github.com/zajann/lcleaner/pkg/easylog"
)

type Target struct {
	path     string
	r        *regexp.Regexp
	baseDate time.Time
	period   string
}

func (t *Target) Clean() error {
	var cnt int
	var total int
	log.Info("Start Clean Log in %s", t.path)
	files, err := ioutil.ReadDir(t.path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if t.r.MatchString(f.Name()) {
			total++
			if f.ModTime().Before(t.baseDate) {
				fullPath := fmt.Sprintf("%s/%s", t.path, f.Name())
				if err := os.Remove(fullPath); err != nil {
					return err
				}
				log.Debug("Delete %s", fullPath)
				cnt++
			}
		}
	}
	t.updateBaseDate()
	log.Info("Finish Clean Log in %s (%d/%d)", t.path, cnt, total)
	return nil
}

func (t *Target) updateBaseDate() {
	t.baseDate, _ = getBaseDate(t.period)
}

func New(path string, reg string, period string) (*Target, error) {
	t := new(Target)
	var err error

	t.path = path
	t.period = period
	t.r, err = regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	t.baseDate, err = getBaseDate(period)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func getBaseDate(period string) (time.Time, error) {
	now := time.Now()
	if period == "" {
		return now, errors.New("period is undefined")
	}
	var t time.Time
	var p int

	dayR := regexp.MustCompile("^([0-9]+)d$")
	monthR := regexp.MustCompile("^([0-9]+)m$")
	yearR := regexp.MustCompile("^([0-9]+)y$")

	if dayR.MatchString(period) {
		p, _ = strconv.Atoi(dayR.FindStringSubmatch(period)[1])
		t = now.AddDate(0, 0, -p)
	} else if monthR.MatchString(period) {
		p, _ = strconv.Atoi(monthR.FindStringSubmatch(period)[1])
		t = now.AddDate(0, -p, 0)
	} else if yearR.MatchString(period) {
		p, _ = strconv.Atoi(yearR.FindStringSubmatch(period)[1])
		t = now.AddDate(-p, 0, 0)
	} else {
		return now, fmt.Errorf("invalid period type: %s", period)
	}
	return t, nil
}
