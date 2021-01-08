# Compared with mercari/tfnotify

suzuki-shunsuke/tfnotify is compatible with mercari/tfnotify.

* Features
  * [Support `keep_duplicate_comments` to keep duplicate comments](#support-keep_duplicate_comments-to-keep-duplicate-comments)
  * [Find the configuration file recursively](#find-the-configuration-file-recursively)
  * [Complement CI and GitHub Repository owner and name from environment variables](#complement-ci-and-github-repository-owner-and-name-from-environment-variables)
  * [Support to configure label colors](#support-to-configure-label-colors)
  * Support template functions [sprig](http://masterminds.github.io/sprig/)
  * [Support to pass variables by -var option](#support-to-pass-variables-by--var-option)
  * [Don't recreate labels](#dont-recreate-labels)
  * [--version option and version command](#--version-option-and-version-command)
* Fixes
  * [Post a comment even if it failed to update GitHub Pull Request labels](#fix-post-a-comment-even-if-it-failed-to-update-github-pull-request-labels)
* Others
  * refactoring
  * update urfave/cli to v2

## Support `keep_duplicate_comments` to keep duplicate comments

tfnotify deletes duplicate comments at GitHub and GitLab.
But by setting `keep_duplicate_comments: true`, tfnotify doesn't remove them.

```yaml
notifier:
  github:
    token: $GITHUB_TOKEN
keep_duplicate_comments: true
# ...
```

## Find the configuration file recursively

[#19](https://github.com/suzuki-shunsuke/tfnotify/pull/19)

tfnotify searches the configuration file from the current directory to the root directory recursively.

## Complement CI and GitHub Repository owner and name from environment variables

Supported platform

* CI
  * CircleCI
  * CodeBuild
  * GitHub Actions
  * Drone
* Notifier
  * GitHub

The configuration of CI and GitHub Repository owner and name is complemented by CI builtin environment variables.
[suzuki-shunsuke/go-ci-env](https://github.com/suzuki-shunsuke/go-ci-env) is used to complement them.

AS IS

```yaml
ci: circleci
notifier:
  github:
    token: $GITHUB_TOKEN
    repository:
      owner: suzuki-shunsuke
      name: tfcmt
```

We can omit `ci` and `repository`.

```yaml
notifier:
  github:
    token: $GITHUB_TOKEN
```

## Support to configure label colors

[98547135a6d37b11b641feb399eec17721fe0bc0](https://github.com/suzuki-shunsuke/tfnotify/commit/98547135a6d37b11b641feb399eec17721fe0bc0)
[49ea5c3a8c01e53cac6d3b529bd5d9907c41e9d3](https://github.com/suzuki-shunsuke/tfnotify/commit/49ea5c3a8c01e53cac6d3b529bd5d9907c41e9d3)

tfnotify supports to configure label colors.
So we don't have to configure label colors manually.
This feature is useful especially for Monorepo.

## Support to pass variables by -var option

[#29](https://github.com/suzuki-shunsuke/tfnotify/pull/29)

tfnotify supports to pass variables to template by `-var <name>:<value>` options.
We can access the variable in the template by `{{.Vars.<variable name>}}`.

## Don't recreate labels

[#32](https://github.com/suzuki-shunsuke/tfnotify/pull/32)

If the label which tfnotify set is already set to a pull request, mercari/tfnotify removes the label from the pull request and re-adds the same label to the pull request.
This is meaningless.

So suzuki-shunsuke/tfnotify doesn't recreate a label.

## --version option and version command

[#4](https://github.com/suzuki-shunsuke/tfnotify/pull/4)
[#44](https://github.com/suzuki-shunsuke/tfnotify/pull/44)

AS IS

```
$ tfnotify --version
tfnotify version unset
```

TO BE

```
$ tfnotify --version
tfnotify version 1.3.3

$ tfnotify version
tfnotify version 1.3.3
```

## Fix: Post a comment even if it failed to update GitHub Pull Request labels

[#56](https://github.com/suzuki-shunsuke/tfnotify/pull/56)

mercari/tfnotify doesn't post a comment when it failed to update GitHub Pull Request labels.
For example, when the label length is too long, tfnotify failed to add the label and the comment isn't posted.

On the other hand, suzuki-shunsuke/tfnotify outputs the error log but the process continues even if it failed to update labels.
