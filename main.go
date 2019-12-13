/*
 * Maintained by jemo from 2019.12.12 to now
 * Created by jemo on 2019.12.12 10:42:52
 * Main
 */

package main

import (
  "log"
  "encoding/json"
  "net/http"
  "github.com/go-redis/redis/v7"
  "github.com/graphql-go/graphql"
  "github.com/jiangtaozy/time-assistant-api/mutation"
)

type PostData struct {
  Query string `json:"query"`
  Variables map[string]interface{} `json:"variables"`
}

var port = ":5000"
var schema graphql.Schema

func main() {
  mutation.RedisClient = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    DB: 1,
  });
  query := graphql.NewObject(graphql.ObjectConfig{
    Name: "query",
    Fields: graphql.Fields{
      "hello": &graphql.Field{
        Type: graphql.String,
        Resolve: func(p graphql.ResolveParams) (interface{}, error) {
          return "world", nil
        },
      },
    },
  })
  rootMutation := graphql.NewObject(graphql.ObjectConfig{
    Name: "mutaion",
    Fields: graphql.Fields{
      "create": &graphql.Field{
        Type: graphql.String,
        Args: graphql.FieldConfigArgument{
          "text": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
        },
        Resolve: func(params graphql.ResolveParams) (interface{}, error) {
          text, _ := params.Args["text"].(string)
          return text, nil
        },
      },
      "userRecordTimes": mutation.UserRecordTimesMutation,
    },
  })
  schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: query,
    Mutation: rootMutation,
  })
  log.Println("listen at ", port)
  http.HandleFunc("/", handle)
  log.Fatal(http.ListenAndServe(port, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var data PostData
  err := decoder.Decode(&data)
  if err != nil {
    log.Println("HandleDecodeError, err: ", err)
    panic(err)
  }
  res := graphql.Do(graphql.Params{
    Schema: schema,
    RequestString: data.Query,
    VariableValues: data.Variables,
  })
  if len(res.Errors) > 0 {
    log.Printf("HandleResError, res.Errors: %v\n", res.Errors)
  }
  json.NewEncoder(w).Encode(res)
}
