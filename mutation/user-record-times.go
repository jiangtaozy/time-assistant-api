/*
 * Maintained by jemo from 2019.12.13 to now
 * Created by jemo on 2019.12.13 12:42:25
 * User record times
 */

package mutation

import (
  "fmt"
  "time"
  "strconv"
  "github.com/graphql-go/graphql"
  "github.com/go-redis/redis/v7"
)

var RedisClient *redis.Client

var UserRecordTimesMutation = &graphql.Field{
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
    offsetString, err := RedisClient.Get(id).Result()
    offset, _ := strconv.ParseInt(offsetString, 0, 64)
    if err == redis.Nil {
      offset, _ = RedisClient.BitCount("user", nil).Result()
      _ = RedisClient.Set(id, offset, 0).Err()
      _ = RedisClient.SetBit("user", offset, 1).Err()
    }
    now := time.Now()
    lastDay := now.AddDate(0, 0, -1)
    recordTimesBitmapKey := fmt.Sprintf("%d.%d.%d-%d", lastDay.Year(), lastDay.Month(), lastDay.Day(), lastDayRecordTimes)
    _ = RedisClient.SetBit(recordTimesBitmapKey, offset, 1).Err()
    return "ok", nil
  },
}
