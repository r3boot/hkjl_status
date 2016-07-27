## hkjl_status
Small tool which polls all Hackenkunjeleren services and generates a html
page displaying the status.

## installation
The procedure below works for Debian Jessie. If you have any other OS, you're
on your own

```
$ sudo apt-get install git golang binutils
$ git clone https://github.com/r3boot/hkjl_status
$ cd hkjl_status
$ ./scripts/build.sh
$ sudo ./scripts/install.sh
```

## TODO
* Add graphs for latency
  + Store latency details in redis
* Add method to add outage log
* Add explicit v4/v6 testing instead of relying on the default
* Move to JS refresh model instead of meta-refresh
* Make buttons into nice spans
* Make page look proper on mobile
