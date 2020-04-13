package handler

import (
	"context"

	pb "github.com/ghosv/open/app/mrms/proto"
	metaPb "github.com/ghosv/open/meta/proto"
	"gopkg.in/mgo.v2/bson"
)

// Room Handler
type Room struct{}

type wrapRoom struct {
	ID bson.ObjectId `bson:"_id"`

	Name string
	Addr string

	Devices []string
}

func (w *wrapRoom) From(d *pb.RoomInfo) *wrapRoom {
	id := bson.NewObjectId()
	if d.ID != "" {
		id = bson.ObjectIdHex(d.ID)
	}
	return &wrapRoom{
		ID: id,

		Name: d.Name,
		Addr: d.Addr,

		Devices: d.Devices,
	}
}

func (w *wrapRoom) To() *pb.RoomInfo {
	return &pb.RoomInfo{
		ID: w.ID.Hex(),

		Name: w.Name,
		Addr: w.Addr,

		Devices: w.Devices,
	}
}

// === Query ===

// BatchFind Room
func (h *Room) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.RoomMap) error {
	c := getMongoRepo(ctx, "Room")
	ids := make([]bson.ObjectId, 0, len(req.UUID))
	for _, v := range req.UUID {
		ids = append(ids, bson.ObjectIdHex(v))
	}
	iter := c.FindId(bson.M{"$in": ids}).Iter()

	data := make(map[string]*pb.RoomInfo, len(req.UUID))
	var result wrapRoom
	for iter.Next(&result) {
		d := result.To()
		data[d.ID] = d
	}
	res.Data = data
	return nil
}

// Search Room
func (h *Room) Search(ctx context.Context, req *pb.SearchForm, res *pb.RoomList) error {
	word, page, size := req.Word, req.Page, req.Size
	c := getMongoRepo(ctx, "Room")
	t := c.Find(bson.M{
		"name": bson.RegEx{Pattern: word, Options: "i"},
	})
	total, _ := t.Count()
	res.Total = int32(total)

	iter := t.Skip(int((page - 1) * size)).Limit(int(size)).Iter()
	data := make([]*pb.RoomInfo, 0, size)
	var result wrapRoom
	for iter.Next(&result) {
		d := result.To()
		data = append(data, d)
	}
	res.List = data
	return nil
}

// === Mutation ===

// Create Room
func (h *Room) Create(ctx context.Context, req *pb.RoomModify, res *pb.RoomInfo) error {
	c := getMongoRepo(ctx, "Room")
	d := new(wrapRoom).From(req.Info)
	if err := c.Insert(d); err != nil {
		return err
	}
	*res = *d.To()
	return nil
}

// Delete Room
func (h *Room) Delete(ctx context.Context, req *pb.RoomModify, res *metaPb.None) error {
	c := getMongoRepo(ctx, "Room")
	return c.RemoveId(bson.ObjectIdHex(req.Info.ID))
}

// Update Room
func (h *Room) Update(ctx context.Context, req *pb.RoomModify, res *pb.RoomInfo) error {
	c := getMongoRepo(ctx, "Room")

	d := bson.M{}
	setStr(d, "name", req.Info.Name)
	setStr(d, "addr", req.Info.Addr)

	setArr(d, "devices", req.Info.Devices) // TODO: Append/Delete?

	*res = *req.Info // TODO: get full
	return c.UpdateId(bson.ObjectIdHex(req.Info.ID), bson.M{"$set": d})
}
