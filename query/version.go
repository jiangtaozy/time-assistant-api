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

type Version struct {
  VersionName string `json:"versionName"`
  VersionNumber string `json:"versionNumber"`
  VersionUrl string `json:"versionUrl"`
  VersionApkUrl string `json:"versionApkUrl"`
}

var VersionQuery = &graphql.Field{
  Type: graphql.NewObject(
    graphql.ObjectConfig{
      Name: "Version",
      Fields: graphql.Fields{
        "versionName": &graphql.Field{
          Type: graphql.String,
        },
        "versionNumber": &graphql.Field{
          Type: graphql.String,
        },
        "versionUrl": &graphql.Field{
          Type: graphql.String,
        },
        "versionApkUrl": &graphql.Field{
          Type: graphql.String,
        },
      },
    },
  ),
  Resolve: func(p graphql.ResolveParams) (interface{}, error) {
    versionName, _ := mutation.RedisClient.Get("versionName").Result()
    versionNumber, _ := mutation.RedisClient.Get("versionNumber").Result()
    versionUrl, _ := mutation.RedisClient.Get("versionUrl").Result()
    versionApkUrl, _ := mutation.RedisClient.Get("versionApkUrl").Result()
    version := Version{
      VersionName: versionName,
      VersionNumber: versionNumber,
      VersionUrl: versionUrl,
      VersionApkUrl: versionApkUrl,
    };
    return version, nil
  },
}
