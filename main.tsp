import "@typespec/http";
import "@typespec/openapi";

using Http;
@service(#{ title: "Todo Service" })
namespace TodoService;

model Todo {
  @TypeSpec.OpenAPI.extension("x-oapi-codegen-extra-tags", #{
    bun: "id,pk,autoincrement"
  })
  id: int32;
  content: string;
}

model TodoList {
  items: Todo[];
}

@error
model Error {
  code: int32;
  message: string;
}

model AnalyzeResult {
  id: string;
  analysis: string;
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

  /** Analyze a Todo */
  @route("{id}/analyze") @post analyze(@path id: string): AnalyzeResult | Error;
}
