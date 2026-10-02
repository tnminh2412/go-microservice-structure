package validator

import (
	"context"
	service_error "go-pkg/helper/error_helper"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/date"
)

type Interval struct {
	FromDate date.Date
	ToDate   date.Date
}

func (v Validator) Interval(ctx context.Context, clientID string, fromDate, toDate *timestamppb.Timestamp) (*Interval, error) {
	if fromDate.IsValid() == false {
		msg := "FromDate is require"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if toDate.IsValid() == false {
		msg := "ToDate is require"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	cf, err := v.GlobalConf.GetClientConf(ctx, clientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "", err)
	}
	loc, err := cf.Location()
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "", err)
	}

	from, err := date.FromTime(fromDate.AsTime(), loc)
	if err != nil {
		msg := "FromDate is invalid"
		return nil, error_helper.New(service_error.InvalidArgument, msg, err)
	}
	to, err := date.FromTime(toDate.AsTime().AddDate(0, 0, 1), loc)
	if err != nil {
		msg := "ToDate is invalid"
		return nil, error_helper.New(service_error.InvalidArgument, msg, err)
	}
	if !to.After(from) {
		msg := "ToDate must be after FromDate"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	i := &Interval{
		FromDate: from,
		ToDate:   to,
	}
	return i, nil
}
