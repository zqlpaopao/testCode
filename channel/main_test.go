package main

import "testing"

func Test_worker(t *testing.T) {
	type args struct {
		id   int
		ch   chan token
		next chan token
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
