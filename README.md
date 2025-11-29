# TeleCraft Framework

A lightweight Telegram Bot framework built on top of `telegram-bot-api` with **routing**, **middleware support**, and **monitoring**. Designed for scalable bots with easy configuration and flexible route handling.

---

## Features

- **Routing**: Register commands, messages, or callback handlers.
- **Middlewares**: Apply global or route-specific middlewares.
- **Monitoring**: Track bot usage, performance, and errors.
- **Flexible Initialization**: Configure bot using `TeleCraftOptions`.
- **Structured Responses**: Handlers return `ResponseHandlerFunc` for advanced control.

---

## Installation

```bash
go get github.com/yourusername/telecraft
```

---

## Getting Started

### Initialize Bot

```go
import "github.com/yourusername/telecraft"

options := telecraft.TeleCraftOptions{
    RepoType:      "memory",       // storage type
    DefaultRoute:  "start",        // default route
    maxGoroutines: 10,             // max concurrent handlers
    Timeout:       30,             // timeout in seconds
    Token:         "YOUR_BOT_TOKEN",
    errorMessage:  "Something went wrong",
}

bot := telecraft.New(options)
```

---

### Register Routes

```go
// Register a route with optional middlewares
bot.Router.Register("/start", func(ctx *telecraft.Context) (*telecraft.ResponseHandlerFunc, error) {
    return &telecraft.ResponseHandlerFunc{
        Path:         "welcome",
        MessageConfigs: []*tgbotapi.MessageConfig{
            tgbotapi.NewMessage(ctx.UserID, "Welcome!"),
        },
        ReleaseState: true,
    }, nil
}, loggingMiddleware)
```

---

### Global Middlewares

```go
// Apply middlewares globally
bot.Router.SetGlobalMiddlewares(loggingMiddleware, authMiddleware)
```

---

## Types Overview

- **Context**: Holds incoming `tgbotapi.Update`, user info, params, and extra data.
- **HandlerFunc**: `func(*Context) (*ResponseHandlerFunc, error)`
- **Middleware**: `func(HandlerFunc) HandlerFunc`
- **ResponseHandlerFunc**: Controls responses, routing, state release, and message configs.

---

## Middlewares Example

```go
func loggingMiddleware(next telecraft.HandlerFunc) telecraft.HandlerFunc {
    return func(ctx *telecraft.Context) (*telecraft.ResponseHandlerFunc, error) {
        log.Println("Incoming message:", ctx.Message.Text)
        return next(ctx)
    }
}
```

---

## Monitoring

Built-in monitoring tracks bot performance, errors, and usage.
Supports integration with logging or metrics systems.

---

## Contributing

Contributions are welcome! Open issues or submit pull requests for bug fixes or new features.

---
