package plugins

import "testing"

func TestNewBackend(t *testing.T) {
	b := NewBackend()

	if b.Port != 0 {
		t.Error("New backend should be 0")
	}
}

func TestNewFrontend(t *testing.T) {
	f := NewFrontend()

	if f.Cols != 0 || f.Rows != 0 {
		t.Error("New frontend should get Cols and Rows set to 0")
	}
}

func TestNewPlugin(t *testing.T) {
	p := NewPlugin()

	if p.EltName != "" {
		t.Error("EltName should be equal to \"\"")
	}

	if p.Description != "" {
		t.Error("Description should be equal to \"\"")
	}
}

func TestNewConfiguration(t *testing.T) {
	c := NewConfiguration()

	if len(c.Plugins) != 0 {
		t.Error("New configuration should have a plugins list of 0 element")
	}
}