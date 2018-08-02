# naighes/imposter [![Build Status](https://travis-ci.com/naighes/imposter.svg?branch=master)](https://travis-ci.com/naighes/imposter)

![imPOSTer Logo](https://raw.githubusercontent.com/naighes/imposter/master/readme_files/logo.png)

**imPOSTer** is a lightweight and versatile tool for the mocking of web applications.

## Source
You need `go` installed and `GOBIN` in your `PATH`. Once that is done, run the
command:
```sh
$ go get -u github.com/naighes/imposter
```

---

## Start command
Run a new instance of **imPOSTer**.

### Arguments

 * `-config-file <string>`: the configuration file path
 * `-graceful-timeout <duration>`: the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m (default 15s)
 * `-port <int>`: the listening TCP port (default 8080)
 * `-tls-cert-file-list <string>`: a comma separated list of x.509 certificates to secure communication
 * `-tls-key-file-list <string>`: a comma separated list of private key files corresponding to the x.509 certificates listed in `-tls-cert-file-list <string>`

### Example

```sh
$ ./imposter start --config-file ./config.yaml --port 3000
```

---

## Validate command
Validate and type-check any `rule_expression` within a configuration file.

### Arguments

 * `-config-file <string>`: the configuration file path
 * `-json <string>`: enable JSON output instead of plain text

### Example

```sh
$ ./imposter validate --config-file ./config.yaml --json
```
```
found 2 errors:
--------------------
could not find a built-in function with name 'eqrequest_http_method'
--------------------
could not find a parser for the current token
...
gex_match(request_url_path(), ^/myredirect$")
                              ^
```

## Configuration file

### Overview

We're going to write our first configuration now to launch an instance with a single catch-all rule:

```json
{
  "pattern_list" : [{
    "rule_expression": "${true}",
    "response": {
      "body": "Hello, default body!",
      "status_code": "${200}"
    }
  }]
}
```

A YAML parser is also available and you can write your configuration by the YAML syntax as well:

```yaml
pattern_list:
- rule_expression: ${true}
  response:
    body: Hello, default body!
    status_code: ${200}
```

`pattern_list` is a list of _rules_ defining how **imPOSTer** will handle incoming requests. Every rule requires a boolean expression. That is, if an incoming request URL matches one of the `rule_expression` the corresponding `response` is served.  

Let's suppose you need to catch all requests issued by `POST` HTTP method, every URL path containing the string `hello` and just a `Content-Type` header of type `application/json`:

```yaml
pattern_list:
- rule_expression: ${
      and(
        contains(request_url_path(), "hello"),
        eq(request_http_method(), "POST"),
        eq(request_http_header("Content-Type"), "application/json")
      )
    }
  response:
    body: Hello, complex rule!
    status_code: ${200}
```

Last but not least, you need to define the HTTP `status_code` to be returned (200 when not specified).
Note that `status_code` is an expression itself and so it needs to be wrapped into the block marker (`${...}`).
Furthermore, a positive integer value is expected from its evaluation (e.g. `status_code: ${if(contains(request_http_header("Accept-Language"), "it")) 200 else 404}`).  
Rules are tested in the order they were added to the `pattern_list` collection. If two rules match, the first one wins:

### The response object

There are two ways of defining a response object and it basically depends on the level of granularity you really need.  
For example, you can define it in a computed manner:

```yaml
pattern_list:
- rule_expression: ${regex_match(request_url_path(), "^/myredirect$")}
  response: ${redirect("http://examp.lecom/foo", 301)}
```

The above snippet shows how you can benefit from built-in functions (`redirect`) to achieve interesting results (e.g. redirecting to different URLs).  
**Note:** the computed version of the response object requires functions returning `HTTPRsp` (e.g. `link`, `redirect`, …).  

Alternatively, you can rely on a full structured version of the response object:

```json
{
  "pattern_list" : [{
    "rule_expression": "${regex_match(request_url_path(), \"^/posts$\")}",
    "response": {
      "body": "Hello, post!",
      "headers": {
        "Content-Type": "text/plain; charset=utf-8"
      },
      "status_code": "${202}"
    }
  }]
}
```

That will match the URL path `/posts` when an HTTP request will be issued by any HTTP method. The match will be handled by returning a body containing the `Hello, post!` string and just the `Content-Type` header.

### Variables

Input variables serve as parameters for built-in functions.  
Example:  

```yaml
pattern_list:
- rule_expression: ${
      and(
        regex_match(request_url_path(), "^/imposter$"),
        eq(request_http_method(), "GET")
      )
    }
  response: ${redirect(var("imposter_link", 301))}
vars:
  imposter_link: https://github.com/naighes/imposter
```

### Built-in functions

You can "combine" values with other values. These combinations are wrapped into the evaluation block marker (`${…}`), such as `${link("https://github.com/naighes/imposter")}`.  
Imposter ships with built-in functions. Functions are called with the syntax `function_name(arg1, arg2, ...)`. For example, to read a file: `${file("path.txt")}`.

#### Supported built-in functions
The supported built-in functions are:  

 * `var(name: string) -> string` - Reads the content of a variable with the specified `name` into a string.
 * `and(arg1: bool, arg2: bool, …) -> bool` - Evaluates all arguments by using the `AND` logical operator.
 * `or(arg1: bool, arg2: bool, …) -> bool` - Evaluates all arguments by using the `OR` logical operator.
 * `not(arg: bool) -> bool` - Negates its argument.
 * `eq(arg1: any, arg2: any) -> bool` - Determines whether the two specified arguments are equal.
 * `ne(arg1: any, arg2: any) -> bool` - Determines whether the two specified arguments are not equal.
 * `contains(source: string, value: string) -> bool` - Determines whether `value` substring occurs within this `source` string.
 * `request_url() -> string` - Returns the URL for the current request.
 * `request_url_path() -> string` - Returns the path component of the URL for the current request.
 * `request_url_query() -> string` - Returns any query information included in the URL for the current request.
 * `request_url_query(name: string) -> string` - Returns the first value associated with the given `name`.
 * `request_http_method() -> string` - Returns the HTTP method for the current request.
 * `request_http_host() -> string` - Returns the HTTP Host for the current request.
 * `request_http_header(name: string) -> string` - Returns the value of the HTTP header with the specified `name` for the current request.
 * `regex_match(source: string, pattern: string) -> bool` - Searches the specified `source` string for the first occurrence of the specified regular expression `pattern` and returns a value indicating whether the match is successful.
 * `file(path: string) -> string` - Reads the content of a file into a string.
 * `link(url: string) -> HTTPRsp` - Forwards a client to a new URL.
 * `redirect(url: string, status_code: int) -> HTTPRsp` - Redirects a client to a new URL with the specified `status_code` (it must be a 3XX value).
 * `in(source: array, item: string|bool|int|flota64) -> bool` - Determines whether the specified `item` exists as an element within the `source` array  object.
 * `to_string(oby: any) -> string` - Returns a string that represents the current object.

#### Conditional statements
A conditional statement identifies which statement to run based on the value of a boolean expression.  

**Syntax**:  

```
if (<boolean_expression>) <expression> else <expression>
```

**Example**:  

```yaml
pattern_list:
- rule_expression: ${regex_match(request_url_path(), "^/myfile$")}
  body: testing if statement
  headers:
    Content-Type: text/plain; charset=utf-8
    Content-Language: ${
      if(contains(request_http_header("Accept-Language"), "en"))
        "en"
      else
        "it"
    }
  status_code: ${200}
```

## License

MIT licensed. See the LICENSE file for details.
