# Graphic Raid

![CI](https://github.com/dannyh79/graphic-raid/actions/workflows/test.yml/badge.svg)

## Status

### Classroom

- [x] Deliverable
- [x] Required spec
    > See test suites under /classroom
- [ ] Bonus: Student giving wrong answer
- [ ] Bonus: Spawn quizzes as independent processes

### Quorum

- [ ] Deliverable
- [ ] Required spec
    > See test suites under /quorum for finished specs

    - [x] Leader election
    - [x] KeepAlive mechanism
    - [x] Pluggable quorum mechanism
    - [ ] Eviction behavior
    - [ ] Re-elect leader behavior
    - [ ] Extract domain-specific (quorum) behaviors from usecase
- [ ] Bonus: Different quorum mechanism

### RAID

- [ ] Deliverable

## Getting Started

```shell
# under project root directory

$ go version
# go version go1.22.1

$ go build -C classroom -o ../bin
$ .bin/classroom
```


## Development

```shell
$ go run -race classroom/main.go
$ go run -race quorum/main.go
```

### Testing

```shell
# under project root directory

$ ginkgo ./...
```

## Backlog

- [ ] Extract actor behaviors as a general-purpose package

### Classroom

- [ ] Use pub-sub pattern in inter-process communications
- [ ] Leverage channels with defined message struct to reduce channels needed in goroutines' params

### Quorum

_N/A._

### RAID

_N/A._
