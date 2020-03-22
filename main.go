package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

var storedResponse *responseBody

const endpoint = "https://covid19.mathdro.id/api/countries/ID"

type responseBody struct {
	Confirmed struct {
		Value  int    `json:"value"`
		Detail string `json:"detail"`
	} `json:"confirmed"`
	Recovered struct {
		Value  int    `json:"value"`
		Detail string `json:"detail"`
	} `json:"recovered"`
	Deaths struct {
		Value  int    `json:"value"`
		Detail string `json:"detail"`
	} `json:"deaths"`
	LastUpdate time.Time `json:"lastUpdate"`
}

type persentasi struct {
	recovered string
	deaths    string
	active    int
}

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("could not read config: %s", err))
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigs)
	defer close(done)

	go func() {
		<-sigs
		done <- true
	}()

	ticker := time.NewTicker(viper.GetDuration("check_every"))
	go func() {
		for {
			select {
			case <-ticker.C:

				doJob()
			}
		}
	}()

	<-done
	ticker.Stop()

	os.Exit(0)
}

func doJob() {

	ctx := context.Background()
	var client = &http.Client{}
	var data responseBody

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatalf("Could not make new request: %s", err.Error())
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Could not Do Request: %s", err.Error())
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatalf("Could not decode response body: %s", err.Error())
	}

	if storedResponse == nil {
		storedResponse = &data
	} else {
		if data.LastUpdate.Unix() == storedResponse.LastUpdate.Unix() {
			return
		}
	}

	p := hitungPersentasi(&data)

	var param = url.Values{}
	content := fmt.Sprintf(
		"-----------\nTerkonfirmasi: %d \nTelah Sembuh: %d / %s%% \nKematian: %d / %s%%\nDalam Pengawasan: %d\nUpdate Terakhir: %s",
		data.Confirmed.Value,
		data.Recovered.Value,
		p.recovered,
		data.Deaths.Value,
		p.deaths,
		p.active,
		data.LastUpdate.Format("2 Jan 2006 15:04:05"),
	)
	param.Set("content", content)
	var payload = bytes.NewBufferString(param.Encode())

	request, err = http.NewRequestWithContext(ctx, http.MethodPost, viper.GetString("discord_webhook"), payload)
	if err != nil {
		log.Fatalf("Could not make new request: %s", err.Error())
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err = client.Do(request)
	if err != nil {
		log.Fatalf("Could not Do Request: %s", err.Error())
	}

	defer response.Body.Close()

}

func hitungPersentasi(r *responseBody) persentasi {
	var p persentasi

	pr := float32(r.Recovered.Value) / float32(r.Confirmed.Value) * 100
	pd := float32(r.Deaths.Value) / float32(r.Confirmed.Value) * 100

	prs := func() string {
		s := fmt.Sprintf("%f", pr)
		if strings.Contains(s, ".") {
			ss := strings.Split(s, ".")
			ss[1] = ss[1][0:1]
			return strings.Join(ss, ".")
		}
		return s
	}()
	pds := func() string {
		s := fmt.Sprintf("%f", pd)
		if strings.Contains(s, ".") {
			ss := strings.Split(s, ".")
			ss[1] = ss[1][0:1]
			return strings.Join(ss, ".")
		}
		return s
	}()

	p.recovered = prs
	p.deaths = pds
	p.active = r.Confirmed.Value - r.Recovered.Value - r.Deaths.Value

	return p
}
