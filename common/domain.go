package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// refer: scheduler/stark/common/domain.go
func GetHostsBydomain(domain string) ([]string, error) {
	uri := domain
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type getdomainByServiceResponse struct {
		Code int `json:"code"`
		Data []struct {
			HostName     string `json:"host_name"`
			ServiceNames string `json:"service_names"`
			HostIP       int64  `json:"host_ip"`
			Instance     struct {
				Port   int `json:"port"`
				Status int `json:"status"`
			} `json:"instance_status"`
		} `json:"data"`
	}
	var res getdomainByServiceResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if res.Code != 0 {
		return nil, errors.New(
			fmt.Sprintf("resp of getdomainByServiceResponse is error. code = %d", res.Code))
	}

	ips := make([]string, 0)
	for _, host := range res.Data {
		if host.Instance.Status == 0 {
			ips = append(ips, inettoString(host.HostIP)+":"+strconv.Itoa(host.Instance.Port))
		}
	}
	return ips, nil
}

func inettoString(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
