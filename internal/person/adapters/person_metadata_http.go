package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type personMetaData struct {
	client         *http.Client
	agifyAPI       string
	genderizeAPI   string
	nationalizeAPI string
}

func PersonMetaData(
	agifyAPI string,
	genderizeAPI string,
	nationalizeAPI string,
) (*personMetaData, error) {
	if len(agifyAPI) == 0 {
		return nil, fmt.Errorf("empty agify service api")
	}
	if len(genderizeAPI) == 0 {
		return nil, fmt.Errorf("empty genderize service api")
	}
	if len(nationalizeAPI) == 0 {
		return nil, fmt.Errorf("empty nationalize service api")
	}
	return &personMetaData{
		agifyAPI:       agifyAPI,
		genderizeAPI:   genderizeAPI,
		nationalizeAPI: nationalizeAPI,
		client: &http.Client{
			Timeout: time.Second * 1,
			Transport: &http.Transport{
				MaxIdleConns: 15,
			},
		},
	}, nil
}

type ageByNameResponse struct {
	Age int
}

func (p *personMetaData) AgeByName(ctx context.Context, name string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, p.agifyAPI, nil)
	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var res ageByNameResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, err
	}
	if res.Age == 0 {
		return 0, fmt.Errorf("no age in response")
	}

	return res.Age, nil
}

type genderByNameResponse struct {
	Gender string
}

func (p *personMetaData) GenderByName(ctx context.Context, name string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, p.genderizeAPI, nil)
	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res genderByNameResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if len(res.Gender) == 0 {
		return "", fmt.Errorf("no gender in response")
	}

	return res.Gender, nil
}

type nationByNameResponse struct {
	Country []struct {
		Id          string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}
}

func (p *personMetaData) NationByName(ctx context.Context, name string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, p.nationalizeAPI, nil)
	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res nationByNameResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if len(res.Country) == 0 {
		return "", fmt.Errorf("no country in response")
	}

	return res.Country[0].Id, nil
}
