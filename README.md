# und

[![LICENSE](https://img.shields.io/github/license/jackysp/und.svg)](https://github.com/jackysp/und/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/jackysp/und?status.svg)](https://godoc.org/github.com/jackysp/und)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackysp/und)](https://goreportcard.com/report/github.com/jackysp/und)
[![GitHub release](https://img.shields.io/github/release/jackysp/und.svg)](https://github.com/jackysp/und/releases/latest)
[![GitHub release date](https://img.shields.io/github/release-date/jackysp/und.svg)](https://github.com/jackysp/und/releases)

Update [namesilo](https://www.namesilo.com/) DNS record dynamically.

## Build from source

Require Go 1.13+.

1. `git clone https://github.com/jackysp/und.git`
1. `cd und`
1. `make`

## Installation

### Download the binary (recommanded)

https://github.com/jackysp/und/releases

### Install from source

`go get -u github.com/jackysp/und`

## Usage

./und -key={your key} -domain={your domain} -host={your hostname} -interval={update interval}
