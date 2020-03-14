package loader

import (
	"context"
	"fmt"

	"github.com/ghosv/open/meta"

	"github.com/ghosv/open/gate/client"
	pb "github.com/ghosv/open/plat/services/self/proto"
	"gopkg.in/nicksrandall/dataloader.v5"
)

type userLoader struct {
	service *client.MicroClient
}

func newUserLoader(s *client.MicroClient) dataloader.BatchFunc {
	return userLoader{s}.loadBatch
}

func (ldr userLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(keys)
		results = make([]*dataloader.Result, n)
	)

	res, err := ldr.service.SelfUser.BatchFind(ctx, &pb.BatchID{
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

// LoadUser Data
func LoadUser(ctx context.Context, key string) (*pb.UserInfo, error) {
	var user *pb.UserInfo

	ldr, err := extract(ctx, userLoaderKey)
	if err != nil {
		return nil, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(key))()
	if err != nil {
		return nil, err
	}
	user, ok := data.(*pb.UserInfo)
	if !ok {
		return nil, fmt.Errorf(fmtWrongType, user, data)
	}

	return user, nil
}

// LoadUsers Data
func LoadUsers(ctx context.Context, keys []string) ([]*pb.UserInfo, error) {
	var results []*pb.UserInfo
	ldr, err := extract(ctx, userLoaderKey)
	if err != nil {
		return nil, err
	}

	data, errs := ldr.LoadMany(ctx, dataloader.NewKeysFromStrings(keys))()
	results = make([]*pb.UserInfo, 0, len(keys))

	for i, d := range data {
		if errs != nil && errs[i] != nil {
			return nil, errs[i]
		}

		user, ok := d.(*pb.UserInfo)
		if !ok {
			return nil, fmt.Errorf(fmtWrongType, user, d)
		}

		results = append(results, user)
	}

	return results, nil
}
