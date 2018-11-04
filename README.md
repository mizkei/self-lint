# self-lint

self-lint is a small linter for Go code.

# install

```
$ go get -u github.com/mizkei/self-lint
```

# usage

Write what you want to prohibit in your project.

- import
  import packages
- ref
  reference to variables, methods, or const
- write
  write builtin function or statements

config file

```yaml
prohibited-matter:
  - target:
      - global
    import:
      - github.com/mizkei/self-lint/_example/test
  - target:
      - github.com/mizkei/self-lint/_example/model
    ref:
      github.com/mizkei/self-lint/_example/time:
        - SetTime
        - ResetTime
  - target:
      - github.com/mizkei/self-lint/_example/data
    write:
      - if
      - switch
      - panic
```

run self-lint
```
$ self-lint -config=config_sample.yml ./...
data/data.go:12:2: forbidden to write if
data/data.go:13:3: forbidden to write panic
model/model.go:4:2: forbidden to import github.com/mizkei/self-lint/_example/test
model/model.go:22:7: forbidden to refer 'github.com/mizkei/self-lint/_example/time'.SetTime
```
