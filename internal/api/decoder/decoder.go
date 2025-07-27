package decoder

import (
	"encoding/json"
	"fmt"
	"io"

	"go.uber.org/zap"

	"effmob/internal/api"
)

var (
	ErrCannotDecodeJSON = fmt.Errorf("DecodeRequest: cannot decode body")
)

func DecodeRequest(logger *zap.Logger, body io.ReadCloser) (*api.Subscription, error) {
	res := &api.Subscription{}

	err := json.NewDecoder(body).Decode(res)
	if err != nil {
		logger.Error(ErrCannotDecodeJSON.Error(), zap.Error(err))
		return nil, fmt.Errorf("%w: %s", ErrCannotDecodeJSON, err.Error())
	}

	return res, nil
}
