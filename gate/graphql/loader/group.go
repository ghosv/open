package loader

import (
	"context"
	"fmt"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/self/proto"
	"gopkg.in/nicksrandall/dataloader.v5"
)

type groupLoader struct {
	service *client.MicroClient
}

func newgroupLoader(s *client.MicroClient) dataloader.BatchFunc {
	return groupLoader{s}.loadBatch
}

func (ldr groupLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(keys)
		results = make([]*dataloader.Result, n)
	)

	res, err := ldr.service.SelfGroup.BatchFind(ctx, &pb.BatchID{
		UUID: keys.Keys(),
	})
	if err != nil {
		return nil
	}
	for i, key := range keys {
		v := res.Data[key.String()]
		if v == nil {
			results[i] = &dataloader.Result{
				Error: meta.ErrRepoRecordNotFound,
			}
		} else {
			results[i] = &dataloader.Result{
				Data: v,
			}
		}
	}

	return results
}

// LoadGroup Data
func LoadGroup(ctx context.Context, key string) (*pb.GroupInfo, error) {
	var group *pb.GroupInfo

	ldr, err := extract(ctx, groupLoaderKey)
	if err != nil {
		return nil, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(key))()
	if err != nil {
		return nil, err
	}
	group, ok := data.(*pb.GroupInfo)
	if !ok {
		return nil, fmt.Errorf(fmtWrongType, group, data)
	}

	return group, nil
}

// LoadGroups Data
func LoadGroups(ctx context.Context, keys []string) ([]*pb.GroupInfo, error) {
	var results []*pb.GroupInfo
	ldr, err := extract(ctx, groupLoaderKey)
	if err != nil {
		return nil, err
	}

	data, errs := ldr.LoadMany(ctx, dataloader.NewKeysFromStrings(keys))()
	results = make([]*pb.GroupInfo, 0, len(keys))

	for i, d := range data {
		if errs != nil && errs[i] != nil {
			return nil, errs[i]
		}

		group, ok := d.(*pb.GroupInfo)
		if !ok {
			return nil, fmt.Errorf(fmtWrongType, group, d)
		}

		results = append(results, group)
	}

	return results, nil
}
