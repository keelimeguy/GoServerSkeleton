A note to the original author: I know that I've lost some functionality by porting your code as aggressively as I did. Ultimately, it's a bit easier to build independent RESTful endpoints and more. I do appreciate your collaboration on this.

# Go Server Skeleton

A skeleton for making website servers in Go. It implements the `just*` method of portable, configurable API attachments to handlers. This uses www.github.com/autopogo/justClaims as an API

## Just

The basic idea of all `just*` packages is to

1. Import and configure the API config structures in `main.go` or wherever you set http handlers.
2. Import the (or an interface that supports a class of them) API when defining your handler.
3. Negotiate with the configurations that `main.go` has supplied when initializing your handler by:
  * Using a closure to return a `HandlerFunc` for `http.HandleFunc()`
  * Define and instantiate a type that fufills the `Handler` interface

## More

Look through the files. Version control in the skeleton sucks, but that's a company policy more than a programming issue so I'll leave it.

# Todo:

### Stuff AJ Put Away (this is honestly for the skeleton and justClaims etc, the whole thing)
  * Update w/ claims, API has broke on next pull
	* maybe document for godocs
  * is error processing really sufficient? - we need an error implementation right?
  * prevent collisions (key names, JWT, cookies, etc)
	* better abstraction of templates
	* I didn't port /home
