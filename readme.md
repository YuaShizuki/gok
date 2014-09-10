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
---------------
`<?go ?>` can scope any valid go code. Code scoped in this tag
executes inside a render function. 

`<?gofn functionName(){ }?>` allows you to define a global function
thats accessible by all .gok scripts

