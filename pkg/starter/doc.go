// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package starter provides quick starters for hiboot cli/web application.

Auto-configured Starter

Hiboot auto-configuration attempts to automatically configure your Hiboot application based on the pkg dependencies that you have added.
For example, if bolt is imported in you main.go, and you have not manually configured any database connection,
then Hiboot auto-configures an database bolt for any service to inject.

You need to opt-in to auto-configuration by embedding app.Configuration in your configuration and
calling the app.Register() function inside the init() function of your configuration pkg.

For more details, see https://godoc.org/github.com/hidevopsio/hiboot/pkg/starter

Creating Your Own Starter

A full Hiboot starter for a library may contain the following structs:
	autoconfigure - object that handle the auto-configuration code.
	properties - object that contains properties which will be injected configurable default values or user specified values
If you work in a company that develops shared go packages, or if you work on an open-source or commercial project,
you might want to develop your own auto-configured starter. starter can be implemented in external packages and
can be imported by any go applications.

Understanding Auto-configured Starter

Under the hood, auto-configuration is implemented with standard struct. Additional embedded field app.Configuration.
AutoConfiguration used to constrain when the auto-configuration should apply. Usually, auto-configuration struct use
`after:"fooConfiguration"` or `missing:"fooConfiguration"` tags. This ensures that auto-configuration applies only
when relevant configuration are found and when you have not declared your own configuration.

*/
package starter
