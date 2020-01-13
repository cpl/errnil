package errnil

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

const (
	errNilString = "err != nil"
	stopTimeout  = time.Second * 10
)

type Counter struct {
	workers int
	count   uint32

	input chan string
	close chan interface{}
	errch chan error
}

func NewCounter(workers int) *Counter {
	return &Counter{
		workers: workers,
		close:   make(chan interface{}),
		input:   make(chan string, workers*100),
		errch:   make(chan error, workers*10),
	}
}

func (c *Counter) Count(dir string) (uint32, error) {
	for iter := 0; iter < c.workers; iter++ {
		go c.work()
	}

	if err := filepath.Walk(dir, c.visit); err != nil {
		c.stop()
		return 0, fmt.Errorf("failed walking path, %w", err)
	}

	if err := c.stop(); err != nil {
		return 0, fmt.Errorf("failed closing workers, %w", err)
	}

	return atomic.LoadUint32(&c.count), nil
}

func (c *Counter) visit(p string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if f.IsDir() || path.Ext(p) != ".go" {
		return nil
	}

	select {
	case err := <-c.errch:
		return err
	default:
		c.input <- p
		return nil
	}
}

func (c *Counter) stop() error {
	ctx, _ := context.WithTimeout(context.Background(), stopTimeout)
	for iter := 0; iter < c.workers; iter++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout trying to stop counter")
		default:
			c.close <- nil
		}
	}
	return nil
}

func (c *Counter) work() {
	for {
		select {
		case <-c.close:
			return
		case file := <-c.input:
			source, err := ioutil.ReadFile(file)
			if err != nil {
				c.errch <- err
				continue
			}

			source, err = format.Source(source)
			if err != nil {
				c.errch <- err
				continue
			}

			count, err := countErrNil(source)
			if err != nil {
				c.errch <- err
				continue
			}

			atomic.AddUint32(&c.count, count)
		}
	}
}

// TODO: add checks against comments and fake Go files
func countErrNil(source []byte) (uint32, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(source))
	scanner.Split(bufio.ScanLines)

	var count uint32
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, errNilString) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed scanning source, %w", err)
	}

	return count, nil
}
