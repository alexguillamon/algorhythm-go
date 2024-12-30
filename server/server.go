package server

import (
	"algorhytm/language"
	"algorhytm/rhymes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

// TODO: make the handleDictionary rhyme take params for the get
// TODO: use the params to find the word in the dictionary and then feed that to they search

// TODO: plan and port the Rhyming service and searches
var decoder = schema.NewDecoder()

type Config struct {
	Host string
	Port string
}

func NewServer(
	config *Config,
	English *language.Language,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, English)

	var handler http.Handler = mux

	return handler
}

// addRoutes attaches all your endpoints to the provided *http.ServeMux.
func addRoutes(
	mux *http.ServeMux,
	English *language.Language,
) {
	mux.Handle("/{lang}/dictionary/rhymes", handleDictionaryRhymes(English))
	mux.Handle("/", http.NotFoundHandler())

}

type RhymesQuery struct {
	Word string `schema:"word"`
}

func handleDictionaryRhymes(lang *language.Language) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var query RhymesQuery
		if err := decodeQueryParams(r, &query); err != nil {
			_ = encode(w, r, http.StatusBadRequest, err)
			return
		}

		pronunciation, ok := (*lang.PhonemeDictionary)[query.Word]
		if !ok {
			_ = encode(w, r, int(http.StatusNotFound), "word not in the dictionary")
			return
		}
		sequence := lang.PhoneticAlphabet.FindStressSquence(pronunciation[0])
		// family, err := rhymes.FindFamily(lang, sequence)
		substractive, err := rhymes.FindAssonance(lang, sequence)
		if err != nil {
			_ = encode(w, r, http.StatusInternalServerError, err)
		}
		// result := lang.Trie.Search(sequence)
		err = encode(w, r, http.StatusOK, substractive)
		if err != nil {
			_ = encode(w, r, http.StatusInternalServerError, err)
		}
	})
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decodeQueryParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("unable to parse form: %v", err)
	}

	if err := decoder.Decode(dst, r.Form); err != nil {
		return fmt.Errorf("unable to decode query parameters: %v", err)
	}

	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// Validator is an object that can be validated.
type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}
