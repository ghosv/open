package handler

import (
	"context"
	"strings"

	"github.com/ghosv/open/plat/conf"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getMongoRepo(ctx context.Context, collection string) *mgo.Collection {
	s, _ := conf.GetMongoRepo(ctx)
	c := strings.ToLower(collection)
	return s.DB("mrms").C(c)
}

func setStr(m bson.M, k string, v string) {
	if v != "" {
		m[k] = v
	}
}

func setArr(m bson.M, k string, v []string) {
	if v != nil {
		m[k] = v
	}
}

func setTime(m bson.M, k string, v *timestamp.Timestamp) {
	if v != nil {
		m[k] = timePb2Mongo(v)
	}
}

func timePb2Mongo(v *timestamp.Timestamp) bson.MongoTimestamp {
	return bson.MongoTimestamp(v.Seconds*1000 + int64(v.Nanos)/1_000_000)
}

func timeMongo2Pb(v bson.MongoTimestamp) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: int64(v / 1000),
		Nanos:   int32(v%1000) * 1_000_000,
	}
}
