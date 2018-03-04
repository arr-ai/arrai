package main

import (
	"context"
	"os"

	pb "github.com/arr-ai/proto"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var observeCommand = cli.Command{
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
	observe, err := client.Observe(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		resp, err := observe.Recv()
		if err != nil {
			return err
		}
		s := resp.Value.String()
		os.Stdout.WriteString(s)
		if s[len(s)-1] != '\n' {
			os.Stdout.Write([]byte{'\n'})
		}
	}
}
