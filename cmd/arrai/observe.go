package main

import (
	"context"
	"fmt"

	pb "github.com/arr-ai/proto"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
)

var observeCommand = &cli.Command{
	Name:    "observe",
	Aliases: []string{"o"},
	Usage:   "observe an expression on a server",
	Action:  observe,
}

func observe(c *cli.Context) error {
	server := arraiAddress(c.Args().Get(0))
	source := c.Args().Get(1)

	logrus.Infof("Server: %s", server)
	logrus.Infof("Source: %s", source)

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewArraiClient(conn)

	req := &pb.ObserveReq{Expr: source}
	observe, err := client.Observe(arraictx.InitRunCtx(context.Background()), req)
	if err != nil {
		return err
	}

	for {
		resp, err := observe.Recv()
		if err != nil {
			return err
		}
		json := resp.Value.Choice.(*pb.Value_Json).Json
		value, err := rel.UnmarshalFromJSON([]byte(json))
		if err != nil {
			return err
		}
		fmt.Println(value.String())
	}
}
