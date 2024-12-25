package types

import "fmt"

// Headers that preserve case and order
type Headers struct {
	m    map[string]string
	keys []string
}

type headerValue func(h *Headers)

func Header(k string, v string) func(headers *Headers) {
	return func(headers *Headers) {
		headers.Set(k, v)
	}
}

func NewHeaders(values ...headerValue) *Headers {
	headers := &Headers{
		m:    map[string]string{},
		keys: []string{},
	}
	headers.SetHeaders(values...)
	return headers
}

func (h *Headers) SetHeaders(values ...headerValue) {
	for _, o := range values {
		o(h)
	}
}

func (h *Headers) Set(k string, v string) {
	_, present := h.m[k]
	h.m[k] = v
	if !present {
		h.keys = append(h.keys, k)
	}
}

func (h *Headers) Get(k string) string {
	return h.m[k]
}

func (h *Headers) Generate(shuffle bool) []string {
	var result []string
	if shuffle {
		for k, v := range h.m {
			result = append(result, "-H", fmt.Sprintf("'%s: %s'", k, v))
		}
	} else {
		for _, k := range h.keys {
			result = append(result, "-H", fmt.Sprintf("'%s: %s'", k, h.m[k]))
		}
	}
	return result
}
