import "@typespec/http";
import "@typespec/openapi";

using Http;
@service(#{ title: "Todo Service" })
namespace TodoService;

@doc("Represent a Todo item")
@example(#{ content: "Hello, World", id: 2 })
model Todo {
    @TypeSpec.OpenAPI.extension(
        "x-oapi-codegen-extra-tags",
        #{ bun: "id,pk,autoincrement" }
    )
    @doc("ID of the item")
    id: int32;

    @doc("Content text")
    @minLength(1)
    @maxLength(32)
    content: string;
}

@doc("Represent a list of Todo items")
model TodoList {
    @doc("Todo items")
    items: Todo[];
}

@error
model Error {
    code: int32;
    message: string;
}

@route("/todos")
@tag("Todos")
interface Todos {
    /** List todos */
    @get list(): TodoList | Error;
    /** Read todos */
    @get read(@path id: string): Todo | Error;
    /** Create a Todo */
    @post create(@body body: Todo): Todo | Error;
    /** Update a Todo */
    @patch update(@path id: string, @body body: Todo): Todo | Error;
    /** Delete a Todo */
    @delete delete(@path id: string): void | Error;
}
