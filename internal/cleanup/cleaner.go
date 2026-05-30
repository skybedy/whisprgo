package cleanup

import "context"

type Cleaner interface {
	Clean(ctx context.Context, input string, model string, prompt string) (string, error)
}
