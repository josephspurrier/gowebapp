# GoWebApp

[![Go Report Card](https://goreportcard.com/badge/github.com/josephspurrier/gowebapp)](https://goreportcard.com/report/github.com/josephspurrier/gowebapp)
[![GoDoc](https://godoc.org/github.com/josephspurrier/gowebapp?status.svg)](https://godoc.org/github.com/josephspurrier/gowebapp) 

Basic MVC Web Application in Go

#### I recommend you use Blue Jay which is the latest version of this project: [https://github.com/blue-jay/blueprint](https://github.com/blue-jay/blueprint).

This project demonstrates how to structure and build a website using the Go language without a framework. There is a blog article you can read at [http://www.josephspurrier.com/go-web-app-example/](http://www.josephspurrier.com/go-web-app-example/). There is a full application I built with an earlier version of the project at [https://github.com/verifiedninja/webapp](https://github.com/verifiedninja/webapp). There is an API version of this project at [https://github.com/josephspurrier/gowebapi](https://github.com/josephspurrier/gowebapi).

To download, run the following command:

~~~
go get github.com/josephspurrier/gowebapp
~~~

If you are on Go 1.5, you need to set GOVENDOREXPERIMENT to 1. If you are on Go 1.4 or earlier, the code will not work because it uses the vendor folder.

### Other Branches

- 2021-02-17: There is a branch using Go modules (requires Go 1.13) available here: https://github.com/josephspurrier/gowebapp/tree/modules.
- 2021-05-23: There is a branch using Go embedded assets (requires Go 1.16) available here: https://github.com/josephspurrier/gowebapp/tree/embedded-templates.

## Quick Start with Bolt

The gowebapp.db file will be created once you start the application.

Build and run from the root directory. Open your web browser to: http://localhost. You should see the welcome page.

Navigate to the login page, and then to the register page. Create a new user and you should be able to login. That's it.

## Quick Start with MongoDB

Start MongoDB.

Open config/config.json and edit the Database section so the connection information matches your MongoDB instance. Also, change Type from Bolt to MongoDB.

Build and run from the root directory. Open your web browser to: http://localhost. You should see the welcome page.

Navigate to the login page, and then to the register page. Create a new user and you should be able to login. That's it.

## Quick Start with MySQL

Start MySQL and import config/mysql.sql to create the database and tables.

Open config/config.json and edit the Database section so the connection information matches your MySQL instance. Also, change Type from Bolt to MySQL.

Build and run from the root directory. Open your web browser to: http://localhost. You should see the welcome page.

Navigate to the login page, and then to the register page. Create a new user and you should be able to login. That's it.

## Overview

The web app has a public home page, authenticated home page, login page, register page,
about page, and a simple notepad to demonstrate the CRUD operations.

The entrypoint for the web app is gowebapp.go. The file loads the application settings, 
starts the session, connects to the database, sets up the templates, loads 
the routes, attaches the middleware, and starts the web server.

The front end is built using Bootstrap with a few small changes to fonts and spacing. The flash 
messages are customized so they show up at the bottom right of the screen.

All of the error and warning messages should be either displayed either to the 
user or in the console. Informational messages are displayed to the user via 
flash messages that disappear after 4 seconds. The flash messages are controlled 
by JavaScript in the static folder.

## Structure

Recently, the folder structure changed. After looking at all the forks 
and reusing my project in different places, I decided to move the Go code to the 
**app** folder inside the **vendor** folder so the github path is not littered 
throughout the many imports. I did not want to use relative paths so the vendor
folder seemed like the best option.

The project is organized into the following folders:

~~~
config		- application settings and database schema
static		- location of statically served files like CSS and JS
template	- HTML templates

vendor/app/controller	- page logic organized by HTTP methods (GET, POST)
vendor/app/shared		- packages for templates, MySQL, cryptography, sessions, and json
vendor/app/model		- database queries
vendor/app/route		- route information and middleware
~~~

There are a few external packages:

~~~
github.com/gorilla/context				- registry for global request variables
github.com/gorilla/sessions				- cookie and filesystem sessions
github.com/go-sql-driver/mysql 			- MySQL driver
github.com/haisum/recaptcha				- Google reCAPTCHA support
github.com/jmoiron/sqlx 				- MySQL general purpose extensions
github.com/josephspurrier/csrfbanana 	- CSRF protection for gorilla sessions
github.com/julienschmidt/httprouter 	- high performance HTTP request router
github.com/justinas/alice				- middleware chaining
github.com/mattn/go-sqlite3				- SQLite driver
golang.org/x/crypto/bcrypt 				- password hashing algorithm
~~~

The templates are organized into folders under the **template** folder:

~~~
about/about.tmpl       - quick info about the app
index/anon.tmpl	       - public home page
index/auth.tmpl	       - home page once you login
login/login.tmpl	   - login page
notepad/create.tmpl    - create note
notepad/read.tmpl      - read a note
notepad/update.tmpl    - update a note
partial/footer.tmpl	   - footer
partial/menu.tmpl	   - menu at the top of all the pages
register/register.tmpl - register page
base.tmpl		       - base template for all the pages
~~~

## Templates

There are a few template funcs that are available to make working with the templates 
and static files easier:

~~~ html
<!-- CSS files with timestamps -->
{{CSS "static/css/normalize3.0.0.min.css"}}
parses to
<link rel="stylesheet" type="text/css" href="/static/css/normalize3.0.0.min.css?1435528339" />

<!-- JS files with timestamps -->
{{JS "static/js/jquery1.11.0.min.js"}}
parses to
<script type="text/javascript" src="/static/js/jquery1.11.0.min.js?1435528404"></script>

<!-- Hyperlinks -->
{{LINK "register" "Create a new account."}}
parses to
<a href="/register">Create a new account.</a>

<!-- Output an unescaped variable (not a safe idea, but I find it useful when troubleshooting) -->
{{.SomeVariable | NOESCAPE}}

<!-- Time format -->
{{.SomeTime | PRETTYTIME}}
parses to format
3:04 PM 01/02/2006
~~~

There are a few variables you can use in templates as well:

~~~ html
<!-- Use AuthLevel=auth to determine if a user is logged in (if session.Values["id"] != nil) -->
{{if eq .AuthLevel "auth"}}
You are logged in.
{{else}}
You are not logged in.
{{end}}

<!-- Use BaseURI to print the base URL of the web app -->
<li><a href="{{.BaseURI}}about">About</a></li>

<!-- Use token to output the CSRF token in a form -->
<input type="hidden" name="token" value="{{.token}}">
~~~

It's also easy to add template-specific code before the closing </head> and </body> tags:

~~~ html
<!-- Code is added before the closing </head> tag -->
{{define "head"}}<meta name="robots" content="noindex">{{end}}

...

<!-- Code is added before the closing </body> tag -->
{{define "foot"}}{{JS "//www.google.com/recaptcha/api.js"}}{{end}}
~~~

## JavaScript

You can trigger a flash notification using JavaScript.

~~~ javascript
flashError("You must type in a username.");

flashSuccess("Record created!");

flashNotice("There seems to be a piece missing.");

flashWarning("Something does not seem right...");
~~~

## Controllers

The controller files all share the same package name. This cuts down on the 
number of packages when you are mapping the routes. It also forces you to use
a good naming convention for each of the funcs so you know where each of the 
funcs are located and what type of HTTP request they each are mapped to.

### These are a few things you can do with controllers.

Access a gorilla session:

~~~ go
// Get the current session
sess := session.Instance(r)
...
// Close the session after you are finished making changes
sess.Save(r, w)
~~~

Trigger 1 of 4 different types of flash messages on the next page load (no other code needed):

~~~ go
sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
sess.Save(r, w) // Ensure you save the session after making a change to it
~~~

Validate form fields are not empty:

~~~ go
// Ensure a user submitted all the required form fields
if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
	sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
	sess.Save(r, w)
	LoginGET(w, r)
	return
}
~~~

Render a template:

~~~ go
// Create a new view
v := view.New(r)

// Set the template name
v.Name = "login/login"

// Assign a variable that is accessible in the form
v.Vars["token"] = csrfbanana.Token(w, r, sess)

// Refill any form fields from a POST operation
view.Repopulate([]string{"email"}, r.Form, v.Vars)

// Render the template
v.Render(w)
~~~

Return the flash messages during an Ajax request:

~~~ go
// Get session
sess := session.Instance(r)

// Set the flash message
sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
sess.Save(r, w)

// Display the flash messages as JSON
v := view.New(r)
v.SendFlashes(w)
~~~

Handle the database query:

~~~ go
// Get database result
result, err := model.UserByEmail(email)

if err == sql.ErrNoRows {
	// User does not exist
} else if err != nil {
	// Display error message
} else if passhash.MatchString(result.Password, password) {
	// Password matches!	
} else {
	// Password does not match
}
~~~

Send an email:

~~~ go
// Email a user
err := email.SendEmail(email.ReadConfig().From, "This is the subject", "This is the body!")
if err != nil {
	log.Println(err)
	sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
	sess.Save(r, w)
	return
}
~~~

Validate a form if the Google reCAPTCHA is enabled in the config:

~~~ go
// Validate with Google reCAPTCHA
if !recaptcha.Verified(r) {
	sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
	sess.Save(r, w)
	RegisterGET(w, r)
	return
}
~~~

## Database

It's a good idea to abstract the database layer out so if you need to make 
changes, you don't have to look through business logic to find the queries. All
the queries are stored in the models folder.

This project supports BoltDB, MongoDB, and MySQL. All the queries are stored in
the same files so you can easily change the database without modifying anything
but the config file.

The user.go and note.go files are at the root of the model directory and are a
compliation of all the queries for each database type. There are a few hacks in
the models to get the structs to work with all the supported databases.

Connect to the database (only once needed in your application):

~~~ go
// Connect to database
database.Connect(config.Database)
~~~

Read from the database:

~~~ go
result := User{}
err := database.DB.Get(&result, "SELECT id, password, status_id, first_name FROM user WHERE email = ? LIMIT 1", email)
return result, err
~~~

Write to the database:

~~~ go
_, err := database.DB.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", firstName, lastName, email, password)
return err
~~~

## Middleware

There are a few pieces of middleware included. The package called csrfbanana 
protects against Cross-Site Request Forgery attacks and prevents double submits. 
The package httprouterwrapper provides helper functions to make funcs compatible 
with httprouter. The package logrequest will log every request made against the 
website to the console. The package pprofhandler enables pprof so it will work 
with httprouter. In route.go, all the individual routes use alice to make 
chaining very easy.

## Configuration

To make the web app a little more flexible, you can make changes to different 
components in one place through the config.json file. If you want to add any 
of your own settings, you can add them to config.json and update the structs
in gowebapp.go and the individual files so you can reference them in your code. 
This is config.json:

~~~ json
{
	"Database": {
		"Type": "Bolt",
		"Bolt": {		
 			"Path": "gowebapp.db"
  		},
		"MongoDB": {
			"URL": "127.0.0.1",
			"Database": "gowebapp"
		},
		"MySQL": {
			"Username": "root",
			"Password": "",
			"Name": "gowebapp",
			"Hostname": "127.0.0.1",
			"Port": 3306,
			"Parameter": "?parseTime=true"
		}
	},
	"Email": {
		"Username": "",
		"Password": "",
		"Hostname": "",
		"Port": 25,
		"From": ""
	},
	"Recaptcha": {
		"Enabled": false,
		"Secret": "",
		"SiteKey": ""
	},
	"Server": {
		"Hostname": "",
		"UseHTTP": true,
		"UseHTTPS": false,
		"HTTPPort": 80,
		"HTTPSPort": 443,
		"CertFile": "tls/server.crt",
		"KeyFile": "tls/server.key"
	},
	"Session": {
		"SecretKey": "@r4B?EThaSEh_drudR7P_hub=s#s2Pah",
		"Name": "gosess",
		"Options": {
			"Path": "/",
			"Domain": "",
			"MaxAge": 28800,
			"Secure": false,
			"HttpOnly": true
		}
	},
	"Template": {
		"Root": "base",
		"Children": [
			"partial/menu",
			"partial/footer"
		]
	},
	"View": {
		"BaseURI": "/",
		"Extension": "tmpl",
		"Folder": "template",
		"Name": "blank",
		"Caching": true
	}
}
~~~

To enable HTTPS, set UseHTTPS to true, create a folder called tls in the root, 
and then place the certificate and key files in that folder.

## Screenshots

Public Home:

![Image of Public Home](https://cloud.githubusercontent.com/assets/2394539/11319464/e2cd0eac-9045-11e5-9b24-5e480240cd69.jpg)

About:

![Image of About](https://cloud.githubusercontent.com/assets/2394539/11319462/e2c4d2d2-9045-11e5-805f-8b40598c92c3.jpg)

Register:

![Image of Register](https://cloud.githubusercontent.com/assets/2394539/11319466/e2d03500-9045-11e5-9c8e-c28fe663ed0f.jpg)

Login:

![Image of Login](https://cloud.githubusercontent.com/assets/2394539/11319463/e2cd1a00-9045-11e5-8b8e-68030d870cbe.jpg)

Authenticated Home:

![Image of Auth Home](https://cloud.githubusercontent.com/assets/2394539/14809208/75f340d2-0b59-11e6-8d2a-cd26ee872281.PNG)

View Notes:

![Image of Notepad View](https://cloud.githubusercontent.com/assets/2394539/14809205/75f08432-0b59-11e6-8737-84ee796bd82e.PNG)

Add Note:

![Image of Notepad Add](https://cloud.githubusercontent.com/assets/2394539/14809207/75f338f8-0b59-11e6-9719-61355957996c.PNG)

Edit Note:

![Image of Notepad Edit](https://cloud.githubusercontent.com/assets/2394539/14809206/75f33970-0b59-11e6-8acf-b3d533477aac.PNG)

## Feedback

All feedback is welcome. Let me know if you have any suggestions, questions, or criticisms. 
If something is not idiomatic to Go, please let me know know so we can make it better.
