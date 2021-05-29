# Hiboot - web/cli application framework 

<p align="center">
  <a href="https://hiboot.hidevops.io">
    <img src="https://hiboot.hidevops.io/images/logo.png" alt="hiboot">
  </a>
</p>

<p align="center">
  <a href="https://travis-ci.org/hidevopsio/hiboot?branch=master">
    <img src="https://travis-ci.org/hidevopsio/hiboot.svg?branch=master" alt="Build Status"/>
  </a>
  <a class="badge-align" href="https://www.codacy.com/app/john-deng/hiboot?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=hidevopsio/hiboot&amp;utm_campaign=Badge_Grade"><img src="https://api.codacy.com/project/badge/Grade/ee8ddbf56ece4f46a6efeb216c351a0f"/></a>
  <a href="https://github.com/hidevopsio/hiboot">
    <img src="https://tokei.rs/b1/github/hidevopsio/hiboot" />
  </a>
  <a href="https://codecov.io/gh/hidevopsio/hiboot">
    <img src="https://codecov.io/gh/hidevopsio/hiboot/branch/master/graph/badge.svg" />
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
      <img src="https://img.shields.io/badge/License-Apache%202.0-green.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/hidevopsio/hiboot">
      <img src="https://goreportcard.com/badge/github.com/hidevopsio/hiboot" />
  </a>
  <a href="https://godoc.org/github.com/hidevopsio/hiboot">
      <img src="https://godoc.org/github.com/golang/gddo?status.svg" />
  </a>
</p>

## About

Hiboot is a cloud native web and cli application framework written in Go.

Hiboot is not trying to reinvent everything, it integrates the popular libraries but make them simpler, easier to use. It borrowed some of the Spring features like dependency injection, aspect oriented programming, and auto configuration. You can integrate any other libraries easily by auto configuration with dependency injection support.

If you are a Java developer, you can start coding in Go without learning curve.

## Overview

* Web MVC (Model-View-Controller).
* Auto Configuration, pre-create instance with properties configs for dependency injection.
* Dependency injection with struct tag name **\`inject:""\`**, **Constructor** func, or **Method**.

## Getting Started

* [Hiboot Documentation in English](https://hiboot.netlify.app/)
* [Hiboot 中文文档](https://hiboot.netlify.app/cn)

## Community Contributions Guide

Thank you for considering contributing to the Hiboot framework, The contribution guide can be found [here](CONTRIBUTING.md).

## License

© John Deng, 2017 ~ time.Now

Released under the [Apache License 2.0](https://github.com/hidevopsio/hiboot/blob/master/LICENSE)

[Jetbrains](https://www.jetbrains.com/?from=hiboot) supports this project with GoLand licenses. We appreciate their support for free and open source software!