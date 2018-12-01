package main

import "testing"

func BenchmarkToTestRaceCondition(b *testing.B) {
	go func() {
		for i := 0; i < b.N; i++ {
			NewTodo("x")
		}
	}()

	for j := 0; j < b.N; j++ {
		ListTodo()
	}
}
