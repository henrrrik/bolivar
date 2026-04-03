package main

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	latRe     = regexp.MustCompile(`\blat(?:itude)?\s*[:=]\s*([-]?\d+\.\d+)`)
	lonRe     = regexp.MustCompile(`\b(?:lng|lon|longitude)\s*[:=]\s*([-]?\d+\.\d+)`)
	senderRe  = regexp.MustCompile(`<title>\s*inReach Message from (.+?)\s*</title>`)
	inreachLinkRe = regexp.MustCompile(`inreachlink\.com/(\w+)`)
)

type GarminResult struct {
	Lat    float64
	Lon    float64
	Sender string
}

func ExtractGarminURL(text string) (url string, extID string, ok bool) {
	m := inreachLinkRe.FindStringSubmatch(text)
	if m == nil {
		return "", "", false
	}
	extID = m[1]
	url = "https://eur.explore.garmin.com/textmessage/viewmsg?extId=" + extID
	return url, extID, true
}

func StripGarminURL(text string) string {
	return strings.TrimSpace(inreachLinkRe.ReplaceAllString(text, ""))
}

var garminClient = &http.Client{Timeout: 15 * time.Second}

func FetchGarminLocation(ctx context.Context, url string) (GarminResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GarminResult{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := garminClient.Do(req)
	if err != nil {
		return GarminResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB cap
	if err != nil {
		return GarminResult{}, err
	}
	htmlBody := string(body)

	var result GarminResult

	if m := latRe.FindStringSubmatch(htmlBody); m != nil {
		result.Lat, _ = strconv.ParseFloat(m[1], 64)
	} else {
		return GarminResult{}, fmt.Errorf("latitude not found in page")
	}

	if m := lonRe.FindStringSubmatch(htmlBody); m != nil {
		result.Lon, _ = strconv.ParseFloat(m[1], 64)
	} else {
		return GarminResult{}, fmt.Errorf("longitude not found in page")
	}

	if m := senderRe.FindStringSubmatch(htmlBody); m != nil {
		result.Sender = html.UnescapeString(m[1])
	}

	return result, nil
}
