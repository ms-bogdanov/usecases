package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
	"usecases/internal/service"
)

type Controller struct {
	responder Responder
	svc       service.Service
}

func NewController(responder Responder, svc service.Service) *Controller {
	return &Controller{
		responder: responder,
		svc:       svc,
	}
}

func (c Controller) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var s RequestAddressSearch
	json.NewDecoder(r.Body).Decode(&s)

	tmp, err := c.svc.Search(s.Query)
	if err != nil {
		c.responder.ErrorInternal(w, err)
	}

	c.responder.OutputJSON(w, tmp)
}

func (c Controller) GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	var s RequestAddressGeocode
	json.NewDecoder(r.Body).Decode(&s)

	tmp, err := c.svc.Geocode(s.Lat, s.Lng)
	if err != nil {
		c.responder.ErrorInternal(w, err)
	}

	c.responder.OutputJSON(w, tmp)
}

type MyController struct {
	responder Responder
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type Respond struct {
	log *zap.Logger
	godecoder.Decoder
}

func NewResponder(decoder godecoder.Decoder, logger *zap.Logger) Responder {
	return &Respond{log: logger, Decoder: decoder}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type",
		"application/json;charset=utf-8")
	if err := r.Encode(w, responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Info("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Info("http response unauthorized status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
