package formsubmission

import "testing"

func TestNormalizeTextareaSubmission(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "trims leading and trailing whitespace",
			input: "\n  hello world  \n",
			want:  "hello world",
		},
		{
			name:  "preserves interior newlines",
			input: "first line\nsecond line\nthird line",
			want:  "first line\nsecond line\nthird line",
		},
		{
			name:  "normalizes windows newlines",
			input: "\r\nfirst line\r\nsecond line\r\n",
			want:  "first line\nsecond line",
		},
		{
			name:  "whitespace only becomes empty",
			input: " \n\t\r\n ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeTextareaSubmission(tt.input)
			if got != tt.want {
				t.Fatalf("normalizeTextareaSubmission(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
