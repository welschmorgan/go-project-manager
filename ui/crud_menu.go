package ui

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/welschmorgan/go-release-manager/config"
)

var errMenuQuit = errors.New("user quit")

type CRUDMenu struct {
	Workspace      *config.Workspace
	Key            string
	SubKey         string
	RefItem        interface{}
	Validators     []ObjValidator
	Actions        []CRUDAction
	ActionLabels   map[uint8]string
	Items          []interface{}
	ItemFieldTypes map[string]ItemFieldType
	Names          []string
	Indices        map[string]int
	Finalizer      func(item interface{}) error
}

func NewCRUDMenu(wksp *config.Workspace, key, subKey string, refItem interface{}, validators []ObjValidator, actions []CRUDAction, actionLabels map[uint8]string, itemFieldTypes map[string]ItemFieldType, finalizer func(item interface{}) error) (*CRUDMenu, error) {
	menu := &CRUDMenu{
		Workspace:      wksp,
		Key:            key,
		SubKey:         subKey,
		RefItem:        refItem,
		Validators:     validators,
		Actions:        actions,
		ActionLabels:   actionLabels,
		Items:          []interface{}{},
		ItemFieldTypes: itemFieldTypes,
		Names:          make([]string, 0),
		Indices:        map[string]int{},
		Finalizer:      finalizer,
	}
	if menu.ActionLabels == nil || len(menu.ActionLabels) == 0 {
		menu.ActionLabels = map[uint8]string{
			ActionAdd.Id:    "Add new item",
			ActionEdit.Id:   "Edit existing item",
			ActionRemove.Id: "Remove existing item",
			ActionClear.Id:  "Clear items",
		}
	}
	rv := reflect.Indirect(reflect.ValueOf(wksp))
	rf := rv.FieldByName(key)
	for i := 0; i < rf.Len(); i++ {
		menu.Items = append(menu.Items, reflect.Indirect(rf.Index(i)).Interface())
	}
	menu.Update()
	if err := menu.Discover(); err != nil {
		return nil, err
	}
	menu.Update()
	return menu, nil
}

func (m *CRUDMenu) Get(name string) interface{} {
	if id, ok := m.Indices[name]; ok {
		return m.Items[id]
	} else {
		return nil
	}
}

func (m *CRUDMenu) Edit(id int, newItem interface{}) error {
	if id < 0 || id >= len(m.Items) {
		return errors.New("invalid project")
	}
	if m.Finalizer != nil {
		if err := m.Finalizer(newItem); err != nil {
			return err
		}
	}
	m.Items[id] = newItem
	m.Update()
	return nil
}

func (m *CRUDMenu) Create(newItem interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(newItem)).Interface()
	if m.Finalizer != nil {
		if err := m.Finalizer(v); err != nil {
			return err
		}
	}
	m.Items = append(m.Items, v)
	m.Update()
	return nil
}

func (m *CRUDMenu) Remove(name string) {
	if id, ok := m.Indices[name]; ok {
		m.Items = append(m.Items[:id], m.Items[id+1:]...)
		m.Update()
	}
}

func (m *CRUDMenu) Clear() {
	m.Items = []interface{}{}
	m.Update()
}

func (m *CRUDMenu) RenderItems() {
	s := fmt.Sprintf("Found %d items: ", len(m.Items))
	for id, item := range m.Items {
		if id > 0 {
			s += ", "
		}
		rv := reflect.ValueOf(item)
		if rv.Kind() == reflect.Ptr {
			rv = reflect.Indirect(rv)
		}
		rf := rv.FieldByName(m.SubKey)
		s += rf.String()
	}
	println(s)
}

func (m *CRUDMenu) Update() {
	if m.Items == nil {
		m.Items = []interface{}{}
	}
	m.Names = []string{}
	m.Indices = map[string]int{}
	for _, p := range m.Items {
		rv := reflect.ValueOf(p)
		if rv.Kind() == reflect.Ptr {
			rv = reflect.Indirect(rv)
		}
		rf := rv.FieldByName(m.SubKey)
		m.Names = append(m.Names, rf.String())
		m.Indices[rf.String()] = len(m.Names) - 1
	}
}

func (m *CRUDMenu) SelectAction() (CRUDAction, error) {
	actionNames := []string{}
	for _, a := range m.Actions {
		actionNames = append(actionNames, a.String())
	}
	if action, err := Select("Action", actionNames); err != nil {
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
				done = true
			} else {
				return err
			}
		}
	}
	wrv := reflect.Indirect(reflect.ValueOf(m.Workspace))
	wrf := wrv.FieldByName(m.Key)
	mrv := reflect.Indirect(reflect.ValueOf(m.Items))
	wrf.Set(reflect.MakeSlice(wrf.Type(), len(m.Items), cap(m.Items)))
	refType := reflect.TypeOf(m.RefItem)
	for i := 0; i < wrf.Len(); i++ {
		itemRT := refType
		if itemRT.Kind() == reflect.Ptr {
			itemRT = itemRT.Elem()
		}
		newItem := reflect.New(itemRT)
		itemRV := mrv.Index(i)
		for j := 0; j < itemRT.NumField(); j++ {
			fv := itemRV.Elem().Field(j)
			for _, v := range m.Validators {
				if err := v(fv.Type().Name(), fv.String()); err != nil {
					return err
				}
			}
			nfv := reflect.Indirect(newItem.Elem().Field(j))
			nfv.Set(fv)
		}
		// newItem.Set(mrv.Index(i))
		wrf.Index(i).Set(newItem)
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
		if project, err = Select(m.ActionLabels[action.Id], m.Names); err != nil {
			return err
		}
	}
	defaultProject := m.RefItem
	if action == ActionEdit {
		defaultProject = m.Get(project)
	}
	if action == ActionEdit || action == ActionAdd {
		if res, err := AskObject(m.ActionLabels[action.Id], defaultProject, m.ItemFieldTypes, m.Validators...); err != nil {
			return err
		} else if action == ActionEdit {
			rv := reflect.Indirect(reflect.ValueOf(defaultProject))
			rf := rv.FieldByName(m.SubKey)
			if err := m.Edit(m.Indices[rf.String()], res); err != nil {
				return err
			}
		} else if action == ActionAdd {
			if err := m.Create(res); err != nil {
				return err
			}
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

func (m *CRUDMenu) Discover() error {
	return nil
}
