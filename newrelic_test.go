package go_newrelic_api

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "fmt"
    "io/ioutil"
)

func TestNewNewrelic(t *testing.T) {
    key := "1234"
    nr := NewNewrelic(key)
    if nr.Key != "1234" {
        t.Errorf("Key was expected to be %s, got %s", key, nr.Key)
    }
    
    baseurl := "https://api.newrelic.com/v2"
    if nr.BaseUrl != baseurl {
        t.Errorf("BaseUrl was expected to be %s, got %s", baseurl, nr.BaseUrl)
    }

    format := "json"
    if nr.Format != format {
        t.Errorf("Format was expected to be %s, got %s", format, nr.Format)
    }
}

func TestGetApplications(t *testing.T) {
    json_out, _ := ioutil.ReadFile("fixtures/get_applications_test.json")

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        expected_url := "/applications.json"
        if r.URL.Path != expected_url {
            t.Errorf("URL was wrong. Expected: %s, got: %s", expected_url, r.URL.Path)
        }

        w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(json_out))
	}))
	defer ts.Close()

    nr := Newrelic{"1234", ts.URL, "json"}

    out := nr.GetApplications()

    if len(out.Applications) != 2 {
        t.Errorf("Length was wrong. Expected: %d, got: %d", 2, len(out.Applications))
    }

    if out.Applications[1].Id != 456 {
        t.Errorf("Expected ID for the second application was: %d, got: %d", 456, out.Applications[1].Id)
    }
}

func TestGetApplication(t *testing.T) {
    json_out, _ := ioutil.ReadFile("fixtures/get_application_test.json")

    app_id := 5678

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        expected_url := fmt.Sprintf("/applications/%d.json", app_id)

        if r.URL.Path != expected_url {
            t.Errorf("URL was wrong. Expected: %s, got: %s", expected_url, r.URL.Path)
        }

        w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(json_out))
	}))
	defer ts.Close()

    nr := Newrelic{"1234", ts.URL, "json"}

    out := nr.GetApplication(app_id)

    if out.Application.Id != 123 {
        t.Errorf("Expected ID for the second application was: %d, got: %d", 123, out.Application.Id)
    }
}

func TestGetMetricNames(t *testing.T) {
    json_out, _ := ioutil.ReadFile("fixtures/get_metric_names.json")

    app_id := 5678

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        expected_url := fmt.Sprintf("/applications/%d/metrics.json", app_id)

        if r.URL.Path != expected_url {
            t.Errorf("URL was wrong. Expected: %s, got: %s", expected_url, r.URL.Path)
        }

        w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(json_out))
	}))
	defer ts.Close()

    nr := Newrelic{"1234", ts.URL, "json"}

    out := nr.GetMetricNames(app_id)

    expected_name := "Agent/MetricsReported/count"
    if out.Metrics[0].Name !=  expected_name {
        t.Errorf("Expected name for the second metric was: %s, got: %s", expected_name, out.Metrics[0].Name)
    }
}
