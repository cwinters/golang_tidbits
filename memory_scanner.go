package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type MemoryScanner struct {
	src      *bufio.Scanner
	remember string
}

func (s *MemoryScanner) Err() error {
	return s.src.Err()
}

func (s *MemoryScanner) Remember(line string) {
	s.remember = line
}

func (s *MemoryScanner) Scan() bool {
	if s.remember == "" {
		return s.src.Scan()
	}
	return true
}

func (s *MemoryScanner) Text() string {
	if s.remember == "" {
		return s.src.Text()
	}
	t := s.remember
	s.remember = ""
	return t
}

func main() {
	targets := "One\nTwo\nThree\nFour\nFive\n"
	reader := bytes.NewBufferString(targets)
	sc := MemoryScanner{src: bufio.NewScanner(reader)}
	remembered := false
	for i := 0; i < 7; i++ {
		if sc.Scan() {
			fmt.Printf("SCAN OK, iteration %d: %s\n", i, sc.Text())
			if !remembered && sc.Text() == "Three" {
				sc.Remember("Three")
				remembered = true
			}
		} else {
			fmt.Printf("SCAN FAIL: iteration %d\n", i)
			break
		}
	}
}
