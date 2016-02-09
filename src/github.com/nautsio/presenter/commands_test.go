package main

import "testing"

func TestDoInit(t *testing.T) {

}

func TestCreateExamplePresentationPath(t *testing.T) {

}

func TestCreateExampleSlides(t *testing.T) {

}

func TestCreateExampleImageDirectory(t *testing.T) {

}

func TestCreateExampleTheme(t *testing.T) {

}

func TestThemeExists(t *testing.T) {

}

func TestFileExists(t *testing.T) {

}

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

func TestDoServe(t *testing.T) {

}
