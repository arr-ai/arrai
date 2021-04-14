// +build !darwin

package main

// TODO: Allow watching a single file to represent a non-tuple database.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	pb "github.com/arr-ai/proto"
	"github.com/go-errors/errors"
	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func sync(c *cli.Context) error {
	server := arraiAddress(c.Args().Get(0))
	dir := c.Args().Get(1)

	template := c.String("template")

	if dir == "" {
		dir = "."
	}

	log.Infof("xServer: %s", server)
	log.Infof("Directory: %s", dir)

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewArraiClient(conn)

	eich := make(chan notify.EventInfo, 1)
	eich <- nil

	watch := path.Join(dir, "...")
	if err := notify.Watch(watch, eich, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(eich)

	update, err := client.Update(arraictx.InitRunCtx(context.Background()))
	if err != nil {
		return err
	}

	log.Printf("Watching %q", watch)
	ctx := arraictx.InitRunCtx(context.TODO())
	for {
		ei := <-eich

		// Give Excel et al a chance to finish their work.
		time.Sleep(300 * time.Millisecond)

		log.Printf("EI: %v", ei)

		tree, err := buildTree(ctx, path.Clean(dir))
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				log.Printf("TREE BUILDING ERROR: %s", e.ErrorStack())
			} else {
				log.Printf("TREE BUILDING ERROR: %s", err)
			}
		}
		var buf bytes.Buffer
		writeTreeToBuffer(tree, &buf)
		log.Printf("TREE: %v", buf.String())
		if err := update.Send(&pb.UpdateReq{Expr: fmt.Sprintf(template, buf.String())}); err != nil {
			return err
		}
		if _, err := update.Recv(); err != nil {
			return err
		}
	}
}

var handlers = map[string]func([]byte) ([]byte, error){
	// TODO: Add XLSX support with ".xlsx": handleXlsx, (see commit faaeff7).
}

func writeTreeToBuffer(tree map[string]interface{}, buf *bytes.Buffer) {
	buf.WriteRune('{')
	for name, value := range tree {
		jname, err := json.Marshal(name)
		if err != nil {
			panic(err)
		}
		buf.Write(jname)
		buf.WriteRune(':')
		switch x := value.(type) {
		case map[string]interface{}:
			writeTreeToBuffer(x, buf)
		case []byte:
			buf.Write(x)
		default:
			panic("Bad value writing tree to buffer")
		}
		buf.WriteRune(',')
	}
	buf.WriteRune('}')
}

func buildTree(ctx context.Context, root string) (map[string]interface{}, error) {
	tree := map[string]interface{}{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("ERROR WALKING TO %s: %s", path, err)
			} else if !info.IsDir() {
				data, err := afero.ReadFile(ctxfs.SourceFsFrom(ctx), path)
				if err != nil {
					return errors.WrapPrefix(err, "reading "+path, 0)
				}

				parts := strings.Split(filepath.ToSlash(path), "/")
				log.Infof("path: %#v %#v", path, parts)
				// Splice off the last part.
				parts, attr := parts[:len(parts)-1], parts[len(parts)-1]
				if attr == ".DS_Store" || strings.HasPrefix(attr, "~$") {
					return nil
				}
				log.Println(path)
				if strings.HasSuffix(attr, ".arrai") {
					attr = attr[:len(attr)-6]
				}
				for ext, handler := range handlers {
					if strings.HasSuffix(attr, ext) {
						data, err = handler(data)
						if err != nil {
							return err
						}
					}
				}
				node := tree
				for _, part := range parts {
					if value, ok := node[part]; ok {
						if subnode, ok := value.(map[string]interface{}); ok {
							node = subnode
						} else {
							return errors.Errorf("Walking to non-node")
						}
					} else {
						subnode := map[string]interface{}{}
						node[part] = subnode
						node = subnode
					}
				}
				node[attr] = data
			}
			return nil
		})
	return tree, err
}
