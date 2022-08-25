package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
//        "fmt"

	"github.com/gatopardo/acad/ch14/fcondo/app/route"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/email"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/jsonconfig"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/recaptcha"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/server"
	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view/plugin"
)

var  file  * os.File

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)
        var  err error
	file, err = os.OpenFile("fcondo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
        log.SetOutput(file)
	route.Flogger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
        route.Flogger.Println("Starting app")
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	model.Configure(config.Session)
	// Connect to database
	model.Connect(config.Database)
        defer model.Db.Close()
        defer file.Close()

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

        // Configure email
        email.Configure(config.Email)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		plugin.DateFormat(),
		plugin.Format64(),
		plugin.ConcatStr(),
		recaptcha.Plugin())
	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  model.Info      `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   model.Session   `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

