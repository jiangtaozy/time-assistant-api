/*
 * Maintained by jemo from 2019.12.13 to now
 * Created by jemo on 2019.12.13 13:46:00
 * User record times
 */

package query

import (
  "fmt"
  "math"
  "time"
  "github.com/graphql-go/graphql"
  "github.com/jiangtaozy/time-assistant-api/mutation"
)

var location, err = time.LoadLocation("Local")
var startDate = time.Date(2019, time.December, 11, 0, 0, 0, 0, location);

type RecordNumber struct {
  Date string `json:"date"`
  UserNumber int64 `json"userNumber"`
}

type RecordTimes struct {
  Times int64 `json"times"`
  RecordData []RecordNumber `json"recordData"`
}

var RecordNumberType = graphql.NewObject(
  graphql.ObjectConfig{
    Name: "RecordNumber",
    Fields: graphql.Fields{
      "date": &graphql.Field{
        Type: graphql.String,
      },
      "userNumber": &graphql.Field{
        Type: graphql.Int,
      },
    },
  },
)

var RecordTimesType = graphql.NewObject(
  graphql.ObjectConfig{
    Name: "RecordTimes",
    Fields: graphql.Fields{
      "times": &graphql.Field{
        Type: graphql.Int,
      },
      "recordData": &graphql.Field{
        Type: graphql.NewList(RecordNumberType),
      },
    },
  },
)

var UserRecordTimesQuery = &graphql.Field{
  Type: graphql.NewList(RecordTimesType),
  Resolve: func(p graphql.ResolveParams) (interface{}, error) {
    now := time.Now()
    duration := now.Sub(startDate)
    days := int(math.Floor(duration.Hours() / 24))
    userRecordTimes := []RecordTimes{}
    for i := 0; i < 30; i++ {
      recordData := []RecordNumber{}
      var totalRecordNumberOfTimes int64 = 0
      for j := 0; j < days; j++ {
        date := startDate.AddDate(0, 0, j)
        recordTimesBitmapKey := fmt.Sprintf("%d.%d.%d-%d", date.Year(), date.Month(), date.Day(), i)
        times, _ := mutation.RedisClient.BitCount(recordTimesBitmapKey, nil).Result()
        totalRecordNumberOfTimes += times
        recordData = append(recordData, RecordNumber{
          Date: date.String(),
          UserNumber: times,
        })
      }
      recordTimes := RecordTimes{
        Times: int64(i),
        RecordData: recordData,
      }
      if totalRecordNumberOfTimes > 0 {
        userRecordTimes = append(userRecordTimes, recordTimes)
      }
    }
    return userRecordTimes, nil
  },
}
