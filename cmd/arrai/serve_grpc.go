package main

import (
	"io"

	"github.com/arr-ai/arrai/engine"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	pb "github.com/arr-ai/proto"
	"github.com/arr-ai/wbnf/parser"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func newGrpcServer(cert, key string, e *engine.Engine) *grpc.Server {
	var opts []grpc.ServerOption
	if cert != "" || key != "" {
		if !(cert != "" && key != "") {
			logrus.Fatal("TLS cert and key must be supplied together")
		}
		creds, err := credentials.NewServerTLSFromFile(cert, key)
		if err != nil {
			logrus.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	grpcServer := grpc.NewServer(opts...)
	server := arraiServer{e}
	pb.RegisterArraiServer(grpcServer, &server)
	return grpcServer
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
			logrus.Errorf("Error in arraiServer.Update: %v", err)
			return err
		}
		logrus.Infof("req.Expr: %s", req.Expr)
		pc := syntax.ParseContext{SourceDir: "-"}
		ast, err := pc.Parse(parser.NewScanner(req.Expr))
		if err != nil {
			logrus.Errorf("Error in arraiServer.Update: %v", err)
			return err
		}
		expr := pc.CompileExpr(ast)
		logrus.Info("Parsed successfully")
		err = s.engine.Update(expr)
		if err != nil {
			logrus.Errorf("Error in arraiServer.Update: %v", err)
			return err
		}
		if err = stream.Send(&ack); err != nil {
			logrus.Errorf("Error in arraiServer.Update: %v", err)
			return err
		}
	}
}

func (s *arraiServer) Observe(
	req *pb.ObserveReq, stream pb.Arrai_ObserveServer,
) error {
	var pc syntax.ParseContext
	ast, err := pc.Parse(parser.NewScanner(req.Expr))
	if err != nil {
		return err
	}
	expr := pc.CompileExpr(ast)
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
