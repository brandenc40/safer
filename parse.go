package safer

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	dotSearchParamsRegex   = regexp.MustCompile(`query_string=([0-9]+)`)
	mcs150MileageYearRegex = regexp.MustCompile(`([0-9,]+) \(([0-9]{4})\)`)
)

func parseInt(text string) int {
	if text == "" {
		return 0
	}
	text = strings.Replace(text, ",", "", -1)
	if parsed, err := strconv.Atoi(text); err == nil {
		return parsed
	}
	return 0
}

func parseDate(text string) *time.Time {
	if text == "" || len(text) < 10 {
		return nil
	}
	if parsed, err := time.Parse("01/02/2006", text[:10]); err == nil {
		return &parsed
	}
	return nil
}

func parsePctToFloat32(text string) float32 {
	if text == "" {
		return 0
	}
	if text[len(text)-1] == '%' {
		text = text[:len(text)-1]
	}
	if f, err := strconv.ParseFloat(text, 64); err == nil {
		return float32(f / 100)
	}
	return 0
}

func parseMCS150MileageYear(text string) (mileage int, year string) {
	if text == "" {
		return
	}
	if res := mcs150MileageYearRegex.FindStringSubmatch(text); len(res) == 3 {
		mileage = parseInt(res[1])
		year = res[2]
	}
	return
}

// parse an address returned by xpath query on html.
// multiline address returns an array of strings in this format:
//	[]string{"3101 S PACKERLAND DR", "GREEN BAY, WI \u00a0 54313", "X"}
func parseAddress(texts ...string) string {
	var b strings.Builder
	var written int
	for _, text := range texts {
		if text != "X" {
			if written > 0 {
				b.WriteString(" ")
			}
			b.WriteString(strings.ReplaceAll(text, "\u00a0 ", "")) // remove &nbsp;
			written++
		}
	}
	return b.String()
}

func parseDotFromSearchParams(params string) string {
	if params == "" {
		return ""
	}
	if res := dotSearchParamsRegex.FindStringSubmatch(params); len(res) == 2 {
		return res[1]
	}
	return ""
}
