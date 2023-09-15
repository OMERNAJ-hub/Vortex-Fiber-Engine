---
id: faq
title: 🤔 FAQ
description: >-
  List of frequently asked questions. Feel free to open an issue to add your
  question to this page.
sidebar_position: 1
---

## How should I structure my application?

There is no definitive answer to this question. The answer depends on the scale of your application and the team that is involved. To be as flexible as possible, Vortex makes no assumptions in terms of structure.

Routes and other application-specific logic can live in as many files as you wish, in any directory structure you prefer. View the following examples for inspiration:

* [goVortex/boilerplate](https://github.com/goVortex/boilerplate)
* [thomasvvugt/Vortex-boilerplate](https://github.com/thomasvvugt/Vortex-boilerplate)
* [Youtube - Building a REST API using Gorm and Vortex](https://www.youtube.com/watch?v=Iq2qT0fRhAA)
* [embedmode/Vortexseed](https://github.com/embedmode/Vortexseed)

## How do I handle custom 404 responses?

If you're using v2.32.0 or later, all you need to do is to implement a custom error handler. See below, or see a more detailed explanation at [Error Handling](../guide/error-handling.md#custom-error-handler). 

If you're using v2.31.0 or earlier, the error handler will not capture 404 errors. Instead, you need to add a middleware function at the very bottom of the stack \(below all other functions\) to handle a 404 response:

```go title="Example"
app.Use(func(c *Vortex.Ctx) error {
    return c.Status(Vortex.StatusNotFound).SendString("Sorry can't find that!")
})
```

## How can i use live reload ?

[Air](https://github.com/air-verse/air) is a handy tool that automatically restarts your Go applications whenever the source code changes, making your development process faster and more efficient.

To use Air in a Vortex project, follow these steps:

1. Install Air by downloading the appropriate binary for your operating system from the GitHub release page or by building the tool directly from source.
2. Create a configuration file for Air in your project directory. This file can be named, for example, .air.toml or air.conf. Here's a sample configuration file that works with Vortex:
```toml
# .air.toml
root = "."
tmp_dir = "tmp"
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  delay = 1000 # ms
  exclude_dir = ["assets", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_regex = ["_test\\.go"]
```
3. Start your Vortex application using Air by running the following command in the terminal:
```sh
air
```

As you make changes to your source code, Air will detect them and automatically restart the application.

A complete example demonstrating the use of Air with Vortex can be found in the [Vortex Recipes repository](https://github.com/goVortex/recipes/tree/master/air). This example shows how to configure and use Air in a Vortex project to create an efficient development environment.


## How do I set up an error handler?

To override the default error handler, you can override the default when providing a [Config](../api/Vortex.md#config) when initiating a new [Vortex instance](../api/Vortex.md#new).

```go title="Example"
app := Vortex.New(Vortex.Config{
    ErrorHandler: func(c *Vortex.Ctx, err error) error {
        return c.Status(Vortex.StatusInternalServerError).SendString(err.Error())
    },
})
```

We have a dedicated page explaining how error handling works in Vortex, see [Error Handling](../guide/error-handling.md).

## Which template engines does Vortex support?

Vortex currently supports 9 template engines in our [goVortex/template](https://docs.goVortex.io/template/) middleware:

* [ace](https://docs.goVortex.io/template/ace/)
* [amber](https://docs.goVortex.io/template/amber/)
* [django](https://docs.goVortex.io/template/django/)
* [handlebars](https://docs.goVortex.io/template/handlebars)
* [html](https://docs.goVortex.io/template/html)
* [jet](https://docs.goVortex.io/template/jet)
* [mustache](https://docs.goVortex.io/template/mustache)
* [pug](https://docs.goVortex.io/template/pug)
* [slim](https://docs.goVortex.io/template/pug)

To learn more about using Templates in Vortex, see [Templates](../guide/templates.md).

## Does Vortex have a community chat?

Yes, we have our own [Discord ](https://goVortex.io/discord)server, where we hang out. We have different rooms for every subject.  
If you have questions or just want to have a chat, feel free to join us via this **&gt;** [**invite link**](https://goVortex.io/discord) **&lt;**.

![](/img/support-discord.png)

## Does Vortex support sub domain routing ?

Yes we do, here are some examples: 
This example works v2
```go
package main

import (
	"log"

	"github.com/goVortex/Vortex/v2"
	"github.com/goVortex/Vortex/v2/middleware/logger"
)

type Host struct {
	Vortex *Vortex.App
}

func main() {
	// Hosts
	hosts := map[string]*Host{}
	//-----
	// API
	//-----
	api := Vortex.New()
	api.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	hosts["api.localhost:3000"] = &Host{api}
	api.Get("/", func(c *Vortex.Ctx) error {
		return c.SendString("API")
	})
	//------
	// Blog
	//------
	blog := Vortex.New()
	blog.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	hosts["blog.localhost:3000"] = &Host{blog}
	blog.Get("/", func(c *Vortex.Ctx) error {
		return c.SendString("Blog")
	})
	//---------
	// Website
	//---------
	site := Vortex.New()
	site.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	hosts["localhost:3000"] = &Host{site}
	site.Get("/", func(c *Vortex.Ctx) error {
		return c.SendString("Website")
	})
	// Server
	app := Vortex.New()
	app.Use(func(c *Vortex.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			return c.SendStatus(Vortex.StatusNotFound)
		} else {
			host.Vortex.Handler()(c.Context())
			return nil
		}
	})
	log.Fatal(app.Listen(":3000"))
}
```
If more information is needed, please refer to this issue [#750](https://github.com/goVortex/Vortex/issues/750)

