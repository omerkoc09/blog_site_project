package utils

import "testing"

func TestToTitle(t *testing.T) {

	var tests = []struct {
		input string
		want  string
	}{
		{"talha", "Talha"},
		{"tAlha karaca", "Talha Karaca"},
		{"tAlha karaca    ", "Talha Karaca"},
		{"tAlha    karaca    ", "Talha Karaca"},
		{"risale-i Nur", "Risale-i Nur"},
		{"ihlas", "İhlas"},
		{"İŞÇİ", "İşçi"},
	}

	for _, test := range tests {
		if got := ToTitle(test.input); got != test.want {
			t.Errorf("ToTitle(%q) = %q; want %q", test.input, got, test.want)
		}
	}

}
