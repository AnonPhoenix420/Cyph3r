package netcheck


import (
"crypto/tls"
"net/http"
"time"
)


func HTTP(url string, insecure bool) (int, int64) {
tr := &http.Transport{}
if insecure {
tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}


client := http.Client{
Timeout: 5 * time.Second,
Transport: tr,
}


start := time.Now()
resp, err := client.Get(url)
if err != nil {
return 0, 0
}
defer resp.Body.Close()


return resp.StatusCode, time.Since(start).Milliseconds()
}
