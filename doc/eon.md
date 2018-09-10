# EON (Extensible Object Notation)

EON is a richly-typed, micro-language that is highly suited for creating DSLs.
It takes inspiration from [Dhall], [REBOL]/[RED], and [UCL]/[HCL2], with a
couple of features taken from [ES6] and [YAML].

```hcl
format = "EON"
created = 2018-09-01
version = 0.0.1

author {
    name = "tav"
    addictions = ["Gauloises", "Nutella"]
}

print `${format} was created by ${author.name}, ${days.since created} days ago`
```

## Execution Modes

EON supports two distinct execution modes:

- Static
- Dynamic

In static mode, the interpreted state of EON is mapped to datatypes that are
provided by the caller, e.g. in the standard Go reference implementation, this
is accessible via the familiar:

```go
Marshal(v interface{}) ([]byte, error)
Unmarshal(data []byte, v interface{}) error
```

## Primitive Types

[dhall]: https://github.com/dhall-lang/dhall-lang
[es6]: https://github.com/lukehoban/es6features
[hcl2]: https://github.com/hashicorp/hcl2
[rebol]: https://en.wikipedia.org/wiki/Rebol
[red]: https://www.red-lang.org/
[ucl]: https://github.com/vstakhov/libucl
[yaml]: http://yaml.org/
