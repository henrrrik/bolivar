package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

var (
	latRe     = regexp.MustCompile(`\blat(?:itude)?\s*[:=]\s*([-]?\d+\.\d+)`)
	lonRe     = regexp.MustCompile(`\b(?:lng|lon|longitude)\s*[:=]\s*([-]?\d+\.\d+)`)
	senderRe  = regexp.MustCompile(`<title>\s*inReach Message from (.+?)\s*</title>`)
	garminURL = regexp.MustCompile(`https?://(?:eur\.)?explore\.garmin\.com/textmessage/viewmsg\?extId=(\w+)`)
)

type GarminResult struct {
	Lat    float64
	Lon    float64
	Sender string
}

func ExtractGarminURL(text string) (url string, extID string, ok bool) {
	m := garminURL.FindStringSubmatch(text)
	if m == nil {
		return "", "", false
	}
	return m[0], m[1], true
}

func FetchGarminLocation(ctx context.Context, url string) (GarminResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GarminResult{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GarminResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB cap
	if err != nil {
		return GarminResult{}, err
	}
	html := string(body)

	var result GarminResult

	if m := latRe.FindStringSubmatch(html); m != nil {
		result.Lat, _ = strconv.ParseFloat(m[1], 64)
	} else {
		return GarminResult{}, fmt.Errorf("latitude not found in page")
	}

	if m := lonRe.FindStringSubmatch(html); m != nil {
		result.Lon, _ = strconv.ParseFloat(m[1], 64)
	} else {
		return GarminResult{}, fmt.Errorf("longitude not found in page")
	}

	if m := senderRe.FindStringSubmatch(html); m != nil {
		result.Sender = m[1]
	}

	return result, nil
}
