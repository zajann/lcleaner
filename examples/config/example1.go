package main

import (
	"fmt"
	"regexp"

	"github.com/zajann/lcleaner/pkg/config"
)

func main() {

	filePath := "./lcleaner_config.yml"

	c, err := config.Load(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)

	s := "nvpmon_snmp.log.bak.20191231124455"

	for i, t := range c.Targets {
		r, err := regexp.Compile(t.Regexp)
		if err != nil {
			panic(err)
		}
		fmt.Println(i, r.MatchString(s))
	}
}
