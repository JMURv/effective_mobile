package http

import (
	"encoding/json"
	"errors"
	"github.com/JMURv/effectiveMobile/internal/ctrl"
	"github.com/JMURv/effectiveMobile/internal/hdl"
	"github.com/JMURv/effectiveMobile/internal/validation"
	"github.com/JMURv/effectiveMobile/pkg/model"
	utils "github.com/JMURv/effectiveMobile/pkg/utils/http"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

// ListSongs
// @Summary Список песен
// @Description Получить список песен с возможностью фильтрации и пагинации
// @Tags songs
// @Param page query int false "Номер страницы" default(1)
// @Param size query int false "Размер страницы" default(40)
// @Param group query string false "Фильтр по имени группы"
// @Param song query string false "Фильтр по названию песни"
// @Param release_date query string false "Фильтр по дате релиза"
// @Param min_release_date query string false "Фильтр по дате релиза (минимальная)"
// @Param release_date query string false "Фильтр по дате релиза (максимальная)"
// @Success 200 {object} model.PaginatedSongs "Список песен с пагинацией"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/songs [get]
func (h *Handler) ListSongs(w http.ResponseWriter, r *http.Request) {
	const op = "songs.ListSongs.hdl"

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 40
	}

	filters := utils.ParseFiltersByURL(r)
	res, err := h.ctrl.ListSongs(r.Context(), page, size, filters)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessPaginatedResponse(w, http.StatusOK, res)
}

// GetSong
// @Summary Получить песню по ID
// @Description Получить информацию о песне и её тексте(пагинация)
// @Tags songs
// @Param id path int true "ID песни"
// @Param page query int false "Номер куплета для пагинации текста песни" default(1)
// @Param size query int false "Размер куплета для пагинации текста песни" default(40)
// @Success 200 {object} model.PaginatedSongs "Детали песни с пагинированным текстом"
// @Failure 400 {object} utils.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} utils.ErrorResponse "Песня не найдена"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/songs/{id} [get]
func (h *Handler) GetSong(w http.ResponseWriter, r *http.Request) {
	const op = "songs.GetSong.hdl"

	id, err := strconv.ParseUint(
		strings.TrimPrefix(r.URL.Path, "/api/songs/"), 10, 64,
	)
	if err != nil {
		zap.L().Debug(
			"failed to extract param",
			zap.Error(err), zap.String("op", op),
		)
		utils.ErrResponse(w, http.StatusBadRequest, hdl.ErrMissingSongID)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 40
	}

	res, err := h.ctrl.GetSong(r.Context(), id, page, size)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		utils.ErrResponse(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
		return
	}

	utils.SuccessPaginatedResponse(w, http.StatusOK, res)
}

type CreateSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// CreateSong
// @Summary Добавить новую песню
// @Description Добавить новую песню
// @Tags songs
// @Accept json
// @Param song body CreateSongRequest true "Данные новой песни"
// @Success 200 {object} int "ID добавленной песни"
// @Failure 400 {object} utils.ErrorResponse "Ошибка валидации или декодирования запроса"
// @Failure 409 {object} utils.ErrorResponse "Песня уже существует"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/songs [post]
func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	const op = "songs.CreateSong.hdl"

	req := &model.Song{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		zap.L().Debug(
			"failed to decode request",
			zap.Error(err), zap.String("op", op),
		)
		utils.ErrResponse(w, http.StatusBadRequest, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidateSong(req); err != nil {
		zap.L().Debug(
			"failed to validate request",
			zap.Error(err), zap.String("op", op),
			zap.Any("req", req),
		)
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.ctrl.CreateSong(r.Context(), req)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		utils.ErrResponse(w, http.StatusConflict, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, res)
}

// UpdateSong
// @Summary Обновить информацию о песне
// @Description Обновить существующую песню по ID
// @Tags songs
// @Accept json
// @Param id path int true "ID песни"
// @Param song body model.Song true "Обновленные данные песни"
// @Success 200 {object} string "OK"
// @Failure 400 {object} utils.ErrorResponse "Ошибка валидации или декодирования запроса"
// @Failure 404 {object} utils.ErrorResponse "Песня не найдена"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/songs/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	const op = "songs.UpdateSong.hdl"

	songID, err := strconv.ParseUint(
		strings.TrimPrefix(r.URL.Path, "/api/songs/"), 10, 64,
	)
	if err != nil {
		zap.L().Debug(
			"failed to extract param",
			zap.Error(err), zap.String("op", op),
		)
		utils.ErrResponse(w, http.StatusBadRequest, hdl.ErrMissingSongID)
		return
	}

	req := &model.Song{ID: songID}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		zap.L().Debug(
			"failed to decode request",
			zap.Error(err), zap.String("op", op),
		)
		utils.ErrResponse(w, http.StatusBadRequest, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidateSong(req); err != nil {
		zap.L().Debug(
			"failed to validate request",
			zap.Error(err), zap.String("op", op),
			zap.Any("req", req),
		)
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.ctrl.UpdateSong(r.Context(), req)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		utils.ErrResponse(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "OK")
}

// DeleteSong
// @Summary Удалить песню
// @Description Удалить существующую песню по ID
// @Tags songs
// @Param id path int true "ID песни"
// @Success 204 {object} string "Песня успешно удалена"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	const op = "songs.DeleteSong.hdl"

	songID, err := strconv.ParseUint(
		strings.TrimPrefix(r.URL.Path, "/api/songs/"), 10, 64,
	)
	if err != nil {
		zap.L().Debug(
			"failed to extract param",
			zap.Error(err), zap.String("op", op),
		)
		utils.ErrResponse(w, http.StatusBadRequest, hdl.ErrMissingSongID)
		return
	}

	err = h.ctrl.DeleteSong(r.Context(), songID)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		utils.ErrResponse(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, http.StatusNoContent, "OK")
}
