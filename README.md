# GH Notification Purge

> [![release](https://github.com/mtwig/gh-notification-purge/actions/workflows/release.yml/badge.svg)](https://github.com/mtwig/gh-notification-purge/actions/workflows/release.yml)

The purpose of this extension is to (you guessed it!) mark Github notifications as done. 
Specifically, this to cleanup notifications about pull requests which have _already been closed_!

There is a limitation with the Github API, which prevents notifications from being marked as done. It is currently only possible to mark notifications as done.
To work around the limitation, create a Githb notification filter for `is:read`, and bulk mark as done.

## 🚀 Installation

```bash
gh extension install mtwig/gh-notification-purge
```

## Usage
```bash
gh notification-purge
```


