//go:build integration

package img_resizer_integration_test

import "github.com/stretchr/testify/suite"

import (
	"github.com/stretchr/testify/require"
	"image"
	_ "image/jpeg"
	"net/http"
	"os"
	"testing"
)

//TODO: use the suite
type imgTestSuite struct {
	suite.Suite
	host   string
	static string
}

func TestHealthz(t *testing.T) {
	z, err := http.Get("http://localhost:8080/healthz")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, z.StatusCode)
}

func TestResize(t *testing.T) {
	//TODO: funcs for integration tests
	z, err := http.Get("http://localhost:8080/500/250/web/_gopher_original_1024x504.jpg")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, z.StatusCode)
	defer z.Body.Close()

	img, _, err := image.DecodeConfig(z.Body)
	require.NoError(t, err)
	require.Equal(t, 500, img.Width)
	require.Equal(t, 250, img.Height)
}

func TestCache(t *testing.T) {
	b, err := os.ReadFile("static/_gopher_original_1024x504.jpg")
	require.NoError(t, err)
	require.NotEmpty(t, b)

	err = os.WriteFile("static/cache_gopher.jpg", b, 0644)
	require.NoError(t, err)

	z, err := http.Get("http://localhost:8080/500/250/web/cache_gopher.jpg")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, z.StatusCode)
	err = z.Body.Close()
	require.NoError(t, err)

	err = os.Remove("static/cache_gopher.jpg")
	require.NoError(t, err)

	z, err = http.Get("http://localhost:8080/500/250/web/cache_gopher.jpg")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, z.StatusCode)
	require.NotEmpty(t, z.Body)
	err = z.Body.Close()
	require.NoError(t, err)
}
