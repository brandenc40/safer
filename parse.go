package safer

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	dotSearchParamsRegex = regexp.MustCompile(`query\.asp\?searchtype=ANY&query_type=queryCarrierSnapshot&query_param=USDOT&original_query_param=NAME&query_string=([0-9]+)&original_query_string=.+`)
)

func parseInt(text string) int {
	if text == "" {
		return 0
	}
	if parsed, err := strconv.Atoi(text); err == nil {
		return parsed
	}
	text = strings.Replace(text, ",", "", -1)
	if parsed, err := strconv.Atoi(text); err == nil {
		return parsed
	}
	return 0
}

func parseDate(text string) *time.Time {
	if text == "" {
		return nil
	}
	if parsed, err := time.Parse("01/02/2006", text); err == nil {
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
	if s := strings.Split(text, " ("); len(s) == 2 {
		mileage = parseInt(s[0])
		year = s[1][:len(s[1])-1]
	}
	return
}

func parseAddress(texts ...string) (fullAddr string) {
	for i, text := range texts {
		if text != "X" {
			if i > 0 {
				fullAddr += " "
			}
			fullAddr += strings.ReplaceAll(text, "\u00a0 ", "") // remove &nbsp;
		}
	}
	return
}

func parseDotFromSearchParams(params string) string {
	if res := dotSearchParamsRegex.FindStringSubmatch(params); len(res) == 2 {
		return res[1]
	}
	return ""
}
