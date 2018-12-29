# uhaul
[![Build Status](https://travis-ci.org/mitchr/uhaul.svg?branch=master)](https://travis-ci.org/mitchr/uhaul)

uhaul is an implementation of Python's [struct](https://docs.python.org/3/library/struct.html) library in go

Since go currently has no way to determine endianness (atleast without using unsafe), uhaul defaults to little-endian. The byte-order characters '>' and '!' still allow for switching to big-endian.

The Python struct library also determine type size by computing sizeof() using C interopability. This is difficult in Go, so uhaul defines the standard type sizes defined by the struct documentation, and allows these sizes to be changed by exporting the variables associated with each type's size (CHAR, SCHAR, UCHAR, etc.). These are defined in uhaul.go.
