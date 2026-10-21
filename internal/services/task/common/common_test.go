package common

import (
	"reflect"
	"testing"
)

// NewService must wire every dependency. A missing urlService made Insert panic
// with a nil pointer dereference for UrlScan targets (#342).
func TestNewServiceInitializesAllDependencies(t *testing.T) {
	s, ok := NewService().(*service)
	if !ok {
		t.Fatalf("NewService() returned %T, want *service", NewService())
	}

	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		if field := v.Field(i); field.Kind() == reflect.Interface && field.IsNil() {
			t.Errorf("NewService() left %s nil", v.Type().Field(i).Name)
		}
	}
}
