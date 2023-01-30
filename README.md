# es - envsubst on steroids

`es` is a drop-in replacement for `envsubst` that adds support for:

* Loading environment variables from files
* Requiring all variables to be set
* Traversing a directory and substituting all files
* Printing variables and their values used in a file (or directory of files)

## Install

### From Source

```bash
make
```

### From GitHub

```bash
curl https://raw.githubusercontent.com/robertlestak/es/main/scripts/install.sh | bash
```

## Usage

```bash
Usage of es:
  -e string
        comma separated list of env files
  -i string
        input file (default "-")
  -l string
        log level (default "info")
  -o string
        output file (default "-")
  -r    require all variables to be set
  -v    print variables
  -vv
        print variables and values
```

If an argument is provided, it will be used as the input file.

## Examples

### Simple

`es` can be used as a drop-in replacement for `envsubst`:

```bash
$ echo "Hello, ${USER}!" | es
Hello, root!
$ echo 'Hello, ${USER}!' > hello.txt
$ es hello.txt
Hello, root!
```

### Directory Traversal

`es` can be used to traverse a directory and substitute all files:

```bash
$ mkdir -p templates/{a,b}
$ echo 'Hello, ${USER}!' > templates/a/hello.txt
$ echo 'Hello, ${USER}!' > templates/b/hello.txt
$ es -i templates -o output
$ cat output/a/hello.txt
Hello, root!
$ cat output/b/hello.txt
Hello, root!
```

If the output directory does not exist, it will be created.

If an output directory is not specified, a temp directory will be created and printed to stdout:

```bash
$ mkdir -p templates/{a,b}
$ echo 'Hello, ${USER}!' > templates/a/hello.txt
$ echo 'Hello, ${USER}!' > templates/b/hello.txt
$ OUTDIR=`es templates`
$ cat $OUTDIR/a/hello.txt
Hello, root!
$ cat $OUTDIR/b/hello.txt
Hello, root!
```

### Environment Files

`es` can be used to load environment variables from files:

```bash
$ echo 'Hello, ${USER} this is ${NAME}!' > hello.txt
$ echo 'USER=es' > foo.env
$ echo 'NAME=envsubst on steroids' > bar.env
$ es -i hello.txt -e foo.env,bar.env
Hello, es this is envsubst on steroids!
```

### Variables

`es` can be used to print the variables used in a file (or directory of files):

```bash
$ echo 'Hello, ${USER} this is ${NAME}!' > hello.txt
$ es -i hello.txt -v
USER
NAME
```

#### Values

`es` can be used to print the variables and their values used in a file (or directory of files):

```bash
$ echo 'Hello, ${USER} this is ${NAME}!' > hello.txt
$ es -i hello.txt -vv 
USER=root
NAME=
```

#### Require All Variables

To require all variables to be set, use the `-r` flag:

```bash
$ echo 'Hello, ${USER} this is ${NAME}!' > hello.txt
$ es -i hello.txt -r
ERROR: missing variable NAME
$ echo 'NAME=envsubst on steroids' > bar.env
$ es -i hello.txt -r -e bar.env
Hello, root this is envsubst on steroids!
```
