/*
 * Maintained by jemo from 2019.12.12 to now
 * Created by jemo on 2019.12.12 10:42:52
 * Main
 */

package main

import (
  "log"
  "fmt"
  "time"
  "strconv"
  "encoding/json"
  "net/http"
  "github.com/go-redis/redis/v7"
  "github.com/graphql-go/graphql"
)

type PostData struct {
  Query string `json:"query"`
  Variables map[string]interface{} `json:"variables"`
}

var port = ":5000"
var redisClient *redis.Client
var schema graphql.Schema

func main() {
  redisClient = redis.NewClient(&redis.Options{
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
  mutation := graphql.NewObject(graphql.ObjectConfig{
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
      "userRecordTimes": &graphql.Field{
        Type: graphql.String,
        Args: graphql.FieldConfigArgument{
          "id": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
          "lastDayRecordTimes": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.Int),
          },
        },
        Resolve: func(params graphql.ResolveParams) (interface{}, error) {
          id, _ := params.Args["id"].(string)
          lastDayRecordTimes, _ := params.Args["lastDayRecordTimes"].(int)
          offsetString, err := redisClient.Get(id).Result()
          offset, _ := strconv.ParseInt(offsetString, 0, 64)
          if err == redis.Nil {
            offset, _ = redisClient.BitCount("user", nil).Result()
            _ = redisClient.Set(id, offset, 0).Err()
            _ = redisClient.SetBit("user", offset, 1).Err()
          }
          now := time.Now()
          lastDay := now.AddDate(0, 0, -1)
          recordTimesBitmapKey := fmt.Sprintf("%d.%d.%d-%d", lastDay.Year(), lastDay.Month(), lastDay.Day(), lastDayRecordTimes)
          _ = redisClient.SetBit(recordTimesBitmapKey, offset, 1).Err()
          return "ok", nil
        },
      },
    },
  })
  schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: query,
    Mutation: mutation,
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
