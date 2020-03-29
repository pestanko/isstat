# IS.MUNI Notepads Statistics

Tool to fetch and parse the points from the IS.MUNI notepads.

## Usage

To install the ``isstat`` you can use the `go get`
```bash
go get github.com/pestanko/isstat
```

To get available options and commands you can use the ``--help``.

```bash
isstat --help
```

## Configuration

Configuration file ``isstat-config.yml`` can by located either in current directory or `$HOME/.config/isstat/`.

To generate default configuration that can be edited just use the
 
```bash
isstat config > isstat-config.yml
```