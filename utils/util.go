package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, message string,  err error, status int) {
	log.Printf(fmt.Sprintf("%s: %%s", message), err)
	http.Error(w, http.StatusText(status), status)
}

type Page struct {
	Page    int
	Limit   int
	OrderBy []string
}

func NewPage(r *http.Request) Page {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}
	var orderBy []string
	orderBy = append(orderBy, r.URL.Query().Get("orderBy"))
	paginator := Page{Limit: limit, Page: page + 1, OrderBy: orderBy}
	log.Printf("vars: %#v", paginator)
	return paginator
}

type Filter struct {
	Key   string
	Value string
}

func NewFilters(r *http.Request, keys []string) []Filter {
	var filters []Filter
	for _, v := range keys {
		if r.URL.Query().Get(v) != "" {
			filters = append(filters, Filter{Key: v, Value: r.URL.Query().Get(v)})
		}
	}
	return filters
}
