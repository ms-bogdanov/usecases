package service

import (
	"fmt"
	"io"
	"net/http"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) Search(query string) (string, error) {
	url := "https://nominatim.openstreetmap.org/search?q=" + query + "&format=json"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(tmp), nil
}

func (s *Service) Geocode(lat, lon float64) (string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(tmp), nil
}
