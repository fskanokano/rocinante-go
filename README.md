## Rocinante

#### Rocinante is a gin inspired web framework built on top of net/http.

## ⚙️ Installation

```bash
$ go get -u github.com/fskanokano/rocinante-go
```

## ⚡️ Quickstart

```go
package main

import (
	"log"

	"github.com/fskanokano/rocinante-go"
)

func main() {
	r := rocinante.Default()
	r.GET("/", func(c *rocinante.Context) {
		c.String("hello world")
	})
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

## 👀 Examples

Some common examples are listed below.

#### 📖 Param

```go
func main() {
	r := rocinante.Default()
	r.GET("/param1/:name/:id", func(c *rocinante.Context) {
		name := c.Param("name")
		id := c.Param("id")
		c.JSON(rocinante.Map{
			"name": name,
			"id":   id,
		})
	})
	r.GET("/param2/file/*file_path", func(c *rocinante.Context) {
		filePath := c.Param("file_path")
		c.JSON(rocinante.Map{
			"file_path": filePath,
		}, http.StatusOK)
	})
}
```

```bash
$ curl http://127.0.0.1:8000/param1/kano/1
# {"name":"kano","id":"1"}

$ curl http://127.0.0.1:8000/param2/file/one/two/three.xxx
# {"file_path":“/one/two/three.xxx”}
```

#### 📖 Serving Static File

```go
func main() {
	r := rocinante.Default()
	r.Static("/image", "image")
}
```

#### 📖 Middleware

```go
func main() {
	r := rocinante.Default()
	//global middleware
	r.Use(func(c *rocinante.Context) {
		fmt.Println("before")
		c.Next()
		fmt.Println("after")
	})
	//specific middleware on handler
	r.GET("/foo/:foo_id",
		func(c *rocinante.Context) {
			fooID := c.Param("foo_id")
			if fooID != "1" {
				c.AbortWithJSON(rocinante.Map{
					"error": "foo_id is not 1",
				}, http.StatusBadRequest)
				return
			}
			c.Next()
		},
		func(c *rocinante.Context) {
			c.String("foo")
		})
}
```

#### 📖 Model binding and validation

Rocinante uses [**go-playground/validator/v10**](https://github.com/go-playground/validator) for validation. Check the full docs on tags usage [here](https://godoc.org/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags).

```go
type Login struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password" validate:"required"`
}

func main() {
	r := rocinante.Default()
	//bind json
	r.POST("/json_login", func(c *rocinante.Context) {
		var loginJSON Login
		if err := c.BindJSON(&loginJSON); err != nil {
			c.JSON(rocinante.Map{
				"error": err.Error(),
			}, http.StatusBadRequest)
			return
		}
		c.String("login success")
	})
	//bind query
	r.POST("/query_login", func(c *rocinante.Context) {
		var loginQuery Login
		if err := c.BindQuery(&loginQuery); err != nil {
			c.JSON(rocinante.Map{
				"error": err.Error(),
			}, http.StatusBadRequest)
			return
		}
		c.String("login success")
	})
	//bind form
	r.POST("/form_login", func(c *rocinante.Context) {
		var loginForm Login
		if err := c.BindForm(&loginForm); err != nil {
			c.JSON(rocinante.Map{
				"error": err.Error(),
			}, http.StatusBadRequest)
			return
		}
		c.String("login success")
	})
}
```

#### 📖 WebSocket

```go
func main() {
	r := rocinante.Default()
	r.WebSocket("/test",
		func(conn *websocket.Conn) {
			for {
				mt, data, err := conn.ReadMessage()
				if err != nil {
					break
				}
				fmt.Println("received message: " + string(data))
				err = conn.WriteMessage(mt, data)
				if err != nil {
					break
				}
			}
		},
		//use specific middleware on websocket handler
		func(c *rocinante.Context) {
			fmt.Println("before")
			c.Next()
			fmt.Println("after")
		})
}
```

#### 📖 MVC

```go
type TestController struct {
	*rocinante.Controller
}

func (t *TestController) GET(c *rocinante.Context) {
	c.String("mvc get")
}

func (t *TestController) POST(c *rocinante.Context) {
	c.String("mvc post")
}

func main() {
	r := rocinante.Default()
	r.Route("/mvc", &TestController{})
}
```

#### 📖 Group

The use of group and router is the same, group can generate unlimited subgroups, the same level of group does not affect each other, the subgroup will inherit some properties of the parent group (url prefix, middleware...).

```go
func main() {
	r := rocinante.Default()
	{
		v1 := r.Group("/v1")
		v1.Use(func(c *rocinante.Context) {
			fmt.Println("v1 middleware")
			c.Next()
		})
		v1.GET("/index", func(c *rocinante.Context) {
			c.String("v1 index")
		})
		//curl http://127.0.0.1:8000/v1/index

		v2 := v1.Group("/v2")
		v2.Use(func(c *rocinante.Context) {
			fmt.Println("v2 middleware")
		})
		v2.GET("/index", func(c *rocinante.Context) {
			c.String("v2 index")
		})
		//curl http://127.0.0.1:8000/v1/v2/index
	}
}
```

#### 📖 CORS

```go
func main() {
	r := rocinante.Default()
	r.Use(cors.New(cors.Option{
		AllowOrigins:     []string{"www.example.com"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders:     []string{"Custom-Header"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Exposed-Header"},
		MaxAge:           3600,
	}))
}
```
