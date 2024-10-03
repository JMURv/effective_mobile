package external

import (
	"encoding/json"
	"fmt"
	"github.com/JMURv/effectiveMobile/internal/ctrl"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	port int
}

func New(port int) *Controller {
	return &Controller{
		port: port,
	}
}

func (c *Controller) FetchSongDetail(group, song string) (*model.SongDetail, error) {
	get, err := http.Get(fmt.Sprintf("http://localhost:%d/info?group=%s&song=%s", c.port, group, song))
	if err != nil {
		zap.L().Error("failed to fetch song detail", zap.Error(err))
		return nil, err
	}

	switch get.StatusCode {
	case http.StatusBadRequest:
		return nil, ctrl.ErrBadExtReq
	case http.StatusInternalServerError:
		return nil, ctrl.ExtSrvErr
	case http.StatusOK:
		defer get.Body.Close()
		res := &model.SongDetail{}
		if err := json.NewDecoder(get.Body).Decode(res); err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, ctrl.ErrExtUnreachable
	}
}
