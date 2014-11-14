package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type PeekingScanner struct {
	src      *bufio.Scanner
	remember string
}

func (s *PeekingScanner) Err() error {
	return s.src.Err()
}

func (s *PeekingScanner) Peek() string {
	if !s.src.Scan() {
		return ""
	}
	s.remember = s.src.Text()
	return s.remember
}

func (s *PeekingScanner) Scan() bool {
	if s.remember == "" {
		return s.src.Scan()
	}
	return true
}

func (s *PeekingScanner) Text() string {
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
	sc := PeekingScanner{src: bufio.NewScanner(reader)}
	remembered := false
	for i := 0; i < 7; i++ {
		if sc.Scan() {
			fmt.Printf("SCAN OK, iteration %d: %s\n", i, sc.Text())
			if !remembered && sc.Peek() == "Three" {
				remembered = true
				fmt.Printf("...peek returned true, next Scan should still return Three...\n")
			}
		} else {
			fmt.Printf("SCAN FAIL: iteration %d\n", i)
			break
		}
	}
}
