package application

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Application struct {
	log     Logger
	grabber Grabber
	resizer Resizer
	cache   Cache
}

type Logger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

type Grabber interface {
	Grab(string, http.Header) ([]byte, error)
}

type Resizer interface {
	Resize([]byte, int, int) ([]byte, error)
}

type Cache interface {
	Get(string) ([]byte, bool, error)
	Set(string, *[]byte) error
	Clear() error
}

func (a *Application) GetServiceHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//Parse input path + checks
		imgWidth, imgHeight, imgUrl, err := parseInput(r.URL.Path)
		if err != nil {
			a.log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//Cache, grab and resize
		resImg, err := a.processImage(imgUrl, imgWidth, imgHeight, r)
		if err != nil {
			a.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Result
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/jpeg")
		w.Header().Add("Content-Length", strconv.Itoa(len(resImg)))
		_, err = w.Write(resImg)
		if err != nil {
			a.log.Error(err)
		}
	})

	return mux
}

func (a *Application) processImage(imgUrl string, imgWidth, imgHeight int, r *http.Request) ([]byte, error) {
	//Grab image form remote source
	srcImg, err := a.grabber.Grab(imgUrl, r.Header)
	if err != nil {
		return nil, err
	}

	//Cache check + update if needed
	if resImg, ok, err := a.cache.Get(r.URL.Path); ok {
		return resImg, err
	}

	//Process image
	resImg, err := a.resizer.Resize(srcImg, imgWidth, imgHeight)
	if err != nil {
		return nil, err
	}

	//Cache update
	err = a.cache.Set(r.URL.Path, &resImg)
	if err != nil {
		return nil, fmt.Errorf("cache set error %w", err)
	}

	return resImg, nil
}

func parseInput(input string) (int, int, string, error) {
	inputSet := strings.SplitN(input, "/", 4)
	if len(inputSet) != 4 {
		return 0, 0, "", fmt.Errorf("not enough parameters, got only %v", len(inputSet))
	}

	width, err := strconv.Atoi(inputSet[1])
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed in parse width: %w", err)
	}

	height, err := strconv.Atoi(inputSet[2])
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed in parse height: %w", err)
	}

	u, err := url.Parse(inputSet[3])
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed in parse URL: %w", err)
	}

	//http-only logic check
	switch u.Scheme {
	case "":
		u.Scheme = "http"
	case "http":
		break
	default:
		return 0, 0, "", fmt.Errorf("only http scheme allowed for this service, got %v", u.Scheme)
	}

	return width, height, u.String(), nil
}

func NewApp(logger Logger, grabber Grabber, resizer Resizer, cache Cache) *Application {
	return &Application{
		log:     logger,
		grabber: grabber,
		resizer: resizer,
		cache:   cache,
	}
}
