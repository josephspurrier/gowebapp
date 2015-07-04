# GoWebApp
Basic Web Application in Golang

This project demonstrates how to structure and build a website using the Go language without a framework.

## Quick Start

To download, run the following command:

~~~
go get github.com/josephspurrier/gowebapp
~~~

Start MySQL and import config/database.sql to create the database and tables.

Open config/config.json and edit the Database section so the connection information matches your MySQL instance.

Build and run from the root directory. Open your web browser to: http://localhost. You should see the welcome page.

Navigate to the login page, and then to the register page. Create a new user and you should be able to login. That's it.

## Overview

The web app has a public home page, authenticated home page, login page, register page, and about page. 

The entrypoint for the web app is gowebapp.go. The file loads the application settings, 
starts the session, connects to the database, sets up the templates, loads 
the routes, attaches the middleware, and starts the web server.

The front end is built using Foundation with a few small changes to fonts and spacing. The flash 
messages are customized so they show up at the bottom right of the screen.

All of the error and warning messages should be either displayed either to the 
user or in the console. Informational messages are displayed to the user via 
flash messages that disappear after 5 seconds. The flash messages are controlled 
by javascript in the static folder.

## Structure

The project is organized into the following folders:

~~~
config		- application settings and database schema
controller	- page logic organized by HTTP methods (GET, POST)
model		- database queries
route		- route information and middleware
shared		- packages for templates, MySQL, cryptography, sessions, and json
static		- location of statically serve files like CSS and JS
template	- HTML templates
~~~

There are a few external packages:

~~~
github.com/gorilla/context				- registry for global request variables
github.com/gorilla/sessions				- cookie and filesystem sessions
github.com/go-sql-driver/mysql 			- MySQL driver
github.com/jmoiron/sqlx 				- MySQL general purpose extensions
github.com/josephspurrier/csrfbanana 	- CSRF protection for gorilla sessions
github.com/julienschmidt/httprouter 	- high performance HTTP request router
github.com/justinas/alice				- middleware chaining
golang.org/x/crypto/bcrypt 				- password hashing algorithm
~~~

The templates are:

~~~
about.tmpl		- quick info about the app
anon_home.tmpl	- public home page
auth_home.tmpl	- home page once you login
footer.tmpl		- footer
index.tmpl		- base template for all the pages
login.tmpl		- login page
menu.tmpl		- menu at the top of all the pages
register.tmpl	- register page
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

<!-- Same page hyperlinks -->
{{LINK "register" "Create a new account."}}
parses to
<a href="/register">Create a new account.</a>
~~~

There are a few variables you can use in templates as well:

~~~ html
<!-- Use AuthLevel=auth to determine if a user is logged in -->
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

Trigger a flash message on the next page load:

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
v.Name = "login"

// Assign a variable that is accessible in the form
v.Vars["token"] = csrfbanana.Token(w, r, sess)

// Refill any form fields from a POST operation
view.Repopulate([]string{"email"}, r.Form, v.Vars)

// Render the template
v.Render(w)
~~~

Handle the database query:

~~~ go
// Get database result
result, err := model.UserByEmail(email)

// Determine if password is correct
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

## Database

It's a good idea to abstract the database layer out so if you need to make 
changes, you don't have to look through business logic to find the queries. All
the queries are stored in the models folder:

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
_, err := database.DB.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", first_name, last_name, email, password)
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
		"Username": "root",
		"Password": "",
		"Name": "webframework",
		"Hostname": "127.0.0.1",
		"Port": 3306,
		"Parameter": "?parseTime=true"
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
		"Root": "index",
		"Children": [
			"menu",
			"footer"
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

![Image of Public Home](http://josephspurrier.com/images/gowebapp-home.JPG)

About:

![Image of About](http://josephspurrier.com/images/gowebapp-about.JPG)

Register:

![Image of Register](http://josephspurrier.com/images/gowebapp-register.JPG)

Login:

![Image of Login](http://josephspurrier.com/images/gowebapp-login.JPG)

Authenticated Home:

![Image of Auth Home](http://josephspurrier.com/images/gowebapp-loggedin.JPG)

## Feedback

All feedback is welcome. Let me know if you have any suggestions, questions, or criticisms. 
If something is not idiomatic to Golang, please let me know know so we can make it better.