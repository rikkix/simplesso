package url

import "github.com/valyala/fasthttp"

func AddQueries(base string, queries map[string]string) string {
	if len(queries) == 0 {
		return base
	}
	args := fasthttp.Args{}
	for k, v := range queries {
		if v == "" {
			continue
		}
		args.Add(k, v)
	}
	return base + "?" + args.String()
}