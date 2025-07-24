# httpx – a minimalistic *Go* library for building HTTP APIs



[![Go Reference](https://pkg.go.dev/badge/github.com/pda-labs/httpx.svg)](https://pkg.go.dev/github.com/pda-labs/httpx) [![Go Report Card](https://goreportcard.com/badge/github.com/pda-labs/httpx)](https://goreportcard.com/report/github.com/pda-labs/httpx) [![License](https://img.shields.io/github/license/pda-labs/httpx.svg)](LICENSE)

**Unified response format + multilingual validation + full set of 2xx / 3xx / 4xx / 5xx shortcuts**

---


- [httpx – a minimalistic *Go* library for building HTTP APIs](#httpx--a-minimalistic-go-library-for-building-http-apis)
  - [What is it?](#what-is-it)
  - [Philosophy](#philosophy)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [Response Format](#response-format)
  - [Built-in HTTP Response Helpers](#built-in-http-response-helpers)
    - [2xx – Success](#2xx--success)
    - [3xx – Redirects](#3xx--redirects)
    - [4xx – Client Errors](#4xx--client-errors)
    - [5xx – Server Errors](#5xx--server-errors)
  - [Multilingual Validation](#multilingual-validation)
  - [License](#license)

---


## What is it?

`httpx` is a lightweight middleware-style library built on top of `net/http`, `chi`, and `github.com/go-playground/validator` that solves  
**three** timeless problems in REST API development:

1. **Standardizes all responses** (both success and error) using a common `Envelope`.
2. **Simplifies validation** of JSON requests + localizes messages (12 languages).
3. **Provides ready-to-use helpers** for _all_ HTTP status codes – from `200 OK` to `511 NETWORK_AUTH_REQUIRED`.

No more copy-pasting `w.Header().Set("Content-Type")`, no more endless `switch status` – just call `httpx.Ok()` or `httpx.ErrorBadRequest()`.

---

## Philosophy

httpx was built to reduce boilerplate and enforce consistency in HTTP API design.
- Strict response contract via unified Envelope
- Human-friendly error messages in multiple languages
- Zero-config, but highly extensible

---

## Installation

```
go get github.com/pda-labs/httpx
```

---

## Quick Start

```go
package main

import (
  "net/http"

  "github.com/go-chi/chi/v5"
  "github.com/your-org/httpx/v1"
)

type SignupDTO struct {
  Email string `json:"email" validate:"required,email"`
}

func main() {
  // Create validator instance
  validatorInstance := validator.New()
  // Initialize httpx
  httpx.Init(validatorInstance)

  r := chi.NewRouter()

  // chi sets X-Request-ID
  r.Use(middleware.RequestID)

  r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
    var dto SignupDTO

    // Line = parsing + validation + localization
    if det, err := httpx.BindValidate(r, &dto); err != nil {
      if det != nil {
        httpx.ErrorValidation(w, r, det) // 400 + details
      } else {
        httpx.ErrorBadRequest(w, r, err.Error())
      }
      return
    }

    // ... business logic ...

    httpx.Created(w, r, "/profile", map[string]string{"status": "ok"})
  })

  http.ListenAndServe(":8080", r)
}
```

---

## Response Format

```json
{
  "success": true,
  "data": { … },
  "error": null,
  "trace_id": "8f8e…"
}
```

**Error:**

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION",
    "message": "Request failed validation",
    "details": { "email": "must be a valid email address" }
  },
  "trace_id": "8f8e…"
}
```

---

## Built-in HTTP Response Helpers

### 2xx – Success

| Code | Helper Function    | Description                              |
| ---- | ------------------ | ---------------------------------------- |
| 200  | `Ok`               | Standard success with data               |
| 201  | `Created`          | Resource created + `Location` header     |
| 202  | `Accepted`         | Async queue accepted                     |
| 203  | `NonAuthoritative` | Fetched from third-party source          |
| 204  | `NoContent`        | Success without body (e.g. DELETE, ping) |
| 205  | `ResetContent`     | Client should reset form/UI              |
| 206  | `PartialContent`   | Partial file response (`Range:`)         |

### 3xx – Redirects

| Code | Helper Function            | Description                        |
| ---- | -------------------------- | ---------------------------------- |
| 300  | `RedirectMultipleChoices`  | Multiple representations           |
| 301  | `RedirectMovedPermanently` | Permanent redirect (SEO)           |
| 302  | `RedirectFound`            | Temporary, method is reset         |
| 303  | `RedirectSeeOther`         | POST → GET after resource creation |
| 304  | `RedirectNotModified`      | ETag/cache unchanged               |
| 307  | `RedirectTemporary`        | Temporary, preserves method/body   |
| 308  | `RedirectPermanent`        | Permanent, preserves method/body   |

### 4xx – Client Errors

| Code | Helper Function             | Description                             |
| ---- | --------------------------- | --------------------------------------- |
| 400  | `ErrorBadRequest`           | Unparsable JSON or bad request          |
| 400  | `ErrorValidation`           | JSON is valid but breaks business rules |
| 401  | `ErrorUnauthorized`         | Missing or invalid auth                 |
| 402  | `ErrorPaymentRequired`      | Quota exceeded or unpaid plan           |
| 403  | `ErrorForbidden`            | Auth valid, but not enough rights       |
| 404  | `ErrorNotFound`             | Endpoint or resource not found          |
| 405  | `ErrorMethodNotAllowed`     | Method not supported                    |
| 406  | `ErrorNotAcceptable`        | Unsupported `Accept` type               |
| 407  | `ErrorProxyAuthRequired`    | Proxy requires authentication           |
| 408  | `ErrorRequestTimeout`       | Client took too long to send data       |
| 409  | `ErrorConflict`             | Conflict (e.g. uniqueness, versioning)  |
| 410  | `ErrorGone`                 | Resource permanently deleted            |
| 411  | `ErrorLengthRequired`       | Missing `Content-Length` header         |
| 412  | `ErrorPreconditionFailed`   | ETag or `If-Match` check failed         |
| 413  | `ErrorPayloadTooLarge`      | Body exceeds max size                   |
| 414  | `ErrorURITooLong`           | URL too long                            |
| 415  | `ErrorUnsupportedMediaType` | Unsupported `Content-Type`              |
| 416  | `ErrorRangeNotSatisfiable`  | File range out of bounds                |
| 417  | `ErrorExpectationFailed`    | `Expect: 100-continue` failed           |
| 418  | `ErrorTeapot`               | I’m a teapot                            |
| 421  | `ErrorMisdirectedRequest`   | Invalid host (HTTP/2 SNI issue)         |
| 422  | `ErrorUnprocessableEntity`  | Semantically invalid input              |
| 423  | `ErrorLocked`               | Resource is locked                      |
| 424  | `ErrorFailedDependency`     | Previous batch step failed              |
| 425  | `ErrorTooEarly`             | 0-RTT request too early                 |
| 426  | `ErrorUpgradeRequired`      | Requires WebSocket or protocol upgrade  |
| 428  | `ErrorPreconditionRequired` | `If-Match` required                     |
| 429  | `ErrorTooManyRequests`      | Rate limit exceeded                     |
| 431  | `ErrorHeaderFieldsTooLarge` | Headers too large                       |
| 451  | `ErrorLegalReasons`         | Blocked for legal reasons               |

### 5xx – Server Errors

| Code | Helper Function                | Description                                |
| ---- | ------------------------------ | ------------------------------------------ |
| 500  | `ErrorInternal`                | Panic or unhandled internal error          |
| 501  | `ErrorNotImplemented`          | Documented feature, not implemented        |
| 502  | `ErrorBadGateway`              | Upstream returned error                    |
| 503  | `ErrorServiceUnavailable`      | Service temporarily offline (deploy, load) |
| 504  | `ErrorTimeout`                 | Gateway timeout                            |
| 505  | `ErrorHTTPVersionNotSupported` | HTTP version not supported                 |
| 506  | `ErrorVariantAlsoNegotiates`   | Loop in content negotiation                |
| 507  | `ErrorInsufficientStorage`     | Disk or quota exceeded                     |
| 508  | `ErrorLoopDetected`            | Cyclic reference (e.g. WebDAV)             |
| 510  | `ErrorNotExtended`             | Missing protocol extension                 |
| 511  | `ErrorNetworkAuthRequired`     | Captive portal — network login required    |

> Note: Every status code has a dedicated shortcut function. It automatically builds an Envelope, attaches the `trace_id`, and sets all headers correctly.

---

## Multilingual Validation

| Code | Language                                       |
| ---- | ---------------------------------------------- |
| en   | English                                        |
| ru   | Russian                                        |
| de   | German                                         |
| lv   | Latvian (**TODO: waiting for native version**) |
| zh   | Chinese                                        |
| fr   | French                                         |
| es   | Spanish                                        |
| it   | Italian                                        |
| pt   | Portuguese                                     |
| ja   | Japanese                                       |
| ko   | Korean                                         |

> Language selection logic (in order of priority):
> 1. `X-Request-Lang: ru`  
> 2. `Accept-Language: ru-RU,ru;q=0.9`  
> 3. fallback - English

## License

MIT © 2025 Pobedinskiy David Arturovich