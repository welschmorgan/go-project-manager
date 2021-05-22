package ui

import (
	"fmt"
	"strings"
)

type CRUDAction struct {
	Id   uint8
	Name string
}

var (
	ActionAdd    = CRUDAction{Id: 0, Name: "Add"}
	ActionEdit   = CRUDAction{Id: 1, Name: "Edit"}
	ActionClear  = CRUDAction{Id: 2, Name: "Clear"}
	ActionRemove = CRUDAction{Id: 3, Name: "Remove"}
	ActionQuit   = CRUDAction{Id: 4, Name: "Quit"}

	All = []CRUDAction{
		ActionAdd, ActionEdit, ActionClear, ActionRemove, ActionQuit,
	}
)

func (a CRUDAction) String() string {
	return a.Name
}

func ParseCRUDAction(s string) (CRUDAction, error) {
	for _, a := range All {
		if strings.ToLower(a.Name) == strings.ToLower(s) {
			return a, nil
		}
	}
	return CRUDAction{0, ""}, fmt.Errorf("unknown CRUD action '%s'", s)
}
