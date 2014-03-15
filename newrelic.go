package newrelic_api

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "net/url"
)

// Newrelic is an object that stores various critical access settings for Newrelic including
// the api key, baseurl for requests and the format of response requested.
type Newrelic struct {
    Key string
    BaseUrl string
    Format string
}

// NewrelicApplications is used to represent the response format from the list applications (/applications) call.
// It aims to encode the response precisely.
type NewrelicApplications struct {
    Applications []struct {
        ID int `json:"id"`
        Name string `json:"name"`
        Language string `json:"language"`
        HealthStatus string `json:"health_status"`
        Reporting bool `json:"reporting"`
        LastReportedAt string `json:"last_reported_at"`
        ApplicationSummary struct {
            ResponseTime float32 `json:"response_time"`
            Throughput float32 `json:"throughput"`
            ErrorRate float32 `json:"error_rate"`
            ApdexTarget float32 `json:"apdex_target"`
            ApdexScore float32 `json:"apdex_score"`
        } `json:"application_summary"`
        EndUserSummary struct {
            ResponseTime float32 `json:"response_time"`
            Throughput float32 `json:"throughput"`
            ApdexTarget float32 `json:"apdex_target"`
            ApdexScore float32 `json:"apdex_score"`
        } `json:"end_user_summary"`
        Settings struct {
            AppApdexThreshold float32 `json:"app_apdex_threshold"`
            EndUserApdexThreshold float32 `json:"end_user_apdex_threshold"`
            EnableRealUserMonitoring bool `json:"enable_real_user_monitoring"`
            UseServerSideConfig bool `json:"use_server_side_config"`
        } `json:"settings"`
        Links struct {
            Servers []int `json:"servers"`
            ApplicationHosts []int `json:"application_hosts"`
            ApplicationInstances []int `json:"application_instances"`
        } `json:"links"`
    } `json:"applications"`
}

// NewrelicMetricData is used to represent the response from the call to display metrics data for a given application (/applications/{application_id}/metrics/data).
type NewrelicMetricData struct {
    MetricData struct {
        From string `json:"from"`
        To string `json:"to"`
        Metrics []struct {
            Name string `json:"name"`
            Timeslices []struct {
                From string `json:"from"`
                To string `json:"to"`
                Values map[string]interface{} `json:"values"`
            } `json:"timeslices"`
        } `json:"metrics"`
    } `json:"metric_data"`
}

// NewNewrelic returns a *Newrelic pointer that can be used to invoke various API endpoints.
// The 'key' argument must be your newrelic api key. The other settings are set as defaults :
// BaseUrl = "https://api.newrelic.com/v2" and Format="json". These should not be changed else
// methods downstream will not work (specically the Format).
func NewNewrelic(key string) *Newrelic {
    nr := new(Newrelic)
    nr.Key = key
    nr.BaseUrl = "https://api.newrelic.com/v2"
    nr.Format = "json"

    return nr
}

func (nr *Newrelic) makeRequest(url string) ([]byte, error) {
    url = fmt.Sprintf("%s/%s.%s", nr.BaseUrl, url, nr.Format)

    client := http.Client{}
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("X-Api-Key", nr.Key)
    resp, _ := client.Do(req)

    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}

func (nr *Newrelic) makeParamsRequest(url string, vals url.Values) ([]byte, error) {
    // TODO: Is there a better way to generate the request?
    url = fmt.Sprintf("%s/%s.%s?%s", nr.BaseUrl, url, nr.Format, vals.Encode())

    client := http.Client{}
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("X-Api-Key", nr.Key)
    resp, _ := client.Do(req)

    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}


func (nr *Newrelic) GetApplications() NewrelicApplications {
    resp, _ := nr.makeRequest("applications")

    var data NewrelicApplications

    json.Unmarshal(resp, &data)

    return data
}

func (nr *Newrelic) getBaseMetricData(invoke_url string, vals url.Values) NewrelicMetricData {
    resp, _ := nr.makeParamsRequest(invoke_url, vals)

    var data NewrelicMetricData

    json.Unmarshal(resp, &data)

    return data
}

func (nr *Newrelic) GetDefaultMetricData(app_id int, names []string) NewrelicMetricData {
    invoke_url := fmt.Sprintf("applications/%d/metrics/data", app_id)

    vals := url.Values{}
    for _, value := range names {
        vals.Add("names[]", value)
    }

    return nr.getBaseMetricData(invoke_url, vals)
}

func (nr *Newrelic) GetMetricData(app_id int, vals url.Values) NewrelicMetricData {
    invoke_url := fmt.Sprintf("applications/%d/metrics/data", app_id)

    return nr.getBaseMetricData(invoke_url, vals)
}
