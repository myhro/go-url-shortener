Go URL Shortener
================

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
