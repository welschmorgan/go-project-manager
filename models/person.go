package models

type Person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func NewPerson(name, email, phone string) *Person {
	return &Person{
		Name:  name,
		Email: email,
		Phone: phone,
	}
}
