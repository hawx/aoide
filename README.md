# aoide

A music library.


## The plan

``` bash
$ cat ~/.config/aoide.toml
musicDir = "/home/.../Music"
playlistDir = "/home/.../.config/mpd/playlists"
dbFile = "/home/.../.cache/aoide.db"
$ aoide index
...[reads through musicDir/**/* adding to dbFile]...
$ aoide organise --dry-run
...[prints changes that will be made]...
$ aoide organise
...[does changes]...
$ aoide autotag
...
$ aoide get-art
...
$ [others?]
...
```
