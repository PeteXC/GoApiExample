package add

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/PeteXC/GoApiExample/api/respond"
	"github.com/a-h/pathvars"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

var MathsAddGetMatcher = pathvars.NewExtractor("*/maths/add")

type Handler struct {
	Log *zapray.Logger
}

type Input struct {
	NumberA int `json:"numberA"`
	NumberB int `json:"numberB"`
}

func (h Handler) Handle(r *http.Request, w http.ResponseWriter) (err error) {

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	var req Input
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		h.Log.Error("failed to decode request", zap.Error(err))
		respond.WithBadRequest(w, "failed to decode request")
		return
	}

	x := calculateAddition(int(req.NumberA), int(req.NumberB))
	fmt.Println(x)
	h.Log.With(zap.String("X", fmt.Sprint(x)))
	h.Log.Info(fmt.Sprint(x))
	respond.WithJSON(w, x)
	return
}

func calculateAddition(a, b int) (x int) {
	x = a + b
	return
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log = h.Log.TraceRequest(r)
	if _, ok := MathsAddGetMatcher.Extract(r.URL); ok {
		h.Handle(r, w)
	} else {
		respond.WithError(w, "not found", http.StatusNotFound)
	}
}

func NewHandler(log *zapray.Logger) (h Handler, err error) {
	return Handler{
		Log: log,
	}, err
}
