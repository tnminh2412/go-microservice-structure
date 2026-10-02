package cache

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	ctx := context.TODO()

	clientID := "demo"
	cacheKey := "cachePkg:unit_test"

	data := cacheData{
		Name: "UnitTest",
	}

	Convey("Set should success", t, func() {
		err := Set(ctx, clientID, cacheKey, data, time.Second*3)
		So(err, ShouldBeNil)
	})
	Convey("Get should success", t, func() {
		var getData *cacheData
		found, err := Get(ctx, clientID, cacheKey, &getData)
		So(err, ShouldBeNil)
		So(found, ShouldEqual, true)
		So(getData.Name, ShouldEqual, data.Name)
	})
}

type cacheData struct {
	Name string
}
