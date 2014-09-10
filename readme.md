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
output
```html
<html>
	<body>
		<p>Iteration 0</p>
		<p>Iteration 1</p>
		<p>Iteration 2</p>
		<p>Iteration 3</p>
		<p>Iteration 4</p>
	</body>
</html>
````