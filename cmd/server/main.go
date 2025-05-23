package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/mattn/tsp-example/api"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// TodoHandler is no-op Handler which returns http.ErrNotImplemented.
type TodoHandler struct {
	bundb *bun.DB
}

var _ api.Handler = &TodoHandler{}

// TodosCreate implements Todos_create operation.
//
// Create a widget.
//
// POST /widgets
func (h *TodoHandler) TodosCreate(ctx context.Context, req *api.Todo) (r *api.Todo, _ error) {
	_, err := h.bundb.NewInsert().Model(req).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// TodosDelete implements Todos_delete operation.
//
// Delete a widget.
//
// DELETE /widgets/{id}
func (h *TodoHandler) TodosDelete(ctx context.Context, params api.TodosDeleteParams) error {
	result, err := h.bundb.NewDelete().Model((*api.Todo)(nil)).Where(`"id" = ?`, params.ID).Exec(ctx)
	if err != nil {
		return h.NewError(ctx, err)
	}
	if num, err := result.RowsAffected(); err != nil || num == 0 {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response:   api.Error{Message: "No records deleted"},
		}
	}
	return nil
}

// TodosList implements Todos_list operation.
//
// List widgets.
//
// GET /widgets
func (h *TodoHandler) TodosList(ctx context.Context) (r *api.TodoList, _ error) {
	var todoList api.TodoList
	err := h.bundb.NewSelect().Model((*api.Todo)(nil)).Scan(ctx, &todoList.Items)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &todoList, nil
}

// TodosRead implements Todos_read operation.
//
// Read widgets.
//
// GET /widgets/{id}
func (h *TodoHandler) TodosRead(ctx context.Context, params api.TodosReadParams) (r *api.Todo, _ error) {
	var todo api.Todo
	err := h.bundb.NewSelect().Model((*api.Todo)(nil)).Where("id = ?", params.ID).Scan(ctx, &todo)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &todo, nil
}

// TodosUpdate implements Todos_update operation.
//
// Update a widget.
//
// PATCH /widgets/{id}
func (h *TodoHandler) TodosUpdate(ctx context.Context, req *api.TodoUpdate, params api.TodosUpdateParams) (r *api.Todo, _ error) {
	var todo api.Todo
	err := h.bundb.NewSelect().Model((*api.Todo)(nil)).Where("id = ?", params.ID).Scan(ctx, &todo)
	if err != nil {
		return nil, &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response:   api.Error{Message: "No records updated"},
		}
	}
	if req.Content.IsSet() {
		todo.SetContent(req.Content.Value)
	}
	if req.Done.IsSet() {
		todo.SetDone(req.Done.Value)
	}
	result, err := h.bundb.NewUpdate().Model(&todo).Where("id = ?", params.ID).Exec(context.Background())
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	if num, err := result.RowsAffected(); err != nil || num == 0 {
		return nil, &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response:   api.Error{Message: "No records updated"},
		}
	}
	return &todo, nil
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (TodoHandler) NewError(ctx context.Context, err error) (r *api.ErrorStatusCode) {
	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			Message: err.Error(),
		},
	}
}

func main() {
	conn, err := sql.Open(sqliteshim.ShimName, "file:./todo.sqlite?cache=shared")
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(1)

	bundb := bun.NewDB(conn, sqlitedialect.New())
	bundb.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	_, err = bundb.NewCreateTable().Model((*api.Todo)(nil)).IfNotExists().ModelTableExpr("todos").
		Exec(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	server, err := api.NewServer(&TodoHandler{
		bundb: bundb,
	})
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(":8888", server)
}
