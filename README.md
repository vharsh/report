# Report
Uses steroids to create github-issues in crazy-times.

## Functionality
- Can dump the output of some program as a gist (DONE)
- Can create an issue (DONE)


## Sample Usage
```bash
user@hostname :~/code/go/src/github.com/vharsh/report $ ./report --repo="vharsh/report" --title="This is a TEST" --desc="This TEST was successful" systemctl status docker
Issue# https://github.com/vharsh/report/issues/14 created
[systemctl status docker]  will be executed
Add some description: The gist thing works too
https://gist.github.com/316bb1a9d9646451a5ce5b15b501bbb9
```
The repo, issue title & description can also be entered interactively if not provided for in the CLI-arguments.

