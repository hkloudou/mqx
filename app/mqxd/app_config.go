package main

import "fmt"

type config struct {
	StrictMode bool
}

func (m config) String() string {
	return fmt.Sprintf("strict_mode:%v", m.StrictMode)
}
