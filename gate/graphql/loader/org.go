package loader

import (
	"fmt"

	"gopkg.in/nicksrandall/dataloader.v5"

	"context"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

type orgLoader struct {
	service *client.MicroClient
}

func neworgLoader(s *client.MicroClient) dataloader.BatchFunc {
	return orgLoader{s}.loadBatch
}

func (ldr orgLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(keys)
		results = make([]*dataloader.Result, n)
	)

	res, err := ldr.service.SelfOrg.BatchFind(ctx, &pb.BatchID{
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

// LoadOrg Data
func LoadOrg(ctx context.Context, key string) (*pb.OrgInfo, error) {
	var org *pb.OrgInfo

	ldr, err := extract(ctx, orgLoaderKey)
	if err != nil {
		return nil, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(key))()
	if err != nil {
		return nil, err
	}
	org, ok := data.(*pb.OrgInfo)
	if !ok {
		return nil, fmt.Errorf(fmtWrongType, org, data)
	}

	return org, nil
}

// LoadOrgs Data
func LoadOrgs(ctx context.Context, keys []string) ([]*pb.OrgInfo, error) {
	var results []*pb.OrgInfo
	ldr, err := extract(ctx, orgLoaderKey)
	if err != nil {
		return nil, err
	}

	data, errs := ldr.LoadMany(ctx, dataloader.NewKeysFromStrings(keys))()
	results = make([]*pb.OrgInfo, 0, len(keys))

	for i, d := range data {
		if errs != nil && errs[i] != nil {
			return nil, errs[i]
		}

		org, ok := d.(*pb.OrgInfo)
		if !ok {
			return nil, fmt.Errorf(fmtWrongType, org, d)
		}

		results = append(results, org)
	}

	return results, nil
}
