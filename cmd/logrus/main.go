package main

import "github.com/jayvib/golog"

func main() {
	l := golog.NewLogrusLogger(golog.InfoLevel)
	l.Fatal("Hello")
}
