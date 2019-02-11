package main

import "testing"

// TestDefaultRoute ..
func TestDefaultRoute(t *testing.T) {
	resp := DefaultRoute(make(map[string]string))
	if resp.S.Code != 200 {
		t.Error("Default route doesn't return HTTP status 200.")
	}
}

// TestCreateUser ..
func TestCreateUser(t *testing.T) {

}
