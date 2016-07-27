## hkjl_status
Small tool which polls all Hackenkunjeleren services and generates a html
page displaying the status.

## installation
The procedure below works for Debian Jessie. If you have any other OS, you're
on your own

```bash
sudo apt-get install git golang binutils
git clone https://github.com/r3boot/hkjl_status
cd hkjl_status
sudo ./scripts/install.sh
```

## TODO
* Add graphs for latency
* Add method to add outage log
