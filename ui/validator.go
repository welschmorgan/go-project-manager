package ui

type Validator func(string) error
type ObjValidator func(k, v string) error

func NewMultiValidator(validators ...Validator) func(string) error {
	return func(s string) error {
		var err error
		for _, v := range validators {
			if v != nil {
				if err = v(s); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func NewMultiObjValidator(validators ...ObjValidator) func(k string) []Validator {
	return func(k string) []Validator {
		ret := []Validator{}
		for _, validator := range validators {
			if validator != nil {
				ret = append(ret, func(v string) error {
					return validator(k, v)
				})
			}
		}
		return ret
	}
}
