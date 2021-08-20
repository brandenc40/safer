package safer

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseInt(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "without comma",
			args: args{"10885"},
			want: 10_885,
		},
		{
			name: "with comma",
			args: args{"10,884"},
			want: 10_884,
		},
		{
			name: "with multi comma",
			args: args{"100,123,884"},
			want: 100_123_884,
		},
		{
			name: "empty string",
			args: args{""},
			want: 0,
		},
		{
			name: "unable to parse",
			args: args{"$1,290"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.args.text); got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDate(t *testing.T) {
	type args struct {
		text string
	}
	date := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		args args
		want *time.Time
	}{
		{
			name: "04/05/2020",
			args: args{"04/05/2020"},
			want: &date,
		},
		{
			name: "empty",
			args: args{""},
			want: nil,
		},
		{
			name: "err",
			args: args{"2020-01-01"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDate(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePctToFloat32(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "with %",
			args: args{"1.4%"},
			want: 0.014,
		},
		{
			name: "without %",
			args: args{"3.45"},
			want: 0.0345,
		},
		{
			name: "empty",
			args: args{""},
			want: 0.0,
		},
		{
			name: "err",
			args: args{"N/A"},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePctToFloat32(tt.args.text); got != tt.want {
				t.Errorf("parsePctToFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMCS150MileageYear(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name        string
		args        args
		wantMileage int
		wantYear    string
	}{
		{
			name:        "expected",
			args:        args{"1,100,158,928 (2020)"},
			wantMileage: 1_100_158_928,
			wantYear:    "2020",
		},
		{
			name:        "invalid",
			args:        args{"1,100,158,928(2020)"},
			wantMileage: 0,
			wantYear:    "",
		},
		{
			name:        "empty",
			args:        args{""},
			wantMileage: 0,
			wantYear:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMileage, gotYear := parseMCS150MileageYear(tt.args.text)
			if gotMileage != tt.wantMileage {
				t.Errorf("parseMCS150MileageYear() gotMileage = %v, want %v", gotMileage, tt.wantMileage)
			}
			if gotYear != tt.wantYear {
				t.Errorf("parseMCS150MileageYear() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
		})
	}
}

func Test_parseAddress(t *testing.T) {
	type args struct {
		texts []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "expected",
			args: args{[]string{"3101 S PACKERLAND DR", "GREEN BAY, WI \u00a0 54313", "X"}},
			want: "3101 S PACKERLAND DR GREEN BAY, WI 54313",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAddress(tt.args.texts...); got != tt.want {
				t.Errorf("parseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDotFromSearchParams(t *testing.T) {
	type args struct {
		params string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "expected",
			args: args{"query.asp?searchtype=ANY&query_type=queryCarrierSnapshot&query_param=USDOT&original_query_param=NAME&query_string=2819773&original_query_string=DALE%20PFEFFER"},
			want: "2819773",
		},
		{
			name: "wrong format",
			args: args{"hello"},
			want: "",
		},
		{
			name: "empty string",
			args: args{""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDotFromSearchParams(tt.args.params); got != tt.want {
				t.Errorf("parseDotFromSearchParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
