# flagger
`flagger` is a simple package wrapping `flag` standard package allowing the declaration of flags via a structure.

# Installation

```
go get -u github.com/inozuma/flagger
```

# Usage

## Simple usage

Supposing we have a `Configuration`:

```go
type Configuration struct {
    Enable  bool   `flag:"enable,false,enable something"`
    Address string `flag:"addr,:8080,address to something"`
}
```

You can call `flagger.Flag` to declare the flag `enable`:

```go
config := &Configuration{}

if err := flagger.Flag(config); err != nil {
    panic(err)
}
```

And like `flag`, parse your flags with `flag.Parse`:

```go
flag.Parse()
```
