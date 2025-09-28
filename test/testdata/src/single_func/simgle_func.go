package single_func

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
)

func insert(ctx context.Context, client *spanner.Client, isNow bool) error {
	var now time.Time
	if isNow {
		now = time.Now()
	} else {
		now = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	_, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.Insert( // want `Insert may use an argument that is a value from time.Now()`
			"Users",
			[]string{"name", "created_at"},
			[]interface{}{"Alice", now},
		),
	})

	return err
}
