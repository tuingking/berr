package berr

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(err error) {
	b, _ := json.MarshalIndent(err, "", "\t")
	fmt.Println(string(b))
}

func PrintStack(err error, opt ...Option) {
	if bErr, ok := err.(*Error); ok && bErr != nil {
		fmt.Println("stacktrace:")
		for i := len(bErr.stack) - 1; i >= 0; i-- {
			fmt.Printf("\t|> %s\n", bErr.stack[i].Print(opt...))
		}
		fmt.Printf("\t|> %v\n", bErr.err)
	}
}

type Option func(*options)

type options struct {
	showFile  bool
	showLine  bool
	shortName bool
	shortFile bool
}

func defaultOptions() *options {
	return &options{
		showFile:  true,
		showLine:  true,
		shortName: true,
		shortFile: false,
	}
}

func PrintWithFile(v bool) Option {
	return func(o *options) {
		o.showFile = v
	}
}
func PrintWithLine(v bool) Option {
	return func(o *options) {
		o.showLine = v
	}
}
func PrintWithShortFunc(v bool) Option {
	return func(o *options) {
		o.shortName = v
	}
}
func PrintWithShortFile(v bool) Option {
	return func(o *options) {
		o.shortFile = v
	}
}
