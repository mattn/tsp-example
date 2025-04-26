package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/mattn/tsp-example/api"
	"github.com/urfave/cli/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client, err := api.NewClient("http://localhost:8888")
	if err != nil {
		log.Fatal(err)
	}

	var app = &cli.Command{
		Name: "client",
		Commands: []*cli.Command{
			&cli.Command{
				Name: "new",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Present() {
						return errors.New("expected exactly an argument")
					}
					_, err := client.TodosCreate(ctx, &api.Todo{
						Content: c.Args().Get(0),
					})
					return err
				},
			},
			&cli.Command{
				Name: "delete",
				Action: func(ctx context.Context, c *cli.Command) error {
					return client.TodosDelete(ctx, api.TodosDeleteParams{
						ID: c.Args().Get(0),
					})
				},
			},
			&cli.Command{
				Name: "update",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Len() != 2 {
						return errors.New("expected exactly two arguments")
					}
					_, err = client.TodosUpdate(ctx, &api.TodoUpdate{
						Content: api.NewOptString(c.Args().Get(1)),
					}, api.TodosUpdateParams{
						ID: c.Args().Get(0),
					})
					return err
				},
			},
			&cli.Command{
				Name: "done",
				Action: func(ctx context.Context, c *cli.Command) error {
					if !c.Args().Present() {
						return errors.New("expected exactly an argument")
					}
					_, err = client.TodosUpdate(ctx, &api.TodoUpdate{
						Done: api.NewOptBool(true),
					}, api.TodosUpdateParams{
						ID: c.Args().Get(0),
					})
					return err
				},
			},
			&cli.Command{
				Name: "list",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "json",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					items, err := client.TodosList(ctx)
					if err != nil {
						return err
					}
					if c.Bool("json") {
						return json.NewEncoder(os.Stdout).Encode(items.Items)
					}
					for _, item := range items.Items {
						if item.GetDone() {
							fmt.Fprintln(color.Output, item.GetID(), color.BlueString(item.GetContent()))
						} else {
							fmt.Fprintln(color.Output, item.GetID(), color.WhiteString(item.GetContent()))
						}
					}
					return err
				},
			},
		},
	}

	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
