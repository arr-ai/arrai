package main

import (
	"net"
	"net/http"

	"github.com/arr-ai/arrai/engine"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var serveCommand = cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Usage:   "start arrai as a gRPC server",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: ":42241",
			Usage: "address to listen on",
		},
		cli.StringFlag{
			Name:  "ws",
			Value: ":42242",
			Usage: "address to listen on",
		},
		cli.StringFlag{
			Name:  "cert",
			Usage: "TLS certificate file",
		},
		cli.StringFlag{
			Name:  "key",
			Usage: "TLS private key file",
		},
	},
	Action: serve,
}

func serve(c *cli.Context) error {
	listen := arraiAddress(c.String("listen"))
	wsListen := arraiAddressWithPort(c.String("ws"), 42242)
	cert := c.String("cert")
	key := c.String("key")

	eng := engine.Start()

	errors := make(chan error)

	go func() {
		lis, err := net.Listen("tcp", wsListen)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		wsFrontend := newWebsocketFrontend(eng)
		srv := &http.Server{
			Addr:    listen,
			Handler: http.HandlerFunc(wsFrontend.ServeHTTP),
		}

		log.Printf("Websocket server listening on " + wsListen)
		errors <- srv.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", listen)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := newGrpcServer(cert, key, eng)

		log.Printf("gRPC server listening on " + listen)
		errors <- grpcServer.Serve(lis)
	}()

	if err := <-errors; err != nil {
		return err
	}
	return <-errors
}
