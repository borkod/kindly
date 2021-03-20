package pkg

import (
	"context"
)

// Check function checks if the packages passed in args are available TODO variadic function
func (k Kindly) Check(ctx context.Context, n string) (ks KindlyStruct, err error) {

	_, yc, err := k.getValidYConfig(ctx, n, false, false)
	if err != nil {
		return yc, err
	}

	return yc, nil
}
