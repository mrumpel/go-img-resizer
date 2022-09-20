package cache

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	c, err := NewCache(3, dir)
	require.NoError(t, err)

	b := []byte("some bytes for test")

	t.Run("base test", func(t *testing.T) {
		err = c.Set("/100/100/test", &b)
		require.NoError(t, err)
		require.Equal(t, c.queue.Len(), 1)
		require.Len(t, c.items, 1)

		d, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, d, 1)
	})

	t.Run("max size test", func(t *testing.T) {
		err = c.Set("/100/200/test", &b)
		require.NoError(t, err)
		err = c.Set("/100/300/test", &b)
		require.NoError(t, err)
		err = c.Set("/100/400/test", &b)
		require.NoError(t, err)
		err = c.Set("/100/500/test", &b)
		require.NoError(t, err)

		require.Equal(t, c.queue.Len(), 3)
		require.Len(t, c.items, 3)

		_, ok, _ := c.Get("/100/200/test")
		require.False(t, ok)

		d, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, d, 3)
	})
}

func TestGet(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	c, err := NewCache(1, dir)
	require.NoError(t, err)

	b := []byte("some bytes for test")

	err = c.Set("/100/100/test", &b)
	require.NoError(t, err)

	res, ok, err := c.Get("/100/100/test")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, res, b)
}

func TestNewCache(t *testing.T) {
	t.Run("dir", func(t *testing.T) {
		dir, err := os.MkdirTemp(os.TempDir(), "")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		c, err := NewCache(1, dir)
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("no dir", func(t *testing.T) {
		c, err := NewCache(1, "")
		require.NoError(t, err)
		require.NotEmpty(t, c.dir)

		_, err = os.ReadDir(c.dir)
		require.NoError(t, err)
	})

	t.Run("dir not exists", func(t *testing.T) {
		dir, err := os.MkdirTemp(os.TempDir(), "")
		require.NoError(t, err)
		err = os.RemoveAll(dir)
		require.NoError(t, err)

		c, err := NewCache(1, dir)
		require.Nil(t, c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cache dir")
	})

	t.Run("wrong size", func(t *testing.T) {
		c, err := NewCache(0, "")
		require.Nil(t, c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cache size")
	})
}

func TestClear(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "")
	require.NoError(t, err)

	c, err := NewCache(1, dir)
	require.NoError(t, err)
	b := []byte("some bytes for test")

	err = c.Set("key", &b)
	require.NoError(t, err)
	d, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, d, 1)

	err = c.Clear()
	require.NoError(t, err)
	_, err = os.ReadDir(dir)
	require.Error(t, err)
}
