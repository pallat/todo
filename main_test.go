package main

import (
	"reflect"
	"testing"
)

func TestNewTodo(t *testing.T) {
	todoList = []todo{}

	NewTodo("Learn golang")

	if todoList[0].Topic != "Learn golang" {
		t.Error("it should store Learn golang but get", todoList[0].Topic)
	}
}

func TestListTodo(t *testing.T) {
	todoList = []todo{}

	NewTodo("Learn golang")
	NewTodo("Learn git")
	NewTodo("Learn VueJS")

	list := ListTodo()

	for i := range list {
		list[i].ID = ""
	}

	expected := []todo{
		{Topic: "Learn golang"},
		{Topic: "Learn git"},
		{Topic: "Learn VueJS"},
	}

	if !reflect.DeepEqual(list, expected) {
		t.Errorf("%v\nis expected but get\n%v\n", expected, list)
	}
}
