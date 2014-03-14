package newrelic_api

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "net/url"
)

type Newrelic struct {
    Key string
    BaseUrl string
    Format string
}

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

type NewRelicMetricData struct {
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

func (nr *Newrelic) getBaseMetricData(invoke_url string, vals url.Values) NewRelicMetricData {
    resp, _ := nr.makeParamsRequest(invoke_url, vals)

    var data NewRelicMetricData

    json.Unmarshal(resp, &data)

    return data
}

/**
 * This will send a request applying the defaults for values, from, to and summarize
 **/
func (nr *Newrelic) GetDefaultMetricData(app_id int, names []string) NewRelicMetricData {
    invoke_url := fmt.Sprintf("applications/%d/metrics/data", app_id)

    vals := url.Values{}
    for _, value := range names {
        vals.Add("names[]", value)
    }

    return nr.getBaseMetricData(invoke_url, vals)
}

/**
 * vals should contain all the filtering mechanisms that you'd like to use
 **/
func (nr *Newrelic) GetMetricData(app_id int, vals url.Values) NewRelicMetricData {
    invoke_url := fmt.Sprintf("applications/%d/metrics/data", app_id)

    return nr.getBaseMetricData(invoke_url, vals)
}
