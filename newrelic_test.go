package newrelic_api

import "testing"

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
