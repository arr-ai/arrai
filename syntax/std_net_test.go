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
	srv := startHTTPServer(t, addr)
	defer func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	wait(context.Background(), t, url)

	AssertCodesEvalToSameValue(t, `(
		status_code: 200,
		status: "200 OK",
		body: <<"hello world">>,
		header: (
			"Content-Length": ["11"],
			"Content-Type": ["text/plain; charset=utf-8"],
		),
	)`, fmt.Sprintf(`//net.http.get((), "%s")`, url))
}

func TestNetGet_NoUrl(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.get((), "")`)
	AssertCodeErrors(t, "", `//net.http.get((), ["localhost"])`)
}

func TestNetGet_BadHeader(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.get((), "")`)
	AssertCodeErrors(t, "", `//net.http.get((), ["localhost"])`)
}

func TestNetPost(t *testing.T) {
	t.Parallel()

	addr := "localhost:57512"
	url := fmt.Sprintf("http://%s", addr)
	srv := startHTTPServer(t, addr)
	defer func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	wait(context.Background(), t, url)

	AssertCodesEvalToSameValue(t, `(
		status_code: 200,
		status: "200 OK",
		body: <<"foo">>,
		header: (
			"Content-Length": ["3"],
			"Content-Type": ["application/sysl"],
		),
	)`, fmt.Sprintf(`//net.http.post((header: {"Content-Type": "application/sysl"}), "%s", "foo")`, url))
}

func TestNetPost_BadConfig(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.post({}, "localhost", "")`)
	AssertCodeErrors(t, "", `//net.http.post((header: 'Content-Type: text/plain'), "localhost", "")`)
	AssertCodeErrors(t, "", `//net.http.post((header: ('Content-Type': 'text/plain')), "localhost", "")`)
	AssertCodeErrors(t, "", `//net.http.post((header: {123: 'text/plain'}), "localhost", "")`)
}

func TestNetPost_NoUrl(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "", `//net.http.post((), "", "")`)
	AssertCodeErrors(t, "", `//net.http.post((), ["localhost"], "")`)
}

// testHTTPHandler is a dummy HTTP handler that serves "hello world".
type testHTTPHandler struct {
	t *testing.T
}

// ServeHTTP writes responses based on the content of the request.
func (h testHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()["Date"] = nil // Remove non-deterministic Date header.

	bs, err := ioutil.ReadAll(r.Body)
	require.NoError(h.t, err)

	if ct, ok := r.Header["Content-Type"]; ok {
		w.Header()["Content-Type"] = ct
	}

	if len(bs) > 0 {
		_, err = w.Write(bs)
		require.NoError(h.t, err)
	} else {
		_, err = w.Write([]byte("hello world"))
		require.NoError(h.t, err)
	}
}

// startHTTPServer creates, starts and returns a test server at addr.
func startHTTPServer(t *testing.T, addr string) *http.Server {
	srv := &http.Server{Addr: addr, Handler: testHTTPHandler{t}}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Fail()
		}
	}()
	return srv
}

// wait pings the server at url until it responds successfully or times out.
func wait(ctx context.Context, t *testing.T, url string) {
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
