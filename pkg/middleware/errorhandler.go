package middleware

import "context"

func DefaultClientErrorHandler(ctx context.Context, err error) error {
	return err
}
