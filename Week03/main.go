package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	fmt.Println(eg.Wait())
}

func ReceiveSignal(ctx context.Context) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-c:
		return errors.New(sig.String())
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
		if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
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
