package log

import "testing"

func TestLogger(t *testing.T) {
	l := GetLogger("testID")
	l.SetFields()
	l.SetOptions()
	l.Error("test")
	l.Info("test")
}
