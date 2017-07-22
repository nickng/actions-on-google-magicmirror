package backend

import "testing"

func TestValidUser(t *testing.T) {
	if validUser("MagicMirror", "password") {
		t.Errorf("MagicMirror:password should not be a valid user")
	}
	if !validUser("MagicMirror", "rorriMcigaM") {
		t.Errorf("MagicMirror:rorriMcigaM should be a valid user")
	}
}
