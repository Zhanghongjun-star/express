package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shunfeng-miniprogram/internal/biz"

	kratoshttp "github.com/go-kratos/kratos/v3/transport/http"
)

type AddrRequest struct {
	AddrID        int64  `json:"addr_id"`
	UserID        int64  `json:"user_id"`
	AddrType      int32  `json:"addr_type"`
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddr    string `json:"detail_addr"`
}

type AddressService struct {
	uc *biz.AddressUsecase
}

func NewAddressService(uc *biz.AddressUsecase) *AddressService {
	return &AddressService{uc: uc}
}

func (s *AddressService) RegisterRoutes(srv *kratoshttp.Server) {
	srv.Handle("/api/user/address/list", http.HandlerFunc(s.List))
	srv.Handle("/api/user/address/add", http.HandlerFunc(s.Create))
	srv.Handle("/api/user/address/edit", http.HandlerFunc(s.Update))
	srv.Handle("/api/user/address/del", http.HandlerFunc(s.Delete))
	srv.Handle("/api/user/address/batch/del", http.HandlerFunc(s.BatchDelete))
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *AddressService) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	addrs, err := s.uc.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"code": 0, "data": addrs})
}

func (s *AddressService) Create(w http.ResponseWriter, r *http.Request) {
	var req AddrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	addr := &biz.Address{
		UserID:        req.UserID,
		AddrType:      biz.AddressType(req.AddrType),
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		DetailAddr:    req.DetailAddr,
	}
	result, err := s.uc.Create(r.Context(), addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"code": 0, "data": result})
}

func (s *AddressService) Update(w http.ResponseWriter, r *http.Request) {
	var req AddrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	addr := &biz.Address{
		AddrID:        req.AddrID,
		AddrType:      biz.AddressType(req.AddrType),
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		DetailAddr:    req.DetailAddr,
	}
	result, err := s.uc.Update(r.Context(), addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"code": 0, "data": result})
}

func (s *AddressService) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	id, _ := strconv.ParseInt(r.URL.Query().Get("addr_id"), 10, 64)
	if err := s.uc.Delete(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"code": 0, "msg": "ok"})
}

func (s *AddressService) BatchDelete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64   `json:"user_id"`
		IDs    []int64 `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := s.uc.BatchDelete(r.Context(), req.IDs, req.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"code": 0, "msg": "ok"})
}
