package zerologhelper

import (
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
)

type helper struct {
}

func New() *helper {
	return &helper{}
}

func (h *helper) SetInterfaceMarshalFunc() *helper {
	zerolog.InterfaceMarshalFunc = json.Marshal
	return h
}

func (h *helper) SetErrorStackMarshaler() *helper {
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		var serr error
		for ; err != nil; err = errors.Unwrap(err) {
			_, ok := err.(errbase.StackTraceProvider)
			if ok {
				serr = err
			}
		}
		safeDetails := errors.GetSafeDetails(serr).SafeDetails
		if len(safeDetails) == 1 {
			return ParsePII(safeDetails[0])
		}
		res := make([][]StackInfo, 0)
		for _, details := range errors.GetSafeDetails(serr).SafeDetails {
			piis := ParsePII(details)
			res = append(res, piis)
		}
		return res
	}
	return h
}

func (h *helper) ConsoleWriter() zerolog.ConsoleWriter {
	return zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime + " -0700"
	})
}

type StackInfo struct {
	Package string `json:"package"`
	Func    string `json:"func"`
	File    string `json:"file"`
	Line    string `json:"line"`
}

func ParsePII(detail string) []StackInfo {
	s := strings.TrimSpace(detail)
	ss := strings.Split(s, "\n")

	var piis []StackInfo
	for i := 0; i < len(ss); i += 2 {
		pkgAndFunc := strings.Split(ss[i], ".")
		pathAndLine := strings.Split(ss[i+1], ":")
		piis = append(piis, StackInfo{
			Package: strings.TrimSpace(pkgAndFunc[0]),
			Func:    strings.TrimSpace(pkgAndFunc[1]),
			File:    strings.TrimSpace(pathAndLine[0]),
			Line:    strings.TrimSpace(pathAndLine[1]),
		})
	}

	return piis
}
