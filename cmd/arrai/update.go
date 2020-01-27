package main

import (
	"context"

	pb "github.com/arr-ai/proto"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var updateCommand = &cli.Command{
	Name:    "update",
	Aliases: []string{"u"},
	Usage:   "update a server with an expression",
	Action:  update,
}

func update(c *cli.Context) error {
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

	update, err := client.Update(context.Background())
	if err != nil {
		return err
	}

	if err := update.Send(&pb.UpdateReq{Expr: source}); err != nil {
		return err
	}

	_, err = update.Recv()
	if err != nil {
		return err
	}

	return nil
}
