package main

import (
	"net/http"

	"github.com/go-on/bootstrap/bs3"
	. "github.com/go-on/html"
	"github.com/go-on/html/h"
	. "github.com/go-on/html/tag"
)

func layout(inner interface{}) http.Handler {
	return HTML5(
		HTML(
			Attrs("lang", "en"),
			HEAD(
				META(Attrs("charset", "utf-8")),
				META(Attrs("http-equiv", "X-UA-Compatible", "content", "IE=edge")),
				META(Attrs("name", "viewport", "content", "width=device-width, initial-scale=1")),
				TITLE("Library Status Checker"),
				h.CssHref("//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css"),
				h.CssHref("//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap-theme.min.css"),
				Comment(" HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries "),
				Comment(" WARNING: Respond.js doesn't work if you view the page via file:// "),
				Comment("[if lt IE 9]>\n      <script src=\"https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js\"></script>\n      <script src=\"https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js\"></script>\n    <![endif]"),
			),
			BODY(
				DIV(bs3.Container_fluid, inner),
				DIV(
					Id("git-diff-modal"),
					bs3.Modal, bs3.Fade,
					Attrs("tabindex", "-1", "role", "dialog", "aria-labelledby", "myModalLabel", "aria-hidden", "true"),
					DIV(bs3.Modal_dialog,
						DIV(bs3.Modal_content,
							DIV(bs3.Modal_body, PRE()),
						),
					),
				),

				h.JsSrc("//ajax.googleapis.com/ajax/libs/jquery/1.11.0/jquery.min.js"),
				h.JsSrc("//netdna.bootstrapcdn.com/bootstrap/3.1.1/js/bootstrap.min.js"),
				SCRIPT(
					`
					$(function(){
						$('.git-diff').click(function (ev) {
							ev.preventDefault();
							$.get($(this).attr("href"), function (data) {
								console.log(data);
								if (data != "") {
									$('#git-diff-modal pre').text(data);
									$('#git-diff-modal').modal();	
								}								
							});
						});
					});
					`,
				),
			),
		),
	)
}
