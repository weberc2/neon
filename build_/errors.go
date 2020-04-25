package build

import (
	"fmt"
	"log"
)

type Error struct {
	UserText string
	Context  []interface{}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s: %v", err.UserText, err.Context)
}

type LogFunc func(v ...interface{})

type ErrorClass struct {
	ID      string
	LogFunc LogFunc
}

func (c ErrorClass) New(id, userText string, v ...interface{}) Error {
	id = c.ID + "." + id
	if c.LogFunc != nil {
		c.LogFunc(append([]interface{}{id}, v...)...)
	}
	return Error{
		UserText: userText,
		Context:  v,
	}
}

var DefaultLogFunc = log.Println

func NewErrorClass(id string) ErrorClass {
	return ErrorClass{ID: id, LogFunc: DefaultLogFunc}
}

func Err(id string, userText string, v ...interface{}) Error {
	return NewErrorClass("").New(id, userText, v...)
}
