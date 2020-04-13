package handler

import (
	"context"

	pb "github.com/ghosv/open/app/mrms/proto"
	metaPb "github.com/ghosv/open/meta/proto"
	"gopkg.in/mgo.v2/bson"
)

// Meeting Handler
type Meeting struct{}

type wrapMeeting struct {
	ID bson.ObjectId `bson:"_id"`

	Name      string
	Desc      string
	StartTime bson.MongoTimestamp
	EndTime   bson.MongoTimestamp

	Room  string
	Host  string
	Users []string
}

func (w *wrapMeeting) From(d *pb.MeetingInfo) *wrapMeeting {
	id := bson.NewObjectId()
	if d.ID != "" {
		id = bson.ObjectIdHex(d.ID)
	}
	return &wrapMeeting{
		ID: id,

		Name:      d.Name,
		Desc:      d.Desc,
		StartTime: timePb2Mongo(d.StartTime),
		EndTime:   timePb2Mongo(d.EndTime),

		Room:  d.Room,
		Host:  d.Host,
		Users: d.Users,
	}
}

func (w *wrapMeeting) To() *pb.MeetingInfo {
	return &pb.MeetingInfo{
		ID: w.ID.Hex(),

		Name:      w.Name,
		Desc:      w.Desc,
		StartTime: timeMongo2Pb(w.StartTime),
		EndTime:   timeMongo2Pb(w.EndTime),

		Room:  w.Room,
		Host:  w.Host,
		Users: w.Users,
	}
}

// === Query ===

// BatchFind Meeting
func (h *Meeting) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.MeetingMap) error {
	c := getMongoRepo(ctx, "Meeting")
	ids := make([]bson.ObjectId, 0, len(req.UUID))
	for _, v := range req.UUID {
		ids = append(ids, bson.ObjectIdHex(v))
	}
	iter := c.FindId(bson.M{"$in": ids}).Iter()

	data := make(map[string]*pb.MeetingInfo, len(req.UUID))
	var result wrapMeeting
	for iter.Next(&result) {
		d := result.To()
		data[d.ID] = d
	}
	res.Data = data
	return nil
}

// Search Meeting
func (h *Meeting) Search(ctx context.Context, req *pb.SearchForm, res *pb.MeetingList) error {
	word, page, size := req.Word, req.Page, req.Size
	c := getMongoRepo(ctx, "Meeting")
	t := c.Find(bson.M{
		"name": bson.RegEx{Pattern: word, Options: "i"},
	})
	total, _ := t.Count()
	res.Total = int32(total)

	iter := t.Skip(int((page - 1) * size)).Limit(int(size)).Iter()
	data := make([]*pb.MeetingInfo, 0, size)
	var result wrapMeeting
	for iter.Next(&result) {
		d := result.To()
		data = append(data, d)
	}
	res.List = data
	return nil
}

// === Mutation ===

// Create Meeting
func (h *Meeting) Create(ctx context.Context, req *pb.MeetingModify, res *pb.MeetingInfo) error {
	c := getMongoRepo(ctx, "Meeting")
	d := new(wrapMeeting).From(req.Info)
	if err := c.Insert(d); err != nil {
		return err
	}
	*res = *d.To()
	return nil
}

// Delete Meeting
func (h *Meeting) Delete(ctx context.Context, req *pb.MeetingModify, res *metaPb.None) error {
	c := getMongoRepo(ctx, "Meeting")
	return c.RemoveId(bson.ObjectIdHex(req.Info.ID))
}

// Update Meeting
func (h *Meeting) Update(ctx context.Context, req *pb.MeetingModify, res *pb.MeetingInfo) error {
	c := getMongoRepo(ctx, "Meeting")

	d := bson.M{}
	setStr(d, "name", req.Info.Name)
	setStr(d, "desc", req.Info.Desc)
	setTime(d, "startTime", req.Info.StartTime)
	setTime(d, "endTime", req.Info.EndTime)

	setStr(d, "room", req.Info.Room)
	setStr(d, "host", req.Info.Host)   // 主持人
	setArr(d, "users", req.Info.Users) // TODO: Append/Delete?

	*res = *req.Info // TODO: get full
	return c.UpdateId(bson.ObjectIdHex(req.Info.ID), bson.M{"$set": d})
}
