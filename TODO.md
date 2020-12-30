- make errors return html by adding an error template
- style it up
- write tests
- add reset function that loads data from json file
- open issue at https://github.com/mattn/go-sqlite3 about static linking problem:
```
/usr/bin/ld: /tmp/go-link-066404605/000010.o: in function `unixDlOpen':
/go/src/github.com/mattn/go-sqlite3/sqlite3-binding.c:40175: warning: Using 'dlopen' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
```
