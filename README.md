## about

This lib is inspired by [s6 fifodirs](http://skarnet.org/software/s6/fifodir.html).
I wrote this because I don't want to have to link in the whole s6 lib, I just want the fifodirs part for another project.
The C code is blatantly copy/pasted from [upstream](https://github.com/skarnet/s6).

This is pretty experimental stuff, it may burn your house down.

## usage
```
:; go get github.com/jdef/gotrig

## assuming that your $GOPATH/bin is in your $PATH
:; mkdir -p /tmp/abc
:; listen /tmp/abc/ &
:; notify /tmp/abc/ billy
:; kill %
```
