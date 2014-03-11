package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	. "github.com/go-on/html/attr"

	"github.com/go-on/bootstrap/bs3"
	. "github.com/go-on/html"
	. "github.com/go-on/html/entity"
	"github.com/go-on/html/h"
	. "github.com/go-on/html/tag"
	"github.com/go-on/queue/q"
)

func libPathsHandler(rw http.ResponseWriter, req *http.Request) {
	list := TBODY()

	libpaths := config.LibPaths

	sort.Strings(libpaths)

	for _, lp := range libpaths {
		list.Add(
			TR(
				TD(CODE(lp)),
				TD(
					// A_Width("200px"),
					h.AHref(routeStatus.MustURL()+"?library_path="+url.QueryEscape(lp), "status", bs3.Btn, bs3.Btn_primary),
				),
				TD(
					// A_Width("100px"),
					h.FormPost(routeDeleteLibraryPath.MustURL(),
						h.InputHidden("library_path", A_Value(lp)),
						h.InputSubmit("submit", A_Value("delete"), bs3.Btn, bs3.Btn_danger),
					),
				),
			),
		)
	}

	TABLE(bs3.Table, bs3.Table_striped, list).WriteTo(rw)
}

func deleteLibraryPathHandler(rw http.ResponseWriter, req *http.Request) {
	q.Err(q.PANIC)(
		req.PostFormValue, "library_path",
	)(
		config.deleteLibraryPath, q.V,
	).Run()

	http.Redirect(rw, req, routeStatus.MustURL(), 302)
}

func addLibraryPathHandler(rw http.ResponseWriter, req *http.Request) {
	q.Err(q.PANIC)(
		req.PostFormValue, "library_path",
	)(
		config.addLibraryPath, q.V,
	).Run()

	http.Redirect(rw, req, routeHome.MustURL(), 302)
}

func gitDiffHandler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Println("gitDiffHandler called")
	library_path := req.URL.Query().Get("library_path")
	library := req.URL.Query().Get("library")

	fullpath := filepath.Join(gopath, "src", library_path)
	d, err := gitDiff(filepath.Join(fullpath, library))

	if err != nil {
		fmt.Println(err.Error(), string(d))
		rw.WriteHeader(500)
	} else {
		rw.Write(d)
	}
}

func gitPushHandler(rw http.ResponseWriter, req *http.Request) {
	library_path := req.PostFormValue("library_path")
	library := req.PostFormValue("library")
	fullpath := filepath.Join(gopath, "src", library_path)
	p, err := gitPush(filepath.Join(fullpath, library))

	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("%s\n", p)
		fmt.Fprintln(rw, err.Error())
		rw.Write(p)
		rw.WriteHeader(500)
	} else {
		http.Redirect(rw, req, routeStatus.MustURL()+"?library_path="+url.QueryEscape(library_path), 302)
	}
}

func depTrackLibraryPath(rw http.ResponseWriter, req *http.Request) {
	library_path := req.PostFormValue("library_path")
	fullpath := filepath.Join(gopath, "src", library_path)

	dir, err := os.Open(fullpath)
	defer dir.Close()

	if err != nil {
		panic(err.Error())
	}

	names, e := dir.Readdirnames(-1)
	if e != nil {
		panic(e.Error())
	}

	for _, n := range names {
		depTrackForRepo(filepath.Join(fullpath, n))
	}

	http.Redirect(rw, req, routeStatus.MustURL()+"?library_path="+library_path, 302)
}

func depGDFLibraryPath(rw http.ResponseWriter, req *http.Request) {
	library_path := req.PostFormValue("library_path")
	fullpath := filepath.Join(gopath, "src", library_path)

	dir, err := os.Open(fullpath)
	defer dir.Close()
	if err != nil {
		panic(err.Error())
	}

	names, e := dir.Readdirnames(-1)
	if e != nil {
		panic(e.Error())
	}

	for _, n := range names {
		depGDFForRepo(filepath.Join(fullpath, n))
	}

	http.Redirect(rw, req, routeStatus.MustURL()+"?library_path="+library_path, 302)
}

func statusLibraryPathHandler(rw http.ResponseWriter, req *http.Request) {
	library_path := req.URL.Query().Get("library_path")
	fullpath := filepath.Join(gopath, "src", library_path)

	dir, err := os.Open(fullpath)
	defer dir.Close()
	if err != nil {
		panic(err.Error())
	}

	names, e := dir.Readdirnames(-1)

	if e != nil {
		panic(e.Error())
	}

	sort.Strings(names)
	tbodyOk := TBODY()
	tbodyModified := TBODY()

	for i, n := range names {
		st, e := gitStatusForRepo(filepath.Join(fullpath, n))
		if e != nil {
			continue
		}

		// b := BUTTON(bs3.Btn, n)
		// state := SPAN()
		if st == "" {
			tbodyOk.Add(
				TR(
					TD(BUTTON(bs3.Btn, bs3.Btn_success, n)),
				),
			)
			// b.Add(bs3.Btn_success)
		} else {
			id := fmt.Sprintf("status-%d", i)
			tbodyModified.Add(
				TR(
					// data-toggle="collapse"
					TD(
						h.AHref("#"+id,
							bs3.Btn, bs3.Btn_danger,
							Attr("data-toggle", "collapse"),
							n,
						),
						BR(),
						BR(),
						h.AHref(routeGitDiff.MustURL()+"?library_path="+url.QueryEscape(library_path)+"&library="+url.QueryEscape(n),
							bs3.Btn, bs3.Btn_default, bs3.Btn_sm,
							Class("git-diff"),
							"git diff",
						),
						E_nbsp,
						h.FormPost(routeGitPush.MustURL(),
							Style("display", "inline"),
							h.InputHidden("library_path", A_Value(library_path)),
							h.InputHidden("library", A_Value(n)),
							h.InputSubmit("submit", A_Value("git push"), bs3.Btn, bs3.Btn_warning, bs3.Btn_sm),
						),
						BR(),
						BR(),
					),
					TD(
						PRE(
							Id(id),
							bs3.Collapse,
							CODE(st),
						),
					),
				),
			)
		}
	}

	colOk := DIV(
		bs3.Col_lg_4, bs3.Col_md_4, bs3.Col_sm_4,
		TABLE(bs3.Table, bs3.Table_striped, tbodyOk),
	)

	colModified := DIV(
		bs3.Col_lg_8, bs3.Col_md_8, bs3.Col_sm_8,
		TABLE(bs3.Table, bs3.Table_striped, tbodyModified),
	)

	Elements(
		DIV(bs3.Row,
			DIV(
				bs3.Col_lg_12, bs3.Col_md_12, bs3.Col_sm_12,
				H1(
					"Status for ",
					CODE(library_path),
					E_nbsp,
					h.AHref(routeHome.MustURL(), "back", bs3.Btn, bs3.Btn_primary, bs3.Btn_sm),
				),
				h.FormPost(routeDepTrack.MustURL(),
					Style("display", "inline"),
					h.InputHidden("library_path", A_Value(library_path)),
					h.InputSubmit("submit", A_Value("dep track for all"), bs3.Btn, bs3.Btn_default, bs3.Btn_sm),
				),
				E_nbsp,
				h.FormPost(routeDepGDF.MustURL(),
					Style("display", "inline"),
					h.InputHidden("library_path", A_Value(library_path)),
					h.InputSubmit("submit", A_Value("dep gdf for all"), bs3.Btn, bs3.Btn_default, bs3.Btn_sm),
				),
				BR(), BR(),
			),
		),
		DIV(bs3.Row,
			DIV(
				bs3.Col_lg_12, bs3.Col_md_12, bs3.Col_sm_12,
				PRE(CODE(fullpath)),
				HR(),
			),
		),
		DIV(bs3.Row,
			colModified, colOk,
		),
	).WriteTo(rw)

}
