package delete_banner

import (
	"context"
)

type BannerService interface {
	DeleteBanner(ctx context.Context, id *int64) error
}
