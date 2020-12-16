package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return Echo(ctx)
	})
	eg.Go(func() error {
		return ReceiveSignal(ctx)
	})
	eg.Go(func() error {
		return Cmd(ctx)
	})
	fmt.Println(eg.Wait())
}

func ReceiveSignal(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-c:
		return errors.New(sig.String() + "\n")
	case <-ctx.Done():
		return nil
	}
}

func Echo(ctx context.Context) error {

	c := make(chan error)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(rw, "%v+", err)
			return
		}
		if reflect.DeepEqual(data, []byte("bye")) {
			c <- errors.New("bye")
			rw.Write([]byte("bye bye\n"))
			return
		}
		rw.Write(data)
		rw.Write([]byte("\n"))
	})

	go func() {
		if err := http.ListenAndServe(":8089", http.DefaultServeMux); err != nil {
			c <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-c:
		return err
	}
}

func Cmd(ctx context.Context) error {
	cmd := exec.Command("sh", "script.sh")
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	if err := cmd.Start(); err != nil {
		return err
	}
	c := make(chan error)
	go func() {
		// TODO Pipe
		// cmd.StdoutPipe()
		c <- cmd.Wait()
	}()
	for {
		select {
		case err := <-c:
			return err
		default:
			if data, err := buf.ReadString('\n'); err != nil && !errors.Is(err, io.EOF) {
				return err
			} else if data != "" {
				fmt.Println(data)
			}
		}
	}
}
