package grabber

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Grabber struct{}

func NewGrabber() *Grabber {
	return &Grabber{}
}

func (g *Grabber) Grab(rawURL string, header http.Header) ([]byte, error) {
	//1. Check URL and prepare request
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("error in the URL %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request %w", err)
	}
	req.URL = u
	req.Header = header

	//2. Do request
	//TODO: а нужен ли клиент внутри граббера? Свериться с требованиями
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in doing request %w", err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error in reading response body %w", err)
	}

	return res, nil
}
