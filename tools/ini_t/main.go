package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

func main() {
	cfg, err := ini.Load("../../vendors.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	fmt.Println(cfg)
	fmt.Println(cfg.SectionStrings())
}
