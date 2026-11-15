package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var ErrInvalidID = errors.New("invalid id")

func ParsePathID(r *http.Request, paramName string) (uint64, error) {
	raw := r.PathValue(paramName)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, ErrInvalidID
	}
	return id, nil
}

func JsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
}
