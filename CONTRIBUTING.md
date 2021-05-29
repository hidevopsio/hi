# Contributing to Hiboot

In order to contribute a feature to Hiboot you'll need to go through the following steps:

* Create a GitHub issue to track the discussion. The issue should include information about the requirements and use cases that it is trying to address. Include a discussion of the proposed design and technical details of the implementation in the issue.
* Fork the Hiboot repository and implement the requirements.
* Submit PRs to hidevopsio/hiboot with your code changes.

## Pull requests

If you're working on an existing issue, simply respond to the issue and express interest in working on it. This helps other people know that the issue is active, and hopefully prevents duplicated efforts.

### To submit a proposed change:

* Fork the affected repository.
* Create a new branch for your changes.
* Develop the code/fix.
* Add new test cases. In the case of a bug fix, the tests should fail without your code changes. For new features try to cover as many variants as reasonably possible.
* Modify the documentation as necessary.
* Verify the entire CI process (building and testing) works.

While there may be exceptions, the general rule is that all PRs should be 100% complete - meaning they should include all test cases and documentation changes related to the change.

## Development Guide

### Git workflow
Below, we outline one of the more common Git workflows that core developers use. Other Git workflows are also valid.

### Fork the main repository

* Go to https://github.com/hidevopsio/hiboot
* Click the "Fork" button (at the top right)

### Clone your fork

The commands below require that you have $GOPATH set ($GOPATH docs). We highly recommend you put Istio's code into your GOPATH. Note: the commands below will not work if there is more than one directory in your $GOPATH.

```bash
export GITHUB_USER=your-github-username
mkdir -p $GOPATH/src/hidevops.io
cd $GOPATH/src/hidevops.io
git clone https://github.com/$GITHUB_USER/hiboot
cd hiboot
git remote add upstream 'https://github.com/hidevopsio/hiboot'
git config --global --add http.followRedirects 1
```

### Create a branch and make changes

```bash
git checkout -b my-feature
# Then make your code changes
```

### Keeping your fork in sync

```bash
git fetch upstream
git rebase upstream/master
```

Note: If you have write access to the main repositories (e.g. github.com/hidevopsio/hiboot), you should modify your Git configuration so that you can't accidentally push to upstream:
