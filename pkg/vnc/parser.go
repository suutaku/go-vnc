package vnc

import (
	"reflect"

	"github.com/suutaku/go-vnc/internal/auth"
	"github.com/suutaku/go-vnc/internal/encodings"
	"github.com/suutaku/go-vnc/internal/events"
)

func addAuthType(name string) *auth.Type {
	for _, present := range auth.DefaultAuthTypes {
		if reflect.TypeOf(present).Elem().Name() == name {
			return &present
		}
	}
	return nil
}

func addEncoding(name string) *encodings.Encoding {
	for _, present := range encodings.DefaultEncodings {
		if reflect.TypeOf(present).Elem().Name() == name {
			return &present
		}
	}
	return nil
}

func addEvent(name string) *events.Event {
	for _, present := range events.DefaultEvents {
		if reflect.TypeOf(present).Elem().Name() == name {
			return &present
		}
	}
	return nil
}

func configureAuthTypes(args []string) []auth.Type {
	validType := make([]auth.Type, 0)
	for _, arg := range args {
		tmp := addAuthType(arg)
		if tmp != nil {
			validType = append(validType, *tmp)
		}
	}
	return validType
}

func configureEncodings(args []string) []encodings.Encoding {
	validType := make([]encodings.Encoding, 0)
	for _, arg := range args {
		tmp := addEncoding(arg)
		if tmp != nil {
			validType = append(validType, *tmp)
		}
	}
	return validType
}

func configureEvents(args []string) []events.Event {
	validType := make([]events.Event, 0)
	for _, arg := range args {
		tmp := addEvent(arg)
		if tmp != nil {
			validType = append(validType, *tmp)
		}
	}
	return validType
}

func authIsEnabled(tt []auth.Type, name string) bool {
	for _, t := range tt {
		if reflect.TypeOf(t).Elem().Name() == name {
			return true
		}
	}
	return false
}
