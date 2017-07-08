package httputil

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Novetta/kerbproxy/kerbtypes"
	"github.com/go-martini/martini"
)

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

func prepareStaticOptions(options []martini.StaticOptions) martini.StaticOptions {
	var opt martini.StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}
	// Normalize the prefix if provided
	if opt.Prefix != "" {
		// Ensure we have a leading '/'
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		// Remove any trailing '/'
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
	return opt
}

// Static returns a middleware handler that serves static files in the given directory.
func Static(directory string, staticOpt ...martini.StaticOptions) martini.Handler {
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(martini.Root, directory)
	}
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)

	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		if req.Method != "GET" && req.Method != "HEAD" {
			return
		}
		if opt.Exclude != "" && strings.HasPrefix(req.URL.Path, opt.Exclude) {
			return
		}
		file := req.URL.Path
		// if we have a prefix, filter requests by stripping the prefix
		if opt.Prefix != "" {
			if !strings.HasPrefix(file, opt.Prefix) {
				return
			}
			file = file[len(opt.Prefix):]
			if file != "" && file[0] != '/' {
				return
			}
		}
		f, err := dir.Open(file)
		if err != nil {
			// try any fallback before giving up
			if opt.Fallback != "" {
				file = opt.Fallback // so that logging stays true
				f, err = dir.Open(opt.Fallback)
			}

			if err != nil {
				// discard the error?
				return
			}
		}
		defer f.Close()
		dirf := f

		fi, err := f.Stat()
		if err != nil {
			return
		}

		// try to serve index file
		if fi.IsDir() {
			u := req.Header.Get(kerbtypes.XForwardedFor)
			if u == "" {
				u = req.URL.Path
			}
			// redirect if missing trailing slash
			if !strings.HasSuffix(u, "/") {
				dest := url.URL{
					Path:     u + "/",
					RawQuery: req.URL.RawQuery,
					Fragment: req.URL.Fragment,
				}
				http.Redirect(res, req, dest.String(), http.StatusFound)
				return
			}

			file = filepath.Join(file, opt.IndexFile)
			f, err = dir.Open(file)
			if err != nil {
				dirList(res, dirf)
				return
			}
			defer f.Close()
			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				return
			}
		}

		if !opt.SkipLogging {
			log.Println("[Static] Serving " + file)
		}

		// Add an Expires header to the static content
		if opt.Expires != nil {
			res.Header().Set("Expires", opt.Expires())
		}

		http.ServeContent(res, req, file, fi.ModTime(), f)
	}
}

func dirList(w http.ResponseWriter, f http.File) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name += "/"
			}
			// name may contain '?' or '#', which must be escaped to remain
			// part of the URL path, and not indicate the start of a query
			// string or fragment.
			url := url.URL{Path: name}
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
		}
	}
	fmt.Fprintf(w, "</pre>\n")
}
