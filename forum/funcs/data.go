package funcs

import (
	"html/template"
)

var tpl *template.Template


type Error struct {
	Code   int
	Status string
}


// Used to indicate when user enters wrong data
type PopUp struct {
    MessageEmpty string
	MessageMail    string
	MessageName    string
	MessagePass    string
	MessageSuccess string
    MessageNoUser string
}

// Sets the client and template to be used inside the package
func SetT(t *template.Template) {
	tpl = t
}
