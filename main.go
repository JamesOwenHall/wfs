package main

import (
	"fmt"
)

func main() {
	c, err := ReadConfig("config.yml")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}
}
