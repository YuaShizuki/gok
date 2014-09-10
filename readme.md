gok
===

this is my creation to build web apps in go, php style.
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
============
`go get https://github.com/YuaShizuki/gok.git`