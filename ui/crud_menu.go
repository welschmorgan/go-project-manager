package ui

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/welschmorgan/go-project-manager/models"
)

var errMenuQuit = errors.New("user quit")

type CRUDMenu struct {
	workspace    *models.Workspace
	key          string
	subKey       string
	refItem      interface{}
	actions      []CRUDAction
	actionLabels map[uint8]string
	items        []interface{}
	names        []string
	indices      map[string]int
}

func NewCRUDMenu(wksp *models.Workspace, key, subKey string, refItem interface{}, actions []CRUDAction, actionLabels map[uint8]string) (*CRUDMenu, error) {
	menu := &CRUDMenu{
		workspace:    wksp,
		key:          key,
		subKey:       subKey,
		refItem:      refItem,
		actions:      actions,
		actionLabels: actionLabels,
		items:        []interface{}{},
		names:        make([]string, 0),
		indices:      map[string]int{},
	}
	if menu.actionLabels == nil || len(menu.actionLabels) == 0 {
		menu.actionLabels = map[uint8]string{
			ActionAdd.Id:    "Add new item",
			ActionEdit.Id:   "Edit existing item",
			ActionRemove.Id: "Remove existing item",
			ActionClear.Id:  "Clear items",
		}
	}
	rv := reflect.Indirect(reflect.ValueOf(wksp))
	rf := rv.FieldByName(key)
	for i := 0; i < rf.Len(); i++ {
		menu.items = append(menu.items, reflect.Indirect(rf.Index(i)).Interface())
	}
	menu.Update()
	return menu, nil
}

func (m *CRUDMenu) Get(name string) interface{} {
	if id, ok := m.indices[name]; ok {
		return m.items[id]
	} else {
		return nil
	}
}

func (m *CRUDMenu) Edit(id int, newItem interface{}) error {
	if id < 0 || id >= len(m.items) {
		return errors.New("invalid project")
	}
	m.items[id] = newItem
	m.Update()
	return nil
}

func (m *CRUDMenu) Create(newItem interface{}) {
	m.items = append(m.items, reflect.Indirect(reflect.ValueOf(newItem)).Interface())
	m.Update()
}

func (m *CRUDMenu) Remove(name string) {
	if id, ok := m.indices[name]; ok {
		m.items = append(m.items[:id], m.items[id+1:]...)
		m.Update()
	}
}

func (m *CRUDMenu) Clear() {
	m.items = []interface{}{}
	m.Update()
}

func (m *CRUDMenu) RenderItems() {
	s := fmt.Sprintf("Found %d items: ", len(m.items))
	for id, item := range m.items {
		if id > 0 {
			s += ", "
		}
		rv := reflect.ValueOf(item)
		if rv.Kind() == reflect.Ptr {
			rv = reflect.Indirect(rv)
		}
		rf := rv.FieldByName(m.subKey)
		s += rf.String()
	}
	println(s)
}

func (m *CRUDMenu) Update() {
	if m.items == nil {
		m.items = []interface{}{}
	}
	m.names = []string{}
	m.indices = map[string]int{}
	for _, p := range m.items {
		rv := reflect.ValueOf(p)
		if rv.Kind() == reflect.Ptr {
			rv = reflect.Indirect(rv)
		}
		rf := rv.FieldByName(m.subKey)
		m.names = append(m.names, rf.String())
		m.indices[rf.String()] = len(m.names) - 1
	}
}

func (m *CRUDMenu) SelectAction() (CRUDAction, error) {
	actionNames := []string{}
	for _, a := range m.actions {
		actionNames = append(actionNames, a.String())
	}
	if action, err := Select("Action", actionNames, nil); err != nil {
		return CRUDAction{Id: 0, Name: ""}, err
	} else if ret, err := ParseCRUDAction(action); err != nil {
		return CRUDAction{Id: 0, Name: ""}, err
	} else {
		return ret, nil
	}
}

func (m *CRUDMenu) Render() error {
	var done bool = false
	var err error
	for !done {
		if err = m.RenderOnce(); err != nil {
			if err == errMenuQuit {
				return nil
			}
			return err
		}
	}
	return nil
}

func (m *CRUDMenu) RenderOnce() error {
	var project string
	var err error
	var action CRUDAction
	m.Update()
	m.RenderItems()
	if action, err = m.SelectAction(); err != nil {
		return err
	}
	if action == ActionRemove || action == ActionEdit {
		if project, err = Select(m.actionLabels[action.Id], m.names, nil); err != nil {
			return err
		}
	}
	defaultProject := m.refItem
	if action == ActionEdit {
		defaultProject = m.Get(project)
	}
	if action == ActionEdit || action == ActionAdd {
		if res, err := AskObject(m.actionLabels[action.Id], defaultProject, nil); err != nil {
			return err
		} else if action == ActionEdit {
			rv := reflect.Indirect(reflect.ValueOf(defaultProject))
			rf := rv.FieldByName(m.subKey)
			if err := m.Edit(m.indices[rf.String()], res); err != nil {
				return err
			}
		} else if action == ActionAdd {
			m.Create(res)
		}
	}
	if action == ActionRemove {
		m.Remove(project)
	}
	if action == ActionClear {
		m.Clear()
	}
	if action == ActionQuit {
		return errMenuQuit
	}
	return nil
}
