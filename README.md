# sup

Command sup is another work tracking tool.

## install

```zsh
% go get github.com/martindrlik/sup
```

## usage

```text
% /path/to/sup
start reading book
start watching netflix
ps
  id      took name
   0   1h30m0s reading book
   1   1h20m0s watching netflix
       2h50m0s
start free time
```

## help

- `start name` starts new or resumes existing task `name`
- `ps pattern` prints tasks filtered by regexp `pattern`
- `fixname i name` sets `i`-th task name to `name`
