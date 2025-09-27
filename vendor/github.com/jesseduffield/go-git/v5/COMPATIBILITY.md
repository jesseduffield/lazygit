# Supported Features

Here is a non-comprehensive table of git commands and features and their
compatibility status with go-git.

## Getting and creating repositories

| Feature | Sub-feature                                                                                                        | Status | Notes | Examples                                                                                                                                                                                                            |
| ------- | ------------------------------------------------------------------------------------------------------------------ | ------ | ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `init`  |                                                                                                                    | ✅     |       |                                                                                                                                                                                                                     |
| `init`  | `--bare`                                                                                                           | ✅     |       |                                                                                                                                                                                                                     |
| `init`  | `--template` <br/> `--separate-git-dir` <br/> `--shared`                                                           | ❌     |       |                                                                                                                                                                                                                     |
| `clone` |                                                                                                                    | ✅     |       | - [PlainClone](_examples/clone/main.go)                                                                                                                                                                             |
| `clone` | Authentication: <br/> - none <br/> - access token <br/> - username + password <br/> - ssh                          | ✅     |       | - [clone ssh (private_key)](_examples/clone/auth/ssh/private_key/main.go) <br/> - [clone ssh (ssh_agent)](_examples/clone/auth/ssh/ssh_agent/main.go) <br/> - [clone access token](_examples/clone/auth/basic/access_token/main.go) <br/> - [clone user + password](_examples/clone/auth/basic/username_password/main.go) |
| `clone` | `--progress` <br/> `--single-branch` <br/> `--depth` <br/> `--origin` <br/> `--recurse-submodules` <br/>`--shared` | ✅     |       | - [recurse submodules](_examples/clone/main.go) <br/> - [progress](_examples/progress/main.go)                                                                                                                      |

## Basic snapshotting

| Feature  | Sub-feature | Status | Notes                                                    | Examples                             |
| -------- | ----------- | ------ | -------------------------------------------------------- | ------------------------------------ |
| `add`    |             | ✅     | Plain add is supported. Any other flags aren't supported |                                      |
| `status` |             | ✅     |                                                          |                                      |
| `commit` |             | ✅     |                                                          | - [commit](_examples/commit/main.go) |
| `reset`  |             | ✅     |                                                          |                                      |
| `rm`     |             | ✅     |                                                          |                                      |
| `mv`     |             | ✅     |                                                          |                                      |

## Branching and merging

| Feature     | Sub-feature | Status       | Notes                                   | Examples                                                                                        |
| ----------- | ----------- | ------------ | --------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `branch`    |             | ✅           |                                         | - [branch](_examples/branch/main.go)                                                            |
| `checkout`  |             | ✅           | Basic usages of checkout are supported. | - [checkout](_examples/checkout/main.go)                                                        |
| `merge`     |             | ⚠️ (partial) | Fast-forward only                       |                                                                                                 |
| `mergetool` |             | ❌           |                                         |                                                                                                 |
| `stash`     |             | ❌           |                                         |                                                                                                 |
| `sparse-checkout`     |             | ✅           |                                         | - [sparse-checkout](_examples/sparse-checkout/main.go)                                                                                               |
| `tag`       |             | ✅           |                                         | - [tag](_examples/tag/main.go) <br/> - [tag create and push](_examples/tag-create-push/main.go) |

## Sharing and updating projects

| Feature     | Sub-feature | Status | Notes                                                                   | Examples                                   |
| ----------- | ----------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------ |
| `fetch`     |             | ✅     |                                                                         |                                            |
| `pull`      |             | ✅     | Only supports merges where the merge can be resolved as a fast-forward. | - [pull](_examples/pull/main.go)           |
| `push`      |             | ✅     |                                                                         | - [push](_examples/push/main.go)           |
| `remote`    |             | ✅     |                                                                         | - [remotes](_examples/remotes/main.go)     |
| `submodule` |             | ✅     |                                                                         | - [submodule](_examples/submodule/main.go) |
| `submodule` | deinit      | ❌     |                                                                         |                                            |

## Inspection and comparison

| Feature    | Sub-feature | Status    | Notes | Examples                       |
| ---------- | ----------- | --------- | ----- | ------------------------------ |
| `show`     |             | ✅        |       |                                |
| `log`      |             | ✅        |       | - [log](_examples/log/main.go) |
| `shortlog` |             | (see log) |       |                                |
| `describe` |             | ❌        |       |                                |

## Patching

| Feature       | Sub-feature | Status | Notes                                                | Examples |
| ------------- | ----------- | ------ | ---------------------------------------------------- | -------- |
| `apply`       |             | ❌     |                                                      |          |
| `cherry-pick` |             | ❌     |                                                      |          |
| `diff`        |             | ✅     | Patch object with UnifiedDiff output representation. |          |
| `rebase`      |             | ❌     |                                                      |          |
| `revert`      |             | ❌     |                                                      |          |

## Debugging

| Feature  | Sub-feature | Status | Notes | Examples                           |
| -------- | ----------- | ------ | ----- | ---------------------------------- |
| `bisect` |             | ❌     |       |                                    |
| `blame`  |             | ✅     |       | - [blame](_examples/blame/main.go) |
| `grep`   |             | ✅     |       |                                    |

## Email

| Feature        | Sub-feature | Status | Notes | Examples |
| -------------- | ----------- | ------ | ----- | -------- |
| `am`           |             | ❌     |       |          |
| `apply`        |             | ❌     |       |          |
| `format-patch` |             | ❌     |       |          |
| `send-email`   |             | ❌     |       |          |
| `request-pull` |             | ❌     |       |          |

## External systems

| Feature       | Sub-feature | Status | Notes | Examples |
| ------------- | ----------- | ------ | ----- | -------- |
| `svn`         |             | ❌     |       |          |
| `fast-import` |             | ❌     |       |          |
| `lfs`         |             | ❌     |       |          |

## Administration

| Feature         | Sub-feature | Status | Notes | Examples |
| --------------- | ----------- | ------ | ----- | -------- |
| `clean`         |             | ✅     |       |          |
| `gc`            |             | ❌     |       |          |
| `fsck`          |             | ❌     |       |          |
| `reflog`        |             | ❌     |       |          |
| `filter-branch` |             | ❌     |       |          |
| `instaweb`      |             | ❌     |       |          |
| `archive`       |             | ❌     |       |          |
| `bundle`        |             | ❌     |       |          |
| `prune`         |             | ❌     |       |          |
| `repack`        |             | ❌     |       |          |

## Server admin

| Feature              | Sub-feature | Status | Notes | Examples                                  |
| -------------------- | ----------- | ------ | ----- | ----------------------------------------- |
| `daemon`             |             | ❌     |       |                                           |
| `update-server-info` |             | ✅     |       | [cli](./cli/go-git/update_server_info.go) |

## Advanced

| Feature    | Sub-feature | Status      | Notes | Examples |
| ---------- | ----------- | ----------- | ----- | -------- |
| `notes`    |             | ❌          |       |          |
| `replace`  |             | ❌          |       |          |
| `worktree` |             | ❌          |       |          |
| `annotate` |             | (see blame) |       |          |

## GPG

| Feature             | Sub-feature | Status | Notes | Examples |
| ------------------- | ----------- | ------ | ----- | -------- |
| `git-verify-commit` |             | ✅     |       |          |
| `git-verify-tag`    |             | ✅     |       |          |

## Plumbing commands

| Feature         | Sub-feature                           | Status       | Notes                                               | Examples                                     |
| --------------- | ------------------------------------- | ------------ | --------------------------------------------------- | -------------------------------------------- |
| `cat-file`      |                                       | ✅           |                                                     |                                              |
| `check-ignore`  |                                       | ❌           |                                                     |                                              |
| `commit-tree`   |                                       | ❌           |                                                     |                                              |
| `count-objects` |                                       | ❌           |                                                     |                                              |
| `diff-index`    |                                       | ❌           |                                                     |                                              |
| `for-each-ref`  |                                       | ✅           |                                                     |                                              |
| `hash-object`   |                                       | ✅           |                                                     |                                              |
| `ls-files`      |                                       | ✅           |                                                     |                                              |
| `ls-remote`     |                                       | ✅           |                                                     | - [ls-remote](_examples/ls-remote/main.go)   |
| `merge-base`    | `--independent` <br/> `--is-ancestor` | ⚠️ (partial) | Calculates the merge-base only between two commits. | - [merge-base](_examples/merge_base/main.go) |
| `merge-base`    | `--fork-point` <br/> `--octopus`      | ❌           |                                                     |                                              |
| `read-tree`     |                                       | ❌           |                                                     |                                              |
| `rev-list`      |                                       | ✅           |                                                     |                                              |
| `rev-parse`     |                                       | ❌           |                                                     |                                              |
| `show-ref`      |                                       | ✅           |                                                     |                                              |
| `symbolic-ref`  |                                       | ✅           |                                                     |                                              |
| `update-index`  |                                       | ❌           |                                                     |                                              |
| `update-ref`    |                                       | ❌           |                                                     |                                              |
| `verify-pack`   |                                       | ❌           |                                                     |                                              |
| `write-tree`    |                                       | ❌           |                                                     |                                              |

## Indexes and Git Protocols

| Feature              | Version                                                                         | Status | Notes |
| -------------------- | ------------------------------------------------------------------------------- | ------ | ----- |
| index                | [v1](https://github.com/git/git/blob/master/Documentation/gitformat-index.txt)  | ❌     |       |
| index                | [v2](https://github.com/git/git/blob/master/Documentation/gitformat-index.txt)  | ✅     |       |
| index                | [v3](https://github.com/git/git/blob/master/Documentation/gitformat-index.txt)  | ❌     |       |
| pack-protocol        | [v1](https://github.com/git/git/blob/master/Documentation/gitprotocol-pack.txt) | ✅     |       |
| pack-protocol        | [v2](https://github.com/git/git/blob/master/Documentation/gitprotocol-v2.txt)   | ❌     |       |
| multi-pack-index     | [v1](https://github.com/git/git/blob/master/Documentation/gitformat-pack.txt)   | ❌     |       |
| pack-\*.rev files    | [v1](https://github.com/git/git/blob/master/Documentation/gitformat-pack.txt)   | ❌     |       |
| pack-\*.mtimes files | [v1](https://github.com/git/git/blob/master/Documentation/gitformat-pack.txt)   | ❌     |       |
| cruft packs          |                                                                                 | ❌     |       |

## Capabilities

| Feature                        | Status       | Notes |
| ------------------------------ | ------------ | ----- |
| `multi_ack`                    | ❌           |       |
| `multi_ack_detailed`           | ❌           |       |
| `no-done`                      | ❌           |       |
| `thin-pack`                    | ❌           |       |
| `side-band`                    | ⚠️ (partial) |       |
| `side-band-64k`                | ⚠️ (partial) |       |
| `ofs-delta`                    | ✅           |       |
| `agent`                        | ✅           |       |
| `object-format`                | ❌           |       |
| `symref`                       | ✅           |       |
| `shallow`                      | ✅           |       |
| `deepen-since`                 | ✅           |       |
| `deepen-not`                   | ❌           |       |
| `deepen-relative`              | ❌           |       |
| `no-progress`                  | ✅           |       |
| `include-tag`                  | ✅           |       |
| `report-status`                | ✅           |       |
| `report-status-v2`             | ❌           |       |
| `delete-refs`                  | ✅           |       |
| `quiet`                        | ❌           |       |
| `atomic`                       | ✅           |       |
| `push-options`                 | ✅           |       |
| `allow-tip-sha1-in-want`       | ✅           |       |
| `allow-reachable-sha1-in-want` | ❌           |       |
| `push-cert=<nonce>`            | ❌           |       |
| `filter`                       | ❌           |       |
| `session-id=<session id>`      | ❌           |       |

## Transport Schemes

| Scheme               | Status       | Notes                                                                  | Examples                                       |
| -------------------- | ------------ | ---------------------------------------------------------------------- | ---------------------------------------------- |
| `http(s)://` (dumb)  | ❌           |                                                                        |                                                |
| `http(s)://` (smart) | ✅           |                                                                        |                                                |
| `git://`             | ✅           |                                                                        |                                                |
| `ssh://`             | ✅           |                                                                        |                                                |
| `file://`            | ⚠️ (partial) | Warning: this is not pure Golang. This shells out to the `git` binary. |                                                |
| Custom               | ✅           | All existing schemes can be replaced by custom implementations.        | - [custom_http](_examples/custom_http/main.go) |

## SHA256

| Feature  | Sub-feature | Status | Notes                              | Examples                             |
| -------- | ----------- | ------ | ---------------------------------- | ------------------------------------ |
| `init`   |             | ✅     | Requires building with tag sha256. | - [init](_examples/sha256/main.go)   |
| `commit` |             | ✅     | Requires building with tag sha256. | - [commit](_examples/sha256/main.go) |
| `pull`   |             | ❌     |                                    |                                      |
| `fetch`  |             | ❌     |                                    |                                      |
| `push`   |             | ❌     |                                    |                                      |

## Other features

| Feature         | Sub-feature                 | Status | Notes                                          | Examples |
| --------------- | --------------------------- | ------ | ---------------------------------------------- | -------- |
| `config`        | `--local`                   | ✅     | Read and write per-repository (`.git/config`). |          |
| `config`        | `--global` <br/> `--system` | ✅     | Read-only.                                     |          |
| `gitignore`     |                             | ✅     |                                                |          |
| `gitattributes` |                             | ✅     |                                                |          |
| `git-worktree`  |                             | ❌     | Multiple worktrees are not supported.          |          |
