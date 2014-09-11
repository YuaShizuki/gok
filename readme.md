gok
===
Build web apps in go, php style.
```text
<html>
	<body>
		<?go for i := 0; i < 5; i++ { ?>
			<p>Iteration <?go gok.Echo(i) ?></p>
		<?go } ?>
	</body>
</html>
```
everything between `<?go` `?>` can be valid go code

installation:
-------------
```bash
	go get https://github.com/YuaShizuki/gok.git
	echo "PATH=$PATH:$GOPATH/bin" >> ~/.bash_profile
```
the last command makes sure that $GOPATH/bin is in your standard path

gok tags:
---------
1.	**<?go ?>**  
	can scope any valid go code. Code scoped in this tag 
	executes inside a render function. 

2.	**<?gofn funcName( args ) ( return types ) { go code }?>** 
	allows you to define a global function thats accessible by all ".gok" scripts.
	functions defined in this tag cannot be methods of objects or types.
	hence:    	
	`<?gofn (type) funcName( args ) (return types) { } ?>`   
	would be incorect.

3. **<?goimp ?>**   
	tag allows for imports in gok scripts. example   
	```text
		<?goimp 
			"fmt"
			"io/ioutil"
		?>
		<html>
			<body>
				...
	```
4. **<?gouse ?>**    
	gouse tag allows to define global variables, functions and types. 

	```text
		<?gouse
			type nstring string 
			var x nstring = "Hello World"
			func (self nstring) PrintMe() { 
				fmt.Println(self)
			}
			...
		?>    
	```
5. **<?go@fn funcName(args []string) ([]string, error) { go code } ?>**      
	go ajax accessed function.      
	this is a special function, that is only accessed through javascript, hence the limitations
	on the parameters. Calling this function from javascript involves including gok.js in html.
	example:

	```html
		<?go@fn Swap(args []string) ([]string, error) {
			return []string{args[1], args[0]}, nil
		}?>
		<script src="/gok.js"></script>
		<script>
			function callback(response) {
				console.log(response) //prints => ["Hello World", "1"]
			}
			gok.Swap(1, "Hello World", callback);
		</script>
	```    
	this allows for quick ajax.

commands:
---------
1. **`$gok build`**   
	converts .gok file to valid go files and them compiles them to the final executable.

2. **`$gok run`**    
	gok run builds the executable and runs it, any new changes made to the source
	would update the the server executable and restart it.

3. **`$gok src`**    
	src converts .gok files to genrated .go files, running `$go build` in this directory
	would result in the final server executable.

4. **`$gok api`**    
	pritnts core api for instance `gok`

API:
----
*	__gok:__ Core struct encapsulates the http.Request and http.ResponseWriter.     
	Instance present in all `<?go ?>` tags as `gok`.    
	```go
		type Gok struct { ... }
	```    

*	__Echo:__ Equivalent to php `echo`. echo accepts any type of parameters.
	```go
		func (self *Gok) Echo(a ...interface{}) { ... }
	```
*	__Redirect:__ redirects the request to a new url    

	```go
		func (self *Gok) Redirect(newUrl string) { ... }
	```
*	__Die:__ sends the `msg` as error, and undoes everything echoed so far.

	```go
		func (self *Gok) Die(msg string) { ... }
	```
*	__Server Functions__ PHP equivalent to $_SERVER[' ']
    
    ```go
        //returns current path, without the leadind '/'
		func (self *Gok) ServerSelf() string { ... }
		func (self *Gok) ServerSelf() string { ... }
		func (self *Gok) ServerHttpUserAgent() string { ... }
		func (self *Gok) ServerHttpReferer() string { ... }
		func (self *Gok) ServerHttps() bool { ... }
		func (self *Gok) ServerRemoteAddr() string { ... }
        func (self *Gok) ServerRemotePort() string { ... }
        func (self *Gok) ServerPort() int { ... }
        func (self *Gok) ServerHttpAcceptEncoding() string { ... }
        func (self *Gok) ServerProtocol() string { ... }
        //returns "POST" or "GET"
        func (self *Gok) ServerRequestMethod() string { ... }
        func (self *Gok) ServerQueryString() string { ... }
        func (self *Gok) ServerHttpAccept() string { ... }
        func (self *Gok) ServerHttpAcceptCharset() string { ... }
        func (self *Gok) ServerHttpAcceptLanguage() string { ... }
        func (self *Gok) ServerHttpConnection() string { ... }
        func (self *Gok) ServerHttpHost() string { ... }
    ```
*	__Get and Post:__ PHP equivalent to $_GET[' '] and $_POST[' ']

	```go
		// parses the post request and returns the post value of `name`,
		// in case for post containing file date. use gok.File('fileName')
		func (self *Gok) Post(name string) string { ... }

		// the same as above but for get requests.
		func (self *Gok) Get(name string) string { ... }
	```

*	__Cookies:__ PHP equivalent of $_COOKIE[' ']
	
	```go
		// returns the cookie value for `name`, if no cookie is set returns an
		// empty string ""
		func (self *Gok) Cookie(name string) string { ... }

		// sets cookie with name, value, and duration to expire, if duration is 0
		// cookie is permanent
		func (self *Gok) SetCookie(name string, value string, duration int64) { ... }

		//delets the cookie set with `name`
		func (self *Gok) DeleteCookie(name string) { ... }

		//set cookie with extra params
		func (self *Gok) SetCookie_4(name string, value string, duration int64,
                                urlPath string){ ... }
        func (self *Gok) SetCookie_5(name string, value string, duration int64,
                                urlPath string, domain string) { ... }
        func (self *Gok) SetCookie_7(name string, value string, duration int64,
                                urlPath string, domain string, secure bool,
                                httpOnly bool) { ... }
	```

*	__File Uploads:__ PHP equivalent of $_FILE[' ']

	```go
		// File saves the file upload to disk and returns 
		// (`file path`,`uploaded file name`, `file content type`, `size`).
		// its important that file uploads occure from form with 
		// enctype='multipart/form-data' set.
		func (self *Gok) File(name string) (string, string, string, int64) { ... }
	```

*	__Header:__ PHP equivalent of $_HEADER[' ']

	```go
		// sets the header for response.	
		// ex: gok.Header("Connection:Close")
		func (self *Gok) Header(header string) { ... }
	```

*	__Go http.Header:__ returns http.Header for Gok instance

	```go
		//go htt.Header for request
		func (self *Gok) RequestHeader() http.Header { ... }

		//go http.Header for Response
		func (self *Gok) ResponseHeader() http.Header { ... }
	``` 

*	__http.ResponseWriter__ and __*http.Request__:

	```go
		// http.ResponseWriter for gok instance
		func (self *Gok) ResponseWriter() http.ResponseWriter { ... }
		
		// *http.Request for gok instance
		func (self *Gok) HttpRequest() *http.Request { ... }
	```

Sublime text syntax:
--------------------
To get syntax higlighting for sublime text, move `gok.tmLanguage` to sublime text data directory.     
*	windows: `%APPDATA%\Sublime Text 3\Package\User\`    
*	osx:` ~/Library/Application Support/Sublime Text 3/Package/User`     
*	linux:`~/.config/sublime-text-3/Package/User`   
