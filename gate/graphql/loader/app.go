package loader

import (
	"fmt"

	"gopkg.in/nicksrandall/dataloader.v5"

	"context"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

type appLoader struct {
	service *client.MicroClient
}

func newAppLoader(s *client.MicroClient) dataloader.BatchFunc {
	return appLoader{s}.loadBatch
}

func (ldr appLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(keys)
		results = make([]*dataloader.Result, n)
	)

	res, err := ldr.service.SelfApp.BatchFind(ctx, &pb.BatchID{
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

// LoadApp Data
func LoadApp(ctx context.Context, key string) (*pb.AppInfo, error) {
	var app *pb.AppInfo

	ldr, err := extract(ctx, appLoaderKey)
	if err != nil {
		return nil, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(key))()
	if err != nil {
		return nil, err
	}
	app, ok := data.(*pb.AppInfo)
	if !ok {
		return nil, fmt.Errorf(fmtWrongType, app, data)
	}

	return app, nil
}

// LoadApps Data
func LoadApps(ctx context.Context, keys []string) ([]*pb.AppInfo, error) {
	var results []*pb.AppInfo
	ldr, err := extract(ctx, appLoaderKey)
	if err != nil {
		return nil, err
	}

	data, errs := ldr.LoadMany(ctx, dataloader.NewKeysFromStrings(keys))()
	results = make([]*pb.AppInfo, 0, len(keys))

	for i, d := range data {
		if errs != nil && errs[i] != nil {
			return nil, errs[i]
		}

		app, ok := d.(*pb.AppInfo)
		if !ok {
			return nil, fmt.Errorf(fmtWrongType, app, d)
		}

		results = append(results, app)
	}

	return results, nil
}
