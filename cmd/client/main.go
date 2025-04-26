package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

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
					_, err := client.TodosCreate(ctx, &api.Todo{
						Content: c.Args().First(),
					})
					return err
				},
			},
			&cli.Command{
				Name: "delete",
				Action: func(ctx context.Context, c *cli.Command) error {
					return client.TodosDelete(ctx, api.TodosDeleteParams{
						ID: c.Args().First(),
					})
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
						fmt.Println(item.GetID(), item.GetContent())
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
