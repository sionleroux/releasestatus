package main

import (
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	const author = "Bob"
	var err error
	cur := Release{}
	err = cur.start(author)

	// Release starts from default state
	if err != nil {
		t.Errorf("Expected successful release start, got: %v", err)
	}

	// Release marked as running
	if cur.Running != true {
		t.Errorf("Expected release to be running (true), got: %v", cur.Running)
	}

	// Consective start attempts fail
	err = cur.start("Alice")
	if err == nil {
		t.Errorf("Expected failed release start, got: %v", err)
	}

	// Release author is still from the first release
	if cur.Author != author {
		t.Errorf("Expected author to be %s, got: %s", author, cur.Author)
	}
}

func TestStop(t *testing.T) {
	var err error
	cur := Release{"Bob", true, false, time.Now()}
	err = cur.stop()

	// Release stops from running state
	if err != nil {
		t.Errorf("Expected successful release stop, got: %v:", err)
	}

	// Release marked not running
	if cur.Running != false {
		t.Errorf("Expected release not to be running (false), got: %v", cur.Running)
	}

	// Consective stop attempts fail
	err = cur.stop()
	if err == nil {
		t.Errorf("Expected failed release stop, got: %v:", err)
	}
}

func TestBlock(t *testing.T) {
	cur := Release{"Alice", false, true, time.Now()}
	err := cur.start("Bob")

	if err == nil {
		t.Errorf("Expected failed release start, got %v", err)
	}

	if err.Error() != "release blocked" {
		t.Errorf("Expected release to be blocked, got other error: %v", err)
	}
}
