package ddosify

import (
	"os"
	"reflect"
	"testing"
)

func TestIntervalTimeToSeconds(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test OK testing minutes",
			args: args{
				"1m",
			},
			want: 60,
		},
		{
			name: "Test KO testing wrong conversion type",
			args: args{
				"testKO",
			},
			want: -1,
		},
		{
			name: "Test KO wrong duration type",
			args: args{
				"1q",
			},
			want: -1,
		},
		{
			name: "Test OK testing seconds",
			args: args{
				"4s",
			},
			want: 4,
		},
		{
			name: "Test OK testing hours",
			args: args{
				"1h",
			},
			want: 3600,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntervalTimeToSeconds(tt.args.interval); got != tt.want {
				t.Errorf("IntervalTimeToSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatencyChecker_getMinimumLatencies(t *testing.T) {
	type fields struct {
		TargetUrl             string
		Runs                  int
		WaitInterval          int
		Locations             []string
		APIKey                string
		ContentType           string
		OutputLocationsNumber int
		ServiceAPITokenURL    string
		ServiceAPIURL         string
	}
	type args struct {
		latencies map[string]float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
		want1  []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LatencyChecker{
				TargetUrl:             tt.fields.TargetUrl,
				Runs:                  tt.fields.Runs,
				WaitInterval:          tt.fields.WaitInterval,
				Locations:             tt.fields.Locations,
				APIKey:                tt.fields.APIKey,
				ContentType:           tt.fields.ContentType,
				OutputLocationsNumber: tt.fields.OutputLocationsNumber,
				ServiceAPITokenURL:    tt.fields.ServiceAPITokenURL,
				ServiceAPIURL:         tt.fields.ServiceAPIURL,
			}
			got, got1 := lc.getMinimumLatencies(tt.args.latencies)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMinimumLatencies() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getMinimumLatencies() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateIntervalTime(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test OK validation interval time seconds",
			args: args{
				interval: "1s",
			},
			want: true,
		},
		{
			name: "Test OK validation interval time minutes",
			args: args{
				interval: "1m",
			},
			want: true,
		},
		{
			name: "Test OK validation interval time hours ",
			args: args{
				interval: "1h",
			},
			want: true,
		},
		{
			name: "Test KO validation interval time wrong type",
			args: args{
				interval: "1q",
			},
			want: false,
		},
		{
			name: "Test KO validation interval time empty",
			args: args{
				interval: "",
			},
			want: false,
		},
		{
			name: "Test KO validation interval time wrong value",
			args: args{
				interval: "kk",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateIntervalTime(tt.args.interval); got != tt.want {
				t.Errorf("ValidateIntervalTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	type args struct {
		inputURL string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test OK Valid URL",
			args: args{
				inputURL: "https://www.google.com",
			},
			want: true,
		},
		{
			name: "Test KO wrong URL no url scheme",
			args: args{
				inputURL: "www.google.com",
			},
			want: false,
		},
		{
			name: "Test KO wrong URL bad url scheme",
			args: args{
				inputURL: "http//www.google.com",
			},
			want: false,
		},
		{
			name: "Test KO wrong URL no host",
			args: args{
				inputURL: "http://",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateURL(tt.args.inputURL); got != tt.want {
				t.Errorf("ValidateURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnv(t *testing.T) {
	os.Setenv("TEST_OK", "myvalue")
	type args struct {
		varName      string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test OK env variable exists",
			args: args{
				varName:      "TEST_OK",
				defaultValue: "mydefaultvalue",
			},
			want: "myvalue",
		},
		{
			name: "Test OK env variable not exists, return default",
			args: args{
				varName:      "TEST_DEFAULT",
				defaultValue: "mydefaultvalue",
			},
			want: "mydefaultvalue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.varName, tt.args.defaultValue); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
