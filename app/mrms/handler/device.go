package handler

import (
	"context"

	pb "github.com/ghosv/open/app/mrms/proto"
	metaPb "github.com/ghosv/open/meta/proto"
	"gopkg.in/mgo.v2/bson"
)

// Device Handler
type Device struct{}

type wrapDevice struct {
	ID bson.ObjectId `bson:"_id"`

	Name string
	Type string

	Owner string
}

func (w *wrapDevice) From(d *pb.DeviceInfo) *wrapDevice {
	id := bson.NewObjectId()
	if d.ID != "" {
		id = bson.ObjectIdHex(d.ID)
	}
	return &wrapDevice{
		ID: id,

		Name: d.Name,
		Type: d.Type,

		Owner: d.Owner,
	}
}

func (w *wrapDevice) To() *pb.DeviceInfo {
	return &pb.DeviceInfo{
		ID: w.ID.Hex(),

		Name: w.Name,
		Type: w.Type,

		Owner: w.Owner,
	}
}

// === Query ===

// BatchFind Device
func (h *Device) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.DeviceMap) error {
	c := getMongoRepo(ctx, "Device")
	ids := make([]bson.ObjectId, 0, len(req.UUID))
	for _, v := range req.UUID {
		ids = append(ids, bson.ObjectIdHex(v))
	}
	iter := c.FindId(bson.M{"$in": ids}).Iter()

	data := make(map[string]*pb.DeviceInfo, len(req.UUID))
	var result wrapDevice
	for iter.Next(&result) {
		d := result.To()
		data[d.ID] = d
	}
	res.Data = data
	return nil
}

// Search Device
func (h *Device) Search(ctx context.Context, req *pb.SearchForm, res *pb.DeviceList) error {
	word, page, size := req.Word, req.Page, req.Size
	c := getMongoRepo(ctx, "Device")
	t := c.Find(bson.M{
		"name": bson.RegEx{Pattern: word, Options: "i"},
	})
	total, _ := t.Count()
	res.Total = int32(total)

	iter := t.Skip(int((page - 1) * size)).Limit(int(size)).Iter()
	data := make([]*pb.DeviceInfo, 0, size)
	var result wrapDevice
	for iter.Next(&result) {
		d := result.To()
		data = append(data, d)
	}
	res.List = data
	return nil
}

// === Mutation ===

// Create Device
func (h *Device) Create(ctx context.Context, req *pb.DeviceModify, res *pb.DeviceInfo) error {
	c := getMongoRepo(ctx, "Device")
	d := new(wrapDevice).From(req.Info)
	if err := c.Insert(d); err != nil {
		return err
	}
	*res = *d.To()
	return nil
}

// Delete Device
func (h *Device) Delete(ctx context.Context, req *pb.DeviceModify, res *metaPb.None) error {
	c := getMongoRepo(ctx, "Device")
	return c.RemoveId(bson.ObjectIdHex(req.Info.ID))
}

// Update Device
func (h *Device) Update(ctx context.Context, req *pb.DeviceModify, res *pb.DeviceInfo) error {
	c := getMongoRepo(ctx, "Device")

	d := bson.M{}
	setStr(d, "name", req.Info.Name)
	setStr(d, "type", req.Info.Type)

	setStr(d, "owner", req.Info.Owner)

	*res = *req.Info // TODO: get full
	return c.UpdateId(bson.ObjectIdHex(req.Info.ID), bson.M{"$set": d})
}
