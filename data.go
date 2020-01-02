package main

import (
	"errors"
	"fmt"
)

type gameAction string

func (ga gameAction) validate() error {

	switch ga {
	case "mark":
	case "clear":
	default:
		return errors.New(fmt.Sprintf("unknown action %s", ga))
	}
	return nil
}

//json data sent by client API
type inputData struct {
	action gameAction
	row    int
	col    int
}
