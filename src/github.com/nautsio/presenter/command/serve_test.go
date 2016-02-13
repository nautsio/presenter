package command

import "testing"

func TestPathIsAbsolute(t *testing.T) {
	var cases = []struct {
		path     string
		absolute bool
	}{
		{"/absolute/path", true},
		{"relative/path", false},
		{".", false},
		{"", false},
	}

	for _, c := range cases {
		v := pathIsAbsolute(c.path)
		if v != c.absolute {
			t.Error(
				"For", c.path,
				"expected", c.absolute,
				"got", v,
			)
		}
	}
}
