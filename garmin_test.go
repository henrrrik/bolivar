package main

import (
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
			name:    "eur subdomain",
			input:   "Hello from the trail! https://eur.explore.garmin.com/textmessage/viewmsg?extId=gDCYYWDrMkokPJUKOOWO8ng",
			wantURL: "https://eur.explore.garmin.com/textmessage/viewmsg?extId=gDCYYWDrMkokPJUKOOWO8ng",
			wantID:  "gDCYYWDrMkokPJUKOOWO8ng",
			wantOK:  true,
		},
		{
			name:    "no subdomain",
			input:   "Check this: https://explore.garmin.com/textmessage/viewmsg?extId=abc123",
			wantURL: "https://explore.garmin.com/textmessage/viewmsg?extId=abc123",
			wantID:  "abc123",
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
