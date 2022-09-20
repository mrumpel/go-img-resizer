package cache

import (
	"fmt"
	"os"
	"sync"
)

type ImgCache struct {
	mu      sync.Mutex
	queue   listInterface
	items   map[string]*listItem
	dir     string
	maxSize int
}

func (c *ImgCache) Get(key string) ([]byte, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false, nil
	}

	c.queue.MoveToFront(item)

	img, err := os.ReadFile(item.Value)
	if err != nil {
		return nil, false, fmt.Errorf("cache error %w", err)
	}

	return img, true, nil
}

func (c *ImgCache) Set(key string, img *[]byte) error {
	c.mu.Lock()
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	file, err := os.CreateTemp(c.dir, "")
	if err != nil {
		return fmt.Errorf("creating file error %w", err)
	}
	_, err = file.Write(*img)
	if err != nil {
		return fmt.Errorf("writing file error %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing file error %w", err)
	}

	c.mu.Lock()
	if c.queue.Len() == c.maxSize {
		err := os.Remove(c.queue.Back().Value)
		if err != nil {
			return fmt.Errorf("deleting file error %w", err)
		}
		delete(c.items, c.queue.Back().Key)
		c.queue.Remove(c.queue.Back())
	}
	c.queue.PushFront(key, file.Name())
	c.items[key] = c.queue.Front()
	c.mu.Unlock()

	return nil
}

func NewCache(maxSize int, dir string) (*ImgCache, error) {
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp(os.TempDir(), "")
		if err != nil {
			return nil, fmt.Errorf("error creating cache dir %w", err)
		}
	}
	_, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error using cache dir %w", err)
	}

	if maxSize <= 0 {
		return nil, fmt.Errorf("wrong cache size limit: %v", maxSize)
	}

	return &ImgCache{
		mu:      sync.Mutex{},
		queue:   newList(),
		items:   make(map[string]*listItem),
		dir:     dir,
		maxSize: maxSize,
	}, nil
}

func (c *ImgCache) Clear() error {
	return os.RemoveAll(c.dir)
}
