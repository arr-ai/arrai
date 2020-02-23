package main

import (
	"net/http"

	"github.com/arr-ai/arrai/engine"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type websocketFrontend struct {
	engine   *engine.Engine
	upgrader websocket.Upgrader
}

func newWebsocketFrontend(eng *engine.Engine) *websocketFrontend {
	return &websocketFrontend{
		eng,
		websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
	}
}

func (wsfe *websocketFrontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsfe.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading websocket: %s", err)
		return
	}
	defer conn.Close()

	log.Info("Websocket connected")

	cancelObservation := func() {}
	alive := true
	for alive {
		msgtype, p, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Error reading websocket: %s", err)
			return
		}
		if msgtype == websocket.TextMessage {
			s := string(p)
			var pc syntax.ParseContext
			ast, err := pc.ParseString(s)
			if err != nil {
				log.Errorf("Error parsing request %#v: %s", s, err)
				if err := conn.WriteJSON(map[string]interface{}{"error": err.Error()}); err != nil {
					panic(err)
				}
				continue
			}
			expr := pc.CompileExpr(ast)

			cancelObservation()
			cancelObservation = wsfe.engine.Observe(
				expr,
				func(value rel.Value) error {
					return conn.WriteMessage(
						websocket.TextMessage,
						rel.MarshalToJSON(value),
					)
				},
				func(err error) {
					alive = false
				},
			)
		}
	}
}
