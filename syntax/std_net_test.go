package syntax

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/require"
)

func TestNetGet(t *testing.T) {
	t.Parallel()

	addr := "localhost:57511"
	url := fmt.Sprintf("http://%s", addr)
	srv := startHttpServer(t, addr)
	defer func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	wait(t, context.Background(), url)

	AssertCodesEvalToSameValue(t, `(
		status_code: 200,
		status: "200 OK",
		body: <<"hello world">>,
		header: (
			"Content-Length": ["11"],
			"Content-Type": ["text/plain; charset=utf-8"],
		),
	)`, fmt.Sprintf(`//net.http.get("%s")`, url))
}

func TestNetGet_NoUrl(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.get("")`)
	AssertCodeErrors(t, "", `//net.http.get(["localhost"])`)
}

func TestNetPost(t *testing.T) {
	t.Parallel()

	addr := "localhost:57512"
	url := fmt.Sprintf("http://%s", addr)
	srv := startHttpServer(t, addr)
	defer func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	wait(t, context.Background(), url)

	AssertCodesEvalToSameValue(t, `(
		status_code: 200,
		status: "200 OK",
		body: <<"foo">>,
		header: (
			"Content-Length": ["3"],
			"Content-Type": ["text/plain; charset=utf-8"],
		),
	)`, fmt.Sprintf(`//net.http.post((body: "foo"), "%s")`, url))
}

func TestNetPost_NoConfig(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.post({}, "localhost")`)
}

func TestNetPost_NoUrl(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.post((), "")`)
	AssertCodeErrors(t, "", `//net.http.post((), ["localhost"])`)
}

// testHttpHandler is a dummy HTTP handler that serves "hello world".
type testHttpHandler struct{}

// ServeHTTP writes responses based on the content of the request.
func (h testHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()["Date"] = nil // Remove non-deterministic Date header.

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if ct, ok := r.Header["Content-Type"]; ok {
		w.Header()["Content-Type"] = ct
	}

	if bs != nil {
		w.Write(bs)
	} else {
		w.Write([]byte("hello world"))
	}
}

// startHttpServer creates, starts and returns a test server at addr.
func startHttpServer(t *testing.T, addr string) *http.Server {
	srv := &http.Server{Addr: addr, Handler: testHttpHandler{}}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Fail()
		}
	}()
	return srv
}

// wait pings the server at url until it responds successfully or times out.
func wait(t *testing.T, ctx context.Context, url string) {
	backoff, err := retry.NewFibonacci(10 * time.Millisecond)
	require.Nil(t, err)
	backoff = retry.WithMaxDuration(3*time.Second, backoff)
	err = retry.Do(ctx, backoff, func(ctx context.Context) error {
		if err := ping(ctx, url); err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	require.NoError(t, err)
}

// ping checks if a server is available at url, and returns an error if not.
func ping(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("not up")
	}
	return nil
}
