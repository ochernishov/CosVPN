package transport

import (
	"testing"

	"github.com/ochernishov/cosvpn/obfs"
)

func TestAutoTransportDefault(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "auto"}
	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "auto" {
		t.Errorf("expected auto, got %s", tr.CurrentMode())
	}
}

func TestDirectMode(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "direct"}
	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "direct" {
		t.Errorf("expected direct, got %s", tr.CurrentMode())
	}
	if tr.NeedsTLS() {
		t.Error("direct mode should not need TLS")
	}
}

func TestTLSMode(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "tls"}
	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "tls" {
		t.Errorf("expected tls, got %s", tr.CurrentMode())
	}
	if !tr.NeedsTLS() {
		t.Error("tls mode should need TLS")
	}
}

func TestEmptyModeDefaultsToAuto(t *testing.T) {
	config := obfs.ObfsConfig{Mode: ""}
	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "auto" {
		t.Errorf("expected auto for empty mode, got %s", tr.CurrentMode())
	}
}

func TestSwitchToTLS(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "auto"}
	tr := NewAutoTransport(config)
	tr.SwitchToTLS()
	if tr.CurrentMode() != "tls" {
		t.Errorf("expected tls after switch, got %s", tr.CurrentMode())
	}
}

func TestSetMode(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "auto"}
	tr := NewAutoTransport(config)
	tr.SetMode("direct")
	if tr.CurrentMode() != "direct" {
		t.Errorf("expected direct, got %s", tr.CurrentMode())
	}
}
