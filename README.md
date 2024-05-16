## multiplexer - A Versatile HTTP Request Router🚀

The multiplexer library provides a robust and customizable HTTP request router built on top of the standard net/http package in Go. It empowers developers to efficiently manage incoming requests, directing them to appropriate handlers based on URL patterns and methods.
## Installation🛠

Prerequisites: Go 1.22 or above


```bash
  go get -u github.com/iarsham/multiplexer
```

## Usage/Examples💡

```go
func main() {
	mux := multiplexer.New(http.NewServeMux(), "/api")
	mux.NotFound = http.HandlerFunc(notfound)
	mux.MethodNotAllowed = http.HandlerFunc(allowed)
	dynamic := multiplexer.NewChain(logMiddleware)
	mux.Handle("GET /root", dynamic.WrapFunc(root))
	authGroup := mux.Group("/user")
	protected := dynamic.Append(authMiddleware)
	authGroup.Handle("GET /home", protected.WrapFunc(home))
	log.Fatal(http.ListenAndServe(":8000", mux))
}
```


## Features🪜

- Clear and Consistent Routing: Define routes using URL patterns that may include named capture groups (e.g., /users/:id).
- Base Path Support: Configure a base path to be prepended to all route patterns, simplifying organization within nested routing scenarios.
- Flexible Handler Registration: Register handlers using HandleFunc, Handle, or custom handler functions.
- Customizable Not Found and Method Not Allowed Handlers: Provide tailored responses for unmatched requests or unsupported methods.
- Sub-Routing with Grouping: Create nested routing hierarchies using the Group function, promoting better code organization.

## Contributing 🤝

Contributions are always welcome!

- Create a fork of the repository.

- Make your changes in a separate branch.

- Ensure your code adheres to Go's formatting and style conventions.

- Add unit tests for your changes.

- Submit a well-structured pull request with a clear description of your changes.
  ❤️‍🔥
## Contributors 👨🏻‍💻👩🏼‍💻

- [iarsham](https://www.github.com/iarsham)


## License ⚠️

[MIT](https://choosealicense.com/licenses/mit/)

