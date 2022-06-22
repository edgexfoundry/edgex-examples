# Camera Management Web UI

The Camera Management Web UI is an Angular interface that provides features such as:

* PTZ
* EVAM Pipeline Start/Stop/Status
* API Logs
* EVAM Event Logs


## How do I build this thing?
This project is intended to be built and run using only `make`, `docker`, and `docker-compose` without the need
for installing `npm` or `nodejs`.

To skip the technical details, jump to **[Initial Setup](#Initial-Setup)**

If you want to know how the heck this thing works and what it is doing to your machine, 
read on to **[Build Ideology](#Build-Ideology)**

## Build Ideology
A lot of work has gone into making this friendlier to build and run. `Docker` and `make` are used to
allow the developer to work without ever installing `nodejs` or `angular-cli`.

A base nodejs docker image is built using your local user and group. When this image
is run, this entire repo is mounted under `/app` inside the container. This means whatever you do inside it
is also done locally as well. However, it also means that whatever you do locally is automatically
available to the docker side of things. This allows the developers to work completely
dockerized, completely native, or somewhere in-between.

Two docker files are provided, a dev one and a production one. By default, the dev one will be used. View the
[Production](#Production) section to see how to build and run the production version.

`make` commands are used to abstract away more complicated docker command line. 

`make debug` will launch an interactive build container to run any commands you would like (`npm`, `ng`, etc).

`make ng="..."` and `make npm="..."` allow you to run ng and npm commands within the confines of the
docker build container.

## Initial Setup

### Create development docker image
```shell
make image
```

### Install npm packages (node_modules)
```shell
make install
```

### Build and serve dev code
```shell
# Foreground
make up
```
```shell
# Background
make deploy

# Tail the logs
make tail
```

### View Website

The website is available at [http://localhost:4200](http://localhost:4200).

Any changes to the code will be hot-reloaded by the angular dev server (except in production mode).

### Stop / Remove Containers
In order to stop docker compose which is in the foreground, use `Ctrl-C`.
If you used `make deploy` to launch it in the background, use these:


```shell
# Stop running containers
make stop
```
OR
```shell
# Stop and remove running containers
make down
```

### Other useful make commands
```shell
# Build and run unit tests
make test

# Open the unit test code coverage in your browser
make view-coverage

# Run linter
make lint

# Check for security issues with node packages
make audit

# Attempt to fix security issues with packages
make audit-fix

# Upgrade angular to the latest stable version
make update-angular

# Clean build cache, coverage info and dist files
make clean

# Clear out all node modules and build artifacts
make nuke

# Run an npm install
make install

# Run an npm install --force
make force-install

# Open an interactive build environment
# Desktop x11 environment is mapped and programs such as google-chrome can be ran in GUI mode
make debug

# Tail docker-compose logs (`n` is optional amount of lines to tail first)
make tail [n=XXX]

# Prints docker-compose logs and exits (`n` is optional amount of lines to print)
make logs [n=XXX] 

# Generate documentation
make docs

# Open documentation in web-browser
make view-docs

# Open web-browser to root page of this app
make view

# [Broken - do not use] Run end-to-end tests
make e2e

# Run ng serve in debug (hot-reload) mode
make serve

```

### Advanced Make Commands
Open an interactive build environment to run any `npm`, `ng` or other command within the confines
of the volume mounted dev container
```shell
make debug
```

To run `npm` commands within the confines of the volume mounted dev container, set the `npm` variable
to your `npm` argument list like so: 
```shell
# This will run "npm audit --fix"
make npm="audit --fix"

# This will run "npm install [package-name] --save"
make npm="install moment --save"

```

To run `angular-cli` (`ng`) commands within the confines of the volume mounted dev container, set the `ng` variable
to your `ng` argument list like so: 
```shell
# This will run "ng generate component Test"
make ng="generate component Test"
```

To run any arbitrary command within the confines of the volume mounted dev container, set the `exec`
variable to your command and argument list like so:
```shell
# This will run "node --version"
make exec="node --version"
```

## Code Generation

Run `make gen="component component-name"` to generate a new component.

You can generate the following schematics:
- app-shell
- application
- class
- component
- directive
- enum
- guard
- interceptor
- interface
- library
- module
- pipe
- resolver
- service
- service-worker
- web-worker

_**Note:** Do not include the suffix (`Component`, `Service`, etc). They are added for you_

## Testing

Run `make test` to execute the unit tests via [Karma](https://karma-runner.github.io).
You can view the Karma tests at the address [http://127.0.0.1:9876](http://127.0.0.1:9876). Changes will be hot-reloaded
and tests will be re-run.

Run `make e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

Run `make view-coverage` to view static code coverage HTML output. Alternatively, open [coverage/report-html/index.html](coverage/report-html/index.html) with your browser.

**Note:** In order for the code coverage to be accurate, you must first _refresh_ the Karma unit test service at [http://127.0.0.1:9876](http://127.0.0.1:9876). 
It will then produce the HTML code coverage output for you to view in your browser as an HTML file.

## Native Tooling
In order to setup a native tooling environment without docker (for instance, use with an IDE), follow these steps to
download and install all dependencies.

### Install NodeJS LTS
- _[Recommended]_ Install from [binary distributions](https://github.com/nodesource/distributions/blob/master/README.md#installation-instructions)
- Install via [snap](https://snapcraft.io/node)
  - _Note: Installations via snap will sometimes cause issues with installing global packages such as `angular-cli`_

### Install Angular CLI
```shell
npm install -g @angular/cli
```

## Production
The production mode generates production built angular static html/js files able to be served via `nginx` or similar web server.

### Build production files
```shell
make prod
```
