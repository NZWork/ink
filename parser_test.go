package main

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func timer(startAt int) int {
	return time.Now().Nanosecond() - startAt
}

func TestParse(t *testing.T) {
	now := time.Now().Nanosecond()
	data, err := ioutil.ReadFile("test/parser.md")
	fmt.Printf("read file %v nanosecond\n", timer(now))
	//fmt.Printf("file content %v\n", string(data))

	if err != nil {
		fmt.Printf("ERROR %v\n", err.Error())
		return
	}
	now = time.Now().Nanosecond()
	out := ParseMarkdownToHTML(data)
	fmt.Printf("parse file %v nanosecond\n", timer(now))

	now = time.Now().Nanosecond()
	ioutil.WriteFile("test.html", out, 0644)
	fmt.Printf("write file %v nanosecond\n`", timer(now))
	//fmt.Printf("result: %s\n", string(out))
}

/*
func TestLoad(t *testing.T) {
	start := time.Now().Nanosecond()
	LoadRawMarkdownFiles()
	end := time.Now().Nanosecond()

	fmt.Printf("use time %vms", float64(end-start)/1000.0/1000.0)
}
*/
