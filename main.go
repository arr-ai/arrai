package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-errors/errors"
	"github.com/gorilla/websocket"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"

	"github.com/arr-ai/arrai/engine"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

var (
	listen = flag.String("listen", ":42241", "Address to listen on!")

	root = "$"
)

func main() {
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	type watcher struct {
		update chan<- *rel.Scope
		id     uint64
	}

	type updateRequest struct {
		expr   rel.Expr
		errors chan<- error
	}

	engine := engine.Start()

	r.Get("/.__ws__", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("New connection")

		updateCh := make(chan rel.Value)
		errorsCh := make(chan error)

		// Deliver updates and errors to client.
		go func() {
			defer func() {
				log.Printf("Connection closing")
				conn.Close()
			}()

			reportError := func(err error) {
				log.Printf("Sending error to client: %v", err)
				j, err := json.Marshal(map[string]string{"error": err.Error()})
				if err != nil {
					panic(err)
				}
				conn.WriteMessage(websocket.TextMessage, j)
			}

			for {
				select {
				case value, ok := <-updateCh:
					if !ok {
						return
					}
					log.Printf("<- %s", value)
					func() {
						defer func() {
							if r := recover(); r != nil {
								reportError(errors.Errorf("panic(%s)", r))
							}
						}()
						err = conn.WriteMessage(
							websocket.TextMessage, rel.MarshalToJSON(value))
					}()
				case err, ok := <-errorsCh:
					if !ok {
						return
					}
					reportError(err)
				}
			}
		}()

		cancelObs := func() {}
		defer func() { cancelObs() }() // Inside func to lazy-eval

		// Process messages from client.
		for {
			_, code, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Received code: %s", code)
			if code[0] == '=' {
				expr, err := syntax.Parse(code[1:])
				if err != nil {
					errorsCh <- err
				} else {
					engine.Update(expr, errorsCh)
				}
			} else {
				expr, err := syntax.Parse(code)
				if err != nil {
					errorsCh <- err
				} else {
					cancelObs()
					cancelObs = engine.Observe(expr, updateCh, errorsCh)
				}
			}
		}
	})

	hup := make(chan os.Signal)
	go func() {
		for {
			<-hup
			log.Printf("Received SIGHUP. Hanging up all connections...")
			engine.Hangup()
		}
	}()
	signal.Notify(hup, syscall.SIGHUP)

	log.Printf("Listening on " + *listen)
	http.ListenAndServe(*listen, r)
}
