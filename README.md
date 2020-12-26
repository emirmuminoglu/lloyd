# red

[![Go Reference](https://pkg.go.dev/badge/github.com/emirmuminoglu/red.svg)](https://pkg.go.dev/github.com/emirmuminoglu/red)
[![Go Report Card](https://goreportcard.com/badge/github.com/emirmuminoglu/red)](https://goreportcard.com/report/github.com/emirmuminoglu/red)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, high performance and [FastHTTP](https://github.com/valyala/fasthttp) based micro web framework.

It's heavily inspired by [atreugo](https://github.com/savsgio/atreugo)

# Features

- Routing:
  - Based on [router](https://github.com/fasthttp/router)
  - Multiple handlers to single path (like express.js)  
- High performance:
  - Uses same stack with [atreugo](https://github.com/savsgio/atreugo) so the performance is almost same. (atreugo's benchmars is availabile in [here](https://github.com/smallnest/go-web-framework-benchmark))
- Middleware support:
  - Normal middlewares
  - Defer middlewares (literally)
- Responses:
  - MarshalJSON interface support (it's very useful if you're using a marshaler other than standart encoding/json)
  - Blob JSON response
