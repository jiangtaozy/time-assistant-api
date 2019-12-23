/*
 * Maintained by jemo from 2019.12.23 to now
 * Created by jemo on 2019.12.23 16:12:28
 * Version
 */

package query

import (
  "github.com/graphql-go/graphql"
  "github.com/jiangtaozy/time-assistant-api/mutation"
)

var VersionQuery = &graphql.Field{
  Type: graphql.String,
  Resolve: func(p graphql.ResolveParams) (interface{}, error) {
    version, _ := mutation.RedisClient.Get("version").Result()
    return version, nil
  },
}
