package ms

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
	Action gameAction `json:"action"`
	Row    int        `json:"row"`
	Col    int        `json:"col"`
}
