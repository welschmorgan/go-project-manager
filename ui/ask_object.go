package ui

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type ItemFieldTypeId uint8

const (
	ItemFieldUnknown ItemFieldTypeId = iota
	ItemFieldText                    = iota
	ItemFieldList                    = iota
)

type ItemFieldType struct {
	Id           ItemFieldTypeId
	DefaultValue interface{}
}

func NewItemFieldType(t ItemFieldTypeId, data interface{}) ItemFieldType {
	return ItemFieldType{
		Id:           t,
		DefaultValue: data,
	}
}

func (t ItemFieldTypeId) String() string {
	switch t {
	case ItemFieldList:
		return "list"
	case ItemFieldText:
		return "text"
	default:
		return "unknown"
	}
}

func AskObjectItem(label string, itemFieldType ItemFieldType, validators ...Validator) (string, error) {
	switch itemFieldType.Id {
	case ItemFieldText:
		if val, err := Ask(label, itemFieldType.DefaultValue, validators...); err != nil {
			return "", err
		} else {
			return val, nil
		}
	case ItemFieldList:
		if val, err := Select(label, itemFieldType.DefaultValue.([]string)); err != nil {
			return "", err
		} else {
			return val, nil
		}
	default:
		return "", fmt.Errorf("cannot ask item for field type '%s'", itemFieldType.Id.String())
	}
}

func AskObject(label string, defaultValue interface{}, itemFieldTypes map[string]ItemFieldType, validators ...ObjValidator) (interface{}, error) {
	reflectedDefaultValue := reflect.Indirect(reflect.ValueOf(defaultValue))
	reflectedDefaultType := reflect.TypeOf(defaultValue)
	if !reflectedDefaultValue.IsValid() {
		return nil, fmt.Errorf("default value is not valid: %+v", defaultValue)
	}
	if reflectedDefaultType.Kind() == reflect.Ptr {
		reflectedDefaultType = reflectedDefaultType.Elem()
	}
	ret := reflect.Indirect(reflect.New(reflectedDefaultType))
	validator := NewMultiObjValidator(validators...)
	fmt.Printf("["+reflectedDefaultType.Name()+"] %s:\n", label)
	for fieldId := 0; fieldId < reflectedDefaultType.NumField(); fieldId++ {
		defaultFieldType := reflectedDefaultType.Field(fieldId)
		defaultFieldValue := reflectedDefaultValue.Field(fieldId)
		retFieldValue := ret.Field(fieldId)
		itemFieldType := NewItemFieldType(ItemFieldText, defaultFieldValue.Interface())
		if itemFieldTypes != nil {
			if v, ok := itemFieldTypes[defaultFieldType.Name]; ok {
				itemFieldType = v
			}
		}
		if ans, err := AskObjectItem("\t"+defaultFieldType.Name, itemFieldType, validator(defaultFieldType.Name)...); err != nil {
			return nil, err
		} else {
			switch defaultFieldType.Type.Kind() {
			case reflect.String:
				retFieldValue.SetString(ans)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
				if i, err := strconv.ParseInt(ans, 10, 32); err != nil {
					return nil, err
				} else {
					retFieldValue.SetInt(i)
				}
			case reflect.Int64:
				if i, err := strconv.ParseInt(ans, 10, 64); err != nil {
					return nil, err
				} else {
					retFieldValue.SetInt(i)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
				if i, err := strconv.ParseUint(ans, 10, 32); err != nil {
					return nil, err
				} else {
					retFieldValue.SetUint(i)
				}
			case reflect.Uint64:
				if i, err := strconv.ParseUint(ans, 10, 64); err != nil {
					return nil, err
				} else {
					retFieldValue.SetUint(i)
				}
			case reflect.Float32, reflect.Float64:
				if f, err := strconv.ParseFloat(ans, 32); err != nil {
					return nil, err
				} else {
					retFieldValue.SetFloat(f)
				}
			case reflect.Bool:
				if b, err := strconv.ParseBool(ans); err != nil {
					return nil, err
				} else {
					retFieldValue.SetBool(b)
				}
			case reflect.Struct:
				return nil, errors.New("Ask on Structs not yet supported")
			case reflect.Map:
				return nil, errors.New("Ask on Maps not yet supported")
			case reflect.Array:
				return nil, errors.New("Ask on Arrays not yet supported")
			case reflect.Interface:
				retFieldValue.Set(reflect.ValueOf(ans))
			}
		}
	}
	return ret.Interface(), nil
}
