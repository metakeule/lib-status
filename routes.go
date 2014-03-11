package main

import (
	"net/http"

	"github.com/go-on/bootstrap/bs3"
	. "github.com/go-on/html"
	. "github.com/go-on/html/attr"
	. "github.com/go-on/html/entity"
	"github.com/go-on/html/h"
	. "github.com/go-on/html/tag"
	"github.com/go-on/router"
)

var (
	routeAddLibraryPath    *router.Route
	routeDeleteLibraryPath *router.Route
	routeHome              *router.Route
	routeStatus            *router.Route
	routeGitPush           *router.Route
	routeGitDiff           *router.Route
	routeDepGDF            *router.Route
	routeDepTrack          *router.Route
)

func setupRouter() {
	rt := router.New()
	routeDepGDF = rt.POST("/dep-gdf-all", http.HandlerFunc(depGDFLibraryPath))
	routeDepTrack = rt.POST("/dep-gdf-track", http.HandlerFunc(depTrackLibraryPath))
	routeGitDiff = rt.GET("/git-diff", http.HandlerFunc(gitDiffHandler))
	routeGitPush = rt.POST("/git-push", http.HandlerFunc(gitPushHandler))
	routeAddLibraryPath = rt.POST("/add_library_path", http.HandlerFunc(addLibraryPathHandler))
	routeDeleteLibraryPath = rt.POST("/delete_library_path", http.HandlerFunc(deleteLibraryPathHandler))
	routeStatus = rt.GET("/status", layout(statusLibraryPathHandler))
	routeHome = rt.GET("/",
		layout(
			Elements(
				DIV(bs3.Row,
					DIV(
						bs3.Col_lg_12, bs3.Col_md_12, bs3.Col_sm_12,
						H1("Library Paths"),
					),
				),
				DIV(bs3.Row,
					DIV(
						bs3.Col_lg_12, bs3.Col_md_12, bs3.Col_sm_12,
						h.FormPost(routeAddLibraryPath.MustURL(), A_Role("form"), bs3.Form_inline,
							DIV(bs3.Form_group,
								h.InputText("library_path", bs3.Form_control),
								E_nbsp, E_nbsp,
								h.InputSubmit("submit", A_Value("add"), bs3.Btn, bs3.Btn_success),
							),
						),
						HR(),
					),
				),
				DIV(bs3.Row,
					DIV(
						bs3.Col_lg_12, bs3.Col_md_12, bs3.Col_sm_12,
						libPathsHandler,
					),
				),
			),
		),
	)
	router.Mount("/", rt)

}
