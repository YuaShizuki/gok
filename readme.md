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
	pritnts basic api useable inside of `<?go ?>` tag.
