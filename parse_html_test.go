package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
    <body>
    <a href="/path/one">
        <span>Boot.dev</span>
    </a>
    <a href="https://other.com/path/one">
        <span>Boot.dev</span>
    </a>
    </body>
</html>
            `,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:      "Page without links",
			inputURL:  "https://blog.boot.dev",
			inputBody: "<html><body><H1>Hello World</h1></body></html>",
			expected:  []string{},
		},
		{
			name:     "Nested tags with different registers",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
    <body>
    <a href="/path/one">
        <span>Boot.dev</span>
    </a>
    <H1>
        <B>
            <A href="http://google.com/search">Search</A>
        </B>
    </H1>
    </body>
</html>
            `,
			expected: []string{"https://blog.boot.dev/path/one", "http://google.com/search"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
