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

installation
------------
`$ go get https://github.com/YuaShizuki/gok.git`

gok tags
--------
1.	**`<?go ?>`**  
	can scope any valid go code. Code scoped in this tag 
	executes inside a render function. 

2.	**`<?gofn funcName( args ) ( return types ) { go code }?>`**  
	allows you to define a global function thats accessible by all ".gok" scripts.
	functions defined in this tag cannot be methods of objects or types.
	hence:    	
	`<?gofn (type) funcName( args ) (return types) { } ?>`   
	would be incorect.

3. **`<?goimp ?>`**  
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
4. **`<?gouse ?>`**   
	gouse tag allows to define global variables, functions   

	```text
		<?gouse
			type nstring string 
			var x nstring = "Hello World"
			func (self nstring) () { 
				fmt.Println(self)
			}
			...
		?>    
	```
5. **`<?go@fn funcName(args []string) ([]string, error) { go code } ?>`**
