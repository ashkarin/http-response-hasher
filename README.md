HTTP Response Hasher
=====================

Tool to request URLs and calculate MD5 hashes of the received responses
in parallel.

# Disclaimer #
This is a toy project made for fun and it will not be supported.
It is implemented using only the standard Go library.

# Build and Test #

To build the project
```
make build
```

To run tests
```
make test
```

# Usage #

The tool has only a `--parallel` parameter to control the number of workers
processing the URLs in parallel.

```
http-response-hasher --parallel 2 google.com http://mail.ru
```

## DEMO ##
[![asciicast](https://asciinema.org/a/keGyVFYlpMyLNCIGgWNzQ8LV2.svg)](https://asciinema.org/a/keGyVFYlpMyLNCIGgWNzQ8LV2)
