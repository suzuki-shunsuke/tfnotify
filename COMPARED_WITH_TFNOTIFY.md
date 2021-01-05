# Compared with mercari/tfnotify

## Features

* don't recreate labels
* support to configure label colors
* support template functions [sprig](http://masterminds.github.io/sprig/)
* support to pass variables by -var option
* support to find the configuration file recursively
* support --version option

### don't recreate labels

If the label which tfnotify set is already set to a pull request, tfnotify removes the label from the pull request and re-adds the same label to the pull request.
This is meaningless.

So tfcmt doesn't recreate a label.

### support to configure label colors

tfcmt supports to configure label colors.
So we don't have to configure label colors manually.
This feature is useful especially for Monorepo.

### support to pass variables by -var option

tfcmt supports to pass variables to template by `-var <name>:<value>` options.
We can access the variable in the template by `{{.Vars.<variable name>}}`.

### support to find the configuration file recursively

tfcmt searches the configuration file from the current directory to the root directory recursively.

### support --version option

AS IS

```
$ tfnotify --version
tfnotify version unset
```

TO BE

```
$ tfnotify --version
tfcmt version 1.3.3
```

## Others

* refactoring
* update urfave/cli to v2
