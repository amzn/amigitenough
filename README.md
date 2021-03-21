# Am I Git Enough?

AmIGitEnough is a tool that verifies that users have configured git for rebase workflows. The settings AmIGitEnough
recommends will help users minimize the number of unnecessary merge commits they create and the number of conflicts they
resolve. Experienced developers point new team members to AmIGitEnough because unnecessary merge commits in shared
branches are a shared burden.

# Usage

```bash
$ amigitenough
VALIDATION SUCCEEDED
```

# What Does AmIGitEnough Check For?

AmIGitEnough checks that the user has the following configurations. If validation fails, odds are the user does not have
these settings in their $HOME/.gitconfig or wherever else they have their configuration file located.

```
[rerere]  
	enabled = true  
[pull]  
	rebase = true
```

Additionally, AmIGitEnough will verify that the user has configured their name and Email address.

# An Example

When the user does not have one of the required settings, AmIGitEnough will provide instructions to remedy the problem.

```bash
$ amigitenough
Rebase when performing "git pull" is not enabled. Correct this by running: git config --global pull.rebase true
VALIDATION FAILED
```

# Why rebase=true?

This setting ensures that `git pull` is actually performing `git pull --rebase`; pull operations rebase the user's local
changes on top of incoming commits fetched from the upstream branch.   
By default, `git pull` uses a *fetch* followed by a *merge*, and the merge will occasionally cause unnecessary
non-linear history even when the user is developing on the main branch.   
With this setting, `git pull` performs a *fetch* followed by a *rebase*.

# Why Enable Rerere?

The effect of `git pull --rebase` is that **all** missing changes starting from something called the `merge-base`
onwards get moved into your history from `upstream` before your changes. Sometimes there are conflicts to be resolved,
and this is expected. Conflicts will only need to be addressed once as long as the `merge-base` moves forward in time
after each `git pull --rebase`. When the `merge-base` gets "stuck" at a fixed point in history the user needs to resolve
the **same** conflicts after each execution `git pull --rebase.`

A scenario where the `merge-base` could get "stuck"  is when the user pulls a change from a teammate's feature branch
into their own feature branch for the purpose of continuing development against the main branch using subsequent
rebases. Developers who work solo all the time are less likely to encounter this issue of a "stuck" `merge-base`.

We include the setting in `amigitenough` because the conflicts are bewildering when they occur, and we know of no
situation where having the setting is a disadvantage.

For more information about `rerere`, see this [book section](https://www.git-scm.com/book/en/v2/Git-Tools-Rerere).
