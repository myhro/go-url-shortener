Go URL Shortener
================

[![Build Status](https://travis-ci.org/myhro/go-url-shortener.svg?branch=master)](https://travis-ci.org/myhro/go-url-shortener)
[![Go Report Card](https://goreportcard.com/badge/github.com/myhro/go-url-shortener)](https://goreportcard.com/report/github.com/myhro/go-url-shortener)

URL Shortener API written in Go.

## Methods

```
resource URL {
    id: number
    hash: string
    url: string
}
```

* `GET /{hash}`: redirects to full URL.
* `GET /{hash}/details`: returns URL details.
* `POST /`: creates a new entry from a URL resource.
