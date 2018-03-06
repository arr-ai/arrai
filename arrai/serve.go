package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/arr-ai/arrai/engine"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	pb "github.com/arr-ai/proto"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
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

type arraiServer struct {
	engine *engine.Engine
}

func (s *arraiServer) Update(stream pb.Arrai_UpdateServer) error {
	ack := pb.UpdateAck{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		expr, err := syntax.Parse(bytes.NewBufferString(req.Expr))
		if err != nil {
			return err
		}
		err = s.engine.Update(expr)
		if err != nil {
			return err
		}
		if err = stream.Send(&ack); err != nil {
			return err
		}
	}
}

func (s *arraiServer) Observe(
	req *pb.ObserveReq, stream pb.Arrai_ObserveServer,
) error {
	expr, err := syntax.Parse(bytes.NewBufferString(req.Expr))
	if err != nil {
		return err
	}
	retch := make(chan error)

	send := func(resp *pb.ObserveResp) error {
		if err = stream.Send(resp); err != nil {
			retch <- err
			return err
		}
		return nil
	}

	onupdate := func(value rel.Value) error {
		return send(
			&pb.ObserveResp{
				Value: &pb.Value{
					Choice: &pb.Value_Json{
						Json: string(rel.MarshalToJSON(value)),
					},
				},
			},
		)
	}

	onclose := func(err error) {
		retch <- err
	}

	s.engine.Observe(expr, onupdate, onclose)

	return <-retch
}

func serve(c *cli.Context) error {
	listen := arraiAddress(c.String("listen"))
	cert := c.String("cert")
	key := c.String("key")

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if cert != "" || key != "" {
		if !(cert != "" && key != "") {
			grpclog.Fatal("TLS cert and key must be supplied together")
		}
		creds, err := credentials.NewServerTLSFromFile(cert, key)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	grpcServer := grpc.NewServer(opts...)
	server := arraiServer{engine.Start()}
	pb.RegisterArraiServer(grpcServer, &server)

	hup := make(chan os.Signal)
	go func() {
		for {
			<-hup
			log.Printf("Received SIGHUP. Hanging up all connections...")
			server.engine.Hangup()
		}
	}()
	signal.Notify(hup, syscall.SIGHUP)

	log.Printf("Listening on " + listen)
	return grpcServer.Serve(lis)
}
