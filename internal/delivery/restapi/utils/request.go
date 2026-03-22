package utils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func ParseIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	idStr := vars["id"]
	if idStr == "" {
		return 0, fmt.Errorf("id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return id, nil
}

func ParseAliasFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	alias := vars["alias"]
	if alias == "" {
		return "", fmt.Errorf("alias is required")
	}

	return alias, nil
}

func ParseFormData(r *http.Request, form interface{}) error {
	var decoder = schema.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		log.Printf("r.ParseForm: %s", err)
		return fmt.Errorf("r.ParseForm, %w", err)
	}
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(form, r.PostForm)
	return err
}
