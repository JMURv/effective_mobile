package validation

import (
	"github.com/JMURv/effectiveMobile/pkg/model"
)

func ValidateSong(req *model.Song) error {
	if req.Group == "" {
		return ErrMissingGroup
	}

	if req.Song == "" {
		return ErrMissingSong
	}

	return nil
}
