package appError

import "fmt"

type Errors []*Element

type Element struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func New(field, message string) *Element {
	return &Element{Field: field, Message: message}
}

func (e Element) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
