package ddosify

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLatencyChecker_RunCommandExec(t *testing.T) {
	var input LatencyAPIRequest
	var outputTokenOK = tokenAPIResponse{
		RequestCount: 100000,
		Duration:     100,
	}
	var outputTokenKOInsufficient = tokenAPIResponse{
		RequestCount: 1,
		Duration:     100,
	}

	var outputLatencyOK = LatencyCheckerOutputList{
		Result: []LatencyCheckerOutput{
			{
				Location:   "us-east-1",
				AvgLatency: 200,
			},
		},
	}

	loc := []string{"us-east-1"}
	locLocationDown := []string{"LocationDown"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.WriteHeader(200)
			outputTok, _ := json.Marshal(outputTokenOK)
			w.Write(outputTok)
			return
		}
		if r.URL.Path == "/tokenInssuficient" {
			w.WriteHeader(200)
			output, _ := json.Marshal(outputTokenKOInsufficient)
			w.Write(output)
			return
		}
		if r.URL.Path == "/tokenBadStatusCode" {
			w.WriteHeader(400)
			output, _ := json.Marshal(outputTokenKOInsufficient)
			w.Write(output)
			return
		}
		if r.URL.Path == "/tokenBadJsonResponse" {
			w.WriteHeader(400)
			output, _ := json.Marshal("{outputTokenKOInsufficient}")
			w.Write(output)
			return
		}

		json.NewDecoder(r.Body).Decode(&input)
		if input.TargetURL == "https://requestOK.test" {

			if input.Locations[0] == loc[0] {
				w.WriteHeader(200)
				jsonResp := `{"us-east-1":{"latency":200,"status_code":200}}`
				w.Write([]byte(jsonResp))
				return
			}
			if input.Locations[0] == locLocationDown[0] {
				w.WriteHeader(200)
				jsonResp := `{"us-east-1":{"latency":200,"status_code":400}}`
				w.Write([]byte(jsonResp))
				return
			}
			w.Write([]byte("error"))
			w.WriteHeader(400)
			return

		}

	},
	))

	defer srv.Close()

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
	tests := []struct {
		name    string
		fields  fields
		want    LatencyCheckerOutputList
		wantErr bool
	}{
		{
			name: "Test OK executing the funtion with no errors",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"us-east-1", "us-west-1"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/token",
				ServiceAPIURL:         srv.URL,
			},
			want:    outputLatencyOK,
			wantErr: false,
		},
		{
			name: "Test OK executing the funtion with no errors, multiples runs",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  2,
				WaitInterval:          1,
				Locations:             []string{"us-east-1", "us-west-1"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/token",
				ServiceAPIURL:         srv.URL,
			},
			want:    outputLatencyOK,
			wantErr: false,
		},
		{
			name: "Test KO to test errors in availableTokens",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"us-east-1", "us-west-1"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/tokenInssuficient",
				ServiceAPIURL:         srv.URL,
			},
			want:    LatencyCheckerOutputList{},
			wantErr: true,
		},
		{
			name: "Test KO to test errors in APIKEY",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"us-east-1", "us-west-1"},
				APIKey:                "NOT_SET",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/token",
				ServiceAPIURL:         srv.URL,
			},
			want:    LatencyCheckerOutputList{},
			wantErr: true,
		},
		{
			name: "Test KO to test errors in StatusCode Response",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"us-east-1", "us-west-1"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/tokenBadStatusCode",
				ServiceAPIURL:         srv.URL,
			},
			want:    LatencyCheckerOutputList{},
			wantErr: true,
		},
		{
			name: "Test KO to test errors in Json unmarshal decode",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"us-east-1", "us-west-1"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/tokenBadJsonResponse",
				ServiceAPIURL:         srv.URL,
			},
			want:    LatencyCheckerOutputList{},
			wantErr: true,
		},
		{
			name: "Test KO to test errors with Bad location",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"badlocation"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/token",
				ServiceAPIURL:         srv.URL,
			},
			want:    LatencyCheckerOutputList{},
			wantErr: true,
		},
		/*{
			name: "Test KO to test errors with Down Location (status code != 200)",
			fields: fields{
				TargetUrl:             "https://requestOK.test",
				Runs:                  1,
				WaitInterval:          1,
				Locations:             []string{"LocationDown"},
				APIKey:                "APIKEY",
				ContentType:           CONTENT_TYPE_REQ,
				OutputLocationsNumber: 1,
				ServiceAPITokenURL:    srv.URL + "/token",
				ServiceAPIURL:         srv.URL,
			},
			want:    outputLatencyOK,
			wantErr: true,
		},*/
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
			log.Println(lc.Locations)
			got, err := lc.RunCommandExec()
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommandExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunCommandExec() got = %v, want %v", got, tt.want)
			}
		})
	}
}
