# Contextualized errors for go

> Errors with custom additional data, easy comparison and no dependencies 

Description
-----------

One of the greatest challenges for golang (at least, before go2 release) is the error handling question.
Someone says, its ugly. Someone says, its uncomfortable.
And I say that its limited: standard error interface provides only the `Error()` method which provides us the string. 
That is bad because such poorness guides us to comparison issues: how can we know which error did we receive?

* We can compare strings (its about standard examples with their `errors.New("my err:" + err.Error())`) - but strings comparison is the imprecise and obviously doubtful way;
* We can compare variables (other standard way with global initialization `var ErrBadSomething = errors.New("bad something")` and comparison seems like `if err == ErrBadSomething {...}`) - but that prevents us from additional data like scope, stacktrace and etc;
* We can define custom error structures for every our package and use error code constants - yep, it will work! But as we can see, there is almost no real examples of this approach. And that's clear: its too verbose and massive approach.  

But we also have another good example of universal dynamic package, its the standard `context.Context`. So, inspired by it, I decided to write a little package which will provide us another error handling way: we can contextualize our errors to populate them with error codes and other necessary like basic key-valued context.

Usage
-----------

Create an error
```
err := errors.New("my error message")
```
Populate an error
```
err = errors.New("invalid user").WithValue("user_id", 15)
```
Get your custom value
```
userId = err.Value("user_id")
```
Use your custom error codes
```
type errKey int8
const (
    ErrorCodeKey errKey = iota
	ErrorA
	ErrorB
	ErrorC
)
err = errors.New("invalid user").WithValue(ErrorCodeKey, ErrorA)
```
Compare your custom error code
```
if err.Value(ErrorCodeKey) == ErrorA {...}
// or
if errors.Value(err, ErrorCodeKey) == ErrorA {...}
```
And etc.        

Installation
-----------
As every other golang library `go get -u github.com/eyudkin/ctx-errors`