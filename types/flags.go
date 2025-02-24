package types

import "fmt"

type Flags struct {
	m map[string]interface{}
}

type flagValue func(f *Flags)

func Flag(name string, v interface{}) func(flags *Flags) {
	return func(flags *Flags) {
		flags.m[name] = v
	}
}

func NewFlags(values ...flagValue) *Flags {
	f := &Flags{
		m: make(map[string]interface{}),
	}
	f.SetFlags(values...)
	return f
}

func (f *Flags) SetFlags(values ...flagValue) {
	for _, o := range values {
		o(f)
	}
}

func (f *Flags) Set(name string, v interface{}) {
	f.m[name] = v
}

func (f *Flags) Generate() []string {
	var result []string
	for k, v := range f.m {
		switch v.(type) {
		case nil:
			result = append(result, fmt.Sprintf("--%s", k))
		case bool:
			if v == true {
				result = append(result, fmt.Sprintf("--%s", k))
			}
		case int:
			result = append(result, fmt.Sprintf("--%s", k), fmt.Sprintf("%d", v))
		case float64:
			result = append(result, fmt.Sprintf("--%s", k), fmt.Sprintf("%.4f", v))
		case string:
			result = append(result, fmt.Sprintf("--%s", k), v.(string))
		default:
			panic("Unsupported type of flag")
		}
	}

	return result
}
