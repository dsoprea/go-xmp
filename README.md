[![Build Status](https://travis-ci.com/dsoprea/go-xmp.svg?branch=master)](https://travis-ci.com/dsoprea/go-xmp)
[![Coverage Status](https://coveralls.io/repos/github/dsoprea/go-xmp/badge.svg?branch=master)](https://coveralls.io/github/dsoprea/go-xmp?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dsoprea/go-xmp)](https://goreportcard.com/report/github.com/dsoprea/go-xmp)
[![GoDoc](https://godoc.org/github.com/dsoprea/go-xmp?status.svg)](https://godoc.org/github.com/dsoprea/go-xmp)


# Overview

This library manages reading and writing XMP data and is written in pure Go. All
standard namespaces are supported, and values are parsed to correct types.

*Write support is incomplete and will be introduced soon.*

A simple tool has been provided that can dump the metadata or print it as a
simple JSON structure. Verbosity can be enabled to show warnings that arose
while parsing.


# Namespace Support

All standard namespaces (described in Part 1 and Part 2 of the specificiation)
are supported. Non-standard namespaces can be easily defined and registered. If
there are non-standard namespaces that you believe should be a part of this
project, post an issue to start a discussion.
