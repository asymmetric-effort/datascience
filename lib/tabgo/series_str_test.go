//go:build unit

package tabgo

import (
	"testing"
)

func TestStrLower(t *testing.T) {
	s := NewSeries("x", []any{"HELLO", "World", nil, "GoLang"})
	result := s.Str().Lower()
	vals := result.Values()
	if vals[0] != "hello" || vals[1] != "world" || vals[2] != nil || vals[3] != "golang" {
		t.Errorf("Lower() = %v", vals)
	}
}

func TestStrUpper(t *testing.T) {
	s := NewSeries("x", []any{"hello", "World", nil, "GoLang"})
	result := s.Str().Upper()
	vals := result.Values()
	if vals[0] != "HELLO" || vals[1] != "WORLD" || vals[2] != nil || vals[3] != "GOLANG" {
		t.Errorf("Upper() = %v", vals)
	}
}

func TestStrContains(t *testing.T) {
	s := NewSeries("x", []any{"hello world", "goodbye", nil, "hello there"})
	result := s.Str().Contains("hello")
	if !result[0] || result[1] || result[2] || !result[3] {
		t.Errorf("Contains('hello') = %v", result)
	}
}

func TestStrReplace(t *testing.T) {
	s := NewSeries("x", []any{"hello world", "hello there", nil})
	result := s.Str().Replace("hello", "hi")
	vals := result.Values()
	if vals[0] != "hi world" || vals[1] != "hi there" || vals[2] != nil {
		t.Errorf("Replace() = %v", vals)
	}
}

func TestStrSplit(t *testing.T) {
	s := NewSeries("x", []any{"a,b,c", "d,e", nil})
	result := s.Str().Split(",")
	vals := result.Values()
	parts, ok := vals[0].([]string)
	if !ok || len(parts) != 3 || parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
		t.Errorf("Split(',')[0] = %v", vals[0])
	}
	if vals[2] != nil {
		t.Errorf("Split(',')[2] = %v, want nil", vals[2])
	}
}

func TestStrStrip(t *testing.T) {
	s := NewSeries("x", []any{"  hello  ", "\tworld\n", nil})
	result := s.Str().Strip()
	vals := result.Values()
	if vals[0] != "hello" || vals[1] != "world" || vals[2] != nil {
		t.Errorf("Strip() = %v", vals)
	}
}

func TestStrLen(t *testing.T) {
	s := NewSeries("x", []any{"hello", "ab", nil, ""})
	result := s.Str().Len()
	vals := result.Values()
	if vals[0] != 5 || vals[1] != 2 || vals[2] != nil || vals[3] != 0 {
		t.Errorf("Len() = %v", vals)
	}
}

func TestStrStartsWith(t *testing.T) {
	s := NewSeries("x", []any{"hello world", "goodbye", nil, "help"})
	result := s.Str().StartsWith("hel")
	if !result[0] || result[1] || result[2] || !result[3] {
		t.Errorf("StartsWith('hel') = %v", result)
	}
}

func TestStrEndsWith(t *testing.T) {
	s := NewSeries("x", []any{"hello world", "goodbye world", nil, "help"})
	result := s.Str().EndsWith("world")
	if !result[0] || !result[1] || result[2] || result[3] {
		t.Errorf("EndsWith('world') = %v", result)
	}
}

func TestStrSlice(t *testing.T) {
	s := NewSeries("x", []any{"hello", "hi", nil, "world"})
	result := s.Str().Slice(0, 3)
	vals := result.Values()
	if vals[0] != "hel" || vals[1] != "hi" || vals[2] != nil || vals[3] != "wor" {
		t.Errorf("Slice(0,3) = %v", vals)
	}
}
