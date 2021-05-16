package ui

type Validator func(string) error
type ObjValidator func(k, v string) error
