# naighes/imposter

![Imposter Logo](https://raw.githubusercontent.com/naighes/imposter/master/readme_files/logo.png)

**Imposter** is a lightweight and versatile tool for the mocking of web applications.

## Source
You need `go` installed and `GOBIN` in your `PATH`. Once that is done, run the
command:
```sh
$ go get -u github.com/naighes/imposter
```

---

## Start command
```console
  -port int
      TCP port the server is listening at
  -config-file
      Configuration file (see documentation below)
```

---

## Configuration file

### Overview

We're going to write our first configuration now to launch an instance with a single catch-all rule:

```json
{
  "pattern_list" : [{
    "pattern": "^.*$",
    "response": {
      "body": "Hello, default body!",
      "status_code": 200
    }
  }]
}
```

A YAML parser is available as well and you can write your configuration by YAML syntax as well:

```yaml
pattern_list:
- pattern: "^.*$"
  response:
    body: Hello, default body!
    status_code: 200
```

`pattern_list` is a list of _rules_ defining how **imposter** will handle incoming requests. Every rule requires a `pattern` at least. That is, if an incoming request URL matches one of the regex patterns the corrisponding `response` is served.  
By default every rule matches against every HTTP method, but you can overwrite that by using the `method` property:

```json
{
  "pattern_list" : [{
    "pattern": "^.*$",
    "method": "POST",
    "response": {
      "body": "Hello, default body!",
      "status_code": 200
    }
  }]
}
```

Last but not least, you need to define the HTTP `status_code` to be returned (200 when not specified).  
Rules are tested in the order they were added to the `pattern_list` collection. If two rules match, the first one wins:

### The response object

There are two ways of definng a response object and it basically depends on the level of granularity you really need.  
For example, you can define it in a computed manner:

```yaml
pattern_list:
- pattern: "^/myredirect$"
  response: ${redirect("http://examp.lecome/foo")
```

The above snippet shows how you can benefit from built-in functions (`redirect`) to achieve interesting results (e.g. redirecting to different URLs).  
Note: the computed version of response object requires functions returning `HttpRsp` (e.g. `link`, `redirect`, …).  

Alternatively, you can rely on a full structured version of the response object:

```json
{
  "pattern_list" : [{
    "pattern": "^/posts$",
    "method": "POST",
    "response": {
      "body": "Hello, post!",
      "headers": {
        "Content-Type": "text/plain; charset=utf-8"
      },
      "status_code": 202
    }
  }]
}
```

That will match the URL path `/posts` when an HTTP request will be issued by the `POST` HTTP method. The match will be handled by returning a body containing the `Hello, post!` string and just the `Content-Type` header.

### Variables

Input variables serve as parameters for built-in functions.  
Example:  

```yaml
pattern_list:
- pattern: "^/posts$"
  method: POST
  response: ${redirect(var("imposter_link"))
vars:
  imposter_link: http://localhost:8080/aaa
```

### Built-in functions

Emedded within strings you can interpolate other values. These interpolations are wrapped in `${}`, such as `${link("https://github.com/naighes/imposter")}`.  
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
 * `request_url() -> string` - Returns the URL of the current request.
 * `http_header(name: string) -> string` - Returns the value of the HTTP header with the specified `name`.
 * `regex_match(source: string, pattern: string) -> bool` - Searches the specified `source` string for the first occurrence of the specified regular expression `pattern` and returns a value indicating whether the match is successful.
 * `file(path: string) -> string` - Reads the content of a file into a string.
 * `link(url: string) -> HttpRsp` - Forwards a client to a new URL.
 * `redirect(url: string) -> HttpRsp` - Redirects a client to a new URL by using a 301 status code.

#### Conditional statements
A conditional statement identifies which statement to run based on the value of a boolean expression.  

**Syntax**:  

```
if (<boolean_expression>) <expression> else <expression>
```

**Example**:  

```yaml
pattern_list:
- pattern: "^/myfile$"
  body: testing if statement
  headers:
    Content-Type: text/plain; charset=utf-8
    Content-Language: ${if(contains(http_header("Accept-Language"), "en")) "en" else "it"}
  status_code: 200
```

## License

MIT licensed. See the LICENSE file for details.
