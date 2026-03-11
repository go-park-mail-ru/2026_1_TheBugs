package parse

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

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
