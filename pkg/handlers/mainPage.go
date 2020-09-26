package handlers

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zherdev/go-panic-beautifier/pkg/cfg"
	"github.com/zherdev/go-panic-beautifier/pkg/parser"
	"html/template"
	"net/http"
)

const examplePanicLog = `panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x109ab63]

goroutine 1 [running]:
main.f(...)
        /Users/i.zherdev/Documents/code/myGo/src/github.com/zherdev/go-error-beautifier/cmd/goErrorBeautifier/main.go:6
main.main()
        /Users/i.zherdev/Documents/code/myGo/src/github.com/zherdev/go-error-beautifier/cmd/goErrorBeautifier/main.go:14 +0x23
exit status 2
`

type mainPageTemplateData struct {
	PanicLog template.HTML
}

type MainPageHandler struct {
	baseHandler
	template *template.Template
}

func NewMainPageHandler(lg *logrus.Logger, config *cfg.Config) (*MainPageHandler, error) {
	tmpl, err := template.New("main_page.html").ParseFiles(config.MainPageTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse template %q in NewMainPageHandler",
			config.MainPageTemplate)
	}

	return &MainPageHandler{
		baseHandler: newBaseHandler(lg, "MainPageHandler"),
		template:    tmpl,
	}, nil
}

func (h *MainPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lg := h.newLoggerFromRequest(r)

	defer h.recover(lg)

	var panicLog string

	err := r.ParseForm()
	if err != nil {
		lg.WithError(err).Warn("can't parse post in MainPageHandler")
	} else {
		panicLog = r.PostForm.Get("panic_log")
	}

	if panicLog == "" {
		panicLog = examplePanicLog
	}

	result := h.processLog(panicLog)

	data := mainPageTemplateData{PanicLog: template.HTML(result)}
	b := &bytes.Buffer{}
	err = h.template.Execute(b, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.WithError(err).Error("error at template execution")
	} else {
		n, err := w.Write(b.Bytes())
		if err != nil {
			lg.WithError(err).WithField("written", n).Error("error at template send")
		}
	}

	lg.Trace("done")
}

func (h *MainPageHandler) processLog(panicLog string) string {
	p := parser.NewPanicLogParser(panicLog)
	return p.Parse()
}
