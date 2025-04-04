# Changes

## Version 1.2

Previously, if a Host declaration or a value had trailing whitespace, that
whitespace would have been included as part of the value. This led to unexpected
consequences. For example:

```
Host example       # A comment
    HostName example.com      # Another comment
```

Prior to version 1.2, the value for Host would have been "example " and the
value for HostName would have been "example.com      ". Both of these are
unintuitive.

Instead, we strip the trailing whitespace in the configuration, which leads to
more intuitive behavior.
