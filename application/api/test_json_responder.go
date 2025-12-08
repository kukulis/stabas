package api

type TestJSONResponder struct {
	StatusCode int
	Response   any
}

func (t *TestJSONResponder) JSON(code int, obj any) {
	t.StatusCode = code
	t.Response = obj
}
