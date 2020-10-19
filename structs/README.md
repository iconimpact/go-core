# Structs

Structs contains various utilities to work with Go structs. It's basically a high level package based on primitives from the reflect package.

Feel free to add new functions or improve the existing code.

## Install

```bash
go get github.com/iconmobile-dev/go-core/structs
```

## Usage and Examples

`structs.Populate` sets fields from 'b' into 'a' struct by matching name and type.

#### Options
- a tag is in the form of: `populate:"name"`
- use pointers in 'b' fields to pass empty values
- use tag to match field from 'b' to 'a'
- ignore fields from 'b' with `populate:"-"` tag


```go
a := struct {
    Name string
    Description string
    Status int
    ID int
    FileName string
}{
    Name: "John",
    Description: "pointer in b, skip",
    Status: 10,
}

b := struct {
    Name *string
    Description *string
    Status int
    ID int `populate:"-"`
    Filename string `populate:"FileName"`
}{
    Name: "John Wick",
    ID: 20,
    Filename: "image.png",
}

err := structs.Populate(&a, b)

result:
{Name:"John Wick", Description:"pointer in b, skip", Status:0, ID: 0, FileName: "image.png"}
```

`structs.Sanitize` removes all leading and trailing white spaces from struct exported string fields.

#### Options
- a tag is in the form of: `sanitize:"option"`
- ignore fields with `sanitize:"-"` tag
- clean and validate email with `sanitize:"email"` tag

```go
a := struct {
    Name string
    Status int
    Password string `sanitize:"-"`
    Email string `sanitize:"email"`
}{
    Name: " John Wick ",
    Status: 10,
    Password: "my password ",
    Email: " email.string+me@googlemail.com",
}

err := structs.Sanitize(&a)

result:
{Name:"John Wick", Status:0, Password:"my password ", Email:"email.string+me@gmail.com"}
```