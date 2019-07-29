package main

import (
	"strings"
	"encoding/json"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)


// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {

	name := m.Name
	nameSlice := strings.Split(name, ".")
	if nameSlice[0] == "component" {
		nameSlice := nameSlice[2:]
		name = strings.Join(nameSlice[:],".")
	}

	switch name {
		case "sample.message3":
			astilog.Info("sample.message3 recieved")
			if err = json.Unmarshal(m.Payload, &payload); err != nil {
				return nil, err	
			}
			return payload, nil
			break
	}
	return nil, err
}
