# ftransfer
Very simple Go program to transfer files between computers in a local LAN. Made to explore Golang a bit

# Compiling

You will need Golang installed in you machine. After that, in the same directory of this software, open a terminal and:

```shell
go build .
```

# Usage

You can use this as command line or by double clicking it. Using with command line you can pass some arguments to the command. Use `ftransfer --help` for a list of parameters.

If you use it without parameters (or just double clicking it) it have a very simple and self explanatory menu.

# Tips

If you are a receiver and your sender haven't executed the program yet, you will see an error. The sender should run the software first and then wait for connections.
