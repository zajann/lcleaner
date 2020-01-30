package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/zajann/lcleaner/pkg/config"
	log "github.com/zajann/lcleaner/pkg/easylog"
	"github.com/zajann/lcleaner/pkg/target"

	"github.com/jasonlvhit/gocron"
	"github.com/zajann/process"
)

const appVersion = "1.0.0"

var (
	configFilePath string
)

func main() {
	ParseFlag()
	c, err := config.Load(configFilePath)
	if err != nil {
		fmt.Printf("[Error] Failed to load config: %s\n", err)
		os.Exit(1)
	}

	running, _ := process.IsRunning(fmt.Sprintf("%s/%s", c.PIDFilePath, c.PIDFileName))
	if running == 1 {
		fmt.Println("Process is already running")
		os.Exit(0)
	}

	if err := log.Init(
		log.SetFilePath(c.Log.FilePath),
		log.SetFileName(c.Log.FileName),
		log.SetLevel(log.LogLevel(c.Log.Level)),
		log.SetMaxSize(c.Log.MaxSize),
	); err != nil {
		fmt.Printf("[Error] Failed to init log: %s\n", err)
		os.Exit(1)
	}
	c.DumpToLog()

	var targets []*target.Target
	for _, t := range c.Targets {
		nt, err := target.New(t.Path, t.Regexp, t.Period)
		if err != nil {
			log.Fatal("Failed to create new Target: %s", err)
		}
		targets = append(targets, nt)
	}

	log.Info("lcleaner Start")
	gocron.Every(1).Days().At("00:00").Do(Start, targets)
	<-gocron.Start()
}

func Start(targets []*target.Target) {
	for _, t := range targets {
		if err := t.Clean(); err != nil {
			log.Error("Failed to clean log: %s", err)
		}
	}
}

func ParseFlag() {
	var regexp string
	var target string

	t := flag.Bool("t", false, "test regular expression")
	flag.StringVar(&regexp, "regexp", "", "regular expression of log file, required for test")
	flag.StringVar(&target, "target", "", "test log file name, required for test")
	v := flag.Bool("v", false, "version")
	flag.StringVar(&configFilePath, "c", "", "configuration file path, required")

	flag.Parse()
	if *v {
		PrintVersion()
	}
	if *t {
		TestRegexp(regexp, target)
	}
	if configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func TestRegexp(r string, t string) {
	if r == "" || t == "" {
		flag.Usage()
		os.Exit(1)
	}
	var result bool

	reg := regexp.MustCompile(r)
	if reg.MatchString(t) {
		result = true
	}
	fmt.Println(result)
	os.Exit(0)
}

func PrintVersion() {
	fmt.Printf("lcleaner Version %s\n", appVersion)
	os.Exit(0)
}
