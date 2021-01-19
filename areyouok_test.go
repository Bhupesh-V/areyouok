package main

import (
    "testing"
)

func TestInA(t *testing.T) {
    ans := in("test", []string{"test", "sample"})
    if ans != true {
        t.Errorf("In() want %s, want %s", "false", "true")
    }
}

func TestInB(t *testing.T) {
    ans := in("nice", []string{"test", "sample"})
    if ans != false {
        t.Errorf("In() want %s, want %s", "true", "false")
    }
}
