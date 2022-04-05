package auth

import (
	"github.com/suutaku/go-vnc/internal/buffer"
)

// Type represents an authentication type.
type Type interface {
	Code() uint8
	Negotiate(wr *buffer.ReadWriter) error
	Response(wr *buffer.ReadWriter) error
}

// DefaultAuthTypes is the default enabled list of auth types.
var DefaultAuthTypes = []Type{
	&None{},
	&VNCAuth{},
	&TightSecurity{},
}

// GetDefaults returns a slice of the default auth handlers.
func GetDefaults() []Type {
	out := make([]Type, len(DefaultAuthTypes))
	copy(out, DefaultAuthTypes)
	return out
}
