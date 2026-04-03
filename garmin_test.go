package main

import (
	"html"
	"testing"
)

func TestExtractGarminURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantURL string
		wantID  string
		wantOK  bool
	}{
		{
			name:    "inreachlink short url",
			input:   "Jag önskar att du var här! inreachlink.com/gRrVf0Zovbj_77GjLmij6wQ   ",
			wantURL: "https://eur.explore.garmin.com/textmessage/viewmsg?extId=gRrVf0Zovbj_77GjLmij6wQ",
			wantID:  "gRrVf0Zovbj_77GjLmij6wQ",
			wantOK:  true,
		},
		{
			name:   "no garmin url",
			input:  "Just a regular text message",
			wantOK: false,
		},
		{
			name:   "empty",
			input:  "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, id, ok := ExtractGarminURL(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok {
				if url != tt.wantURL {
					t.Errorf("url = %q, want %q", url, tt.wantURL)
				}
				if id != tt.wantID {
					t.Errorf("id = %q, want %q", id, tt.wantID)
				}
			}
		})
	}
}

func TestStripGarminURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello! inreachlink.com/abc123   ", "Hello!"},
		{"inreachlink.com/abc123", ""},
		{"No link here", "No link here"},
		{"Before inreachlink.com/xyz After", "Before  After"},
	}
	for _, tt := range tests {
		got := StripGarminURL(tt.input)
		if got != tt.want {
			t.Errorf("StripGarminURL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

const sampleGarminHTML = `
<!DOCTYPE html>
<html>
<head><title>inReach Message from Henrik</title></head>
<body>
<script>
var lat = 59.387261;
var lng = 18.054143;
</script>
</body>
</html>
`

func TestExtractCoords(t *testing.T) {
	// Test regex extraction directly against sample HTML
	m := latRe.FindStringSubmatch(sampleGarminHTML)
	if m == nil {
		t.Fatal("lat regex did not match")
	}
	if m[1] != "59.387261" {
		t.Errorf("lat = %q, want 59.387261", m[1])
	}

	m = lonRe.FindStringSubmatch(sampleGarminHTML)
	if m == nil {
		t.Fatal("lon regex did not match")
	}
	if m[1] != "18.054143" {
		t.Errorf("lon = %q, want 18.054143", m[1])
	}

	m = senderRe.FindStringSubmatch(sampleGarminHTML)
	if m == nil {
		t.Fatal("sender regex did not match")
	}
	if m[1] != "Henrik" {
		t.Errorf("sender = %q, want Henrik", m[1])
	}
}

const sampleGarminHTML2 = `
<html>
<head><title>inReach Message from Anna Sjökvist</title></head>
<body>
<script>
latitude = -33.856159;
longitude = 151.215256;
</script>
</body>
</html>
`

func TestExtractCoordsVariant(t *testing.T) {
	m := latRe.FindStringSubmatch(sampleGarminHTML2)
	if m == nil {
		t.Fatal("lat regex did not match")
	}
	if m[1] != "-33.856159" {
		t.Errorf("lat = %q, want -33.856159", m[1])
	}

	m = lonRe.FindStringSubmatch(sampleGarminHTML2)
	if m == nil {
		t.Fatal("lon regex did not match")
	}
	if m[1] != "151.215256" {
		t.Errorf("lon = %q, want 151.215256", m[1])
	}

	m = senderRe.FindStringSubmatch(sampleGarminHTML2)
	if m[1] != "Anna Sjökvist" {
		t.Errorf("sender = %q, want Anna Sjökvist", m[1])
	}
}

const sampleGarminHTMLEncoded = `
<html>
<head><title>inReach Message from Henrik Sj&#246;kvist</title></head>
<body>
<script>
lat : 59.387273;
lon : 17.847655;
</script>
</body>
</html>
`

func TestExtractCoordsHTMLEncoded(t *testing.T) {
	m := senderRe.FindStringSubmatch(sampleGarminHTMLEncoded)
	if m == nil {
		t.Fatal("sender regex did not match")
	}
	unescaped := html.UnescapeString(m[1])
	if unescaped != "Henrik Sjökvist" {
		t.Errorf("sender = %q, want Henrik Sjökvist", unescaped)
	}
}
