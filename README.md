# HTTP Pipe

Collects and runs multiple http.Handler filters until the first one writes a
response.

```go
pipe := httppipe.New(
  invalidcookiedropper.New(),
  contenttypecleaner.New()
  requestid.New(),
  myapp.New(),
)

http.Handle("/", pipe)
```
