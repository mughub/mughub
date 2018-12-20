# mHub

A modular suite of packages, which when composed together provide a highly
customizable Git service.

# Getting Started

- [Design](#design)
    * [Modules](#modules)
    * [Git](#git)
    * [Database](#database)
    * [User Interface](#user-interface)
- [Self Hosting](#self-hosting)
    * [Mono](#mono)
    * [Distributed](#distributed)

## Design
The following section describes the reasoning and high-level details of μgHub's
design.

### Modules
The core design concept of modulation for GoHub is a direct result of the
main goals for μgHub:

* Support open source development with large contribution base
* Be highly configurable, customizable, and extendable
* Use idiomatic Go and adhere to [SOLID](https://www.youtube.com/watch?v=zzAdEt3xZ1M) Go design
* Integrate into modern service based architectures

These goals combined with some practical thinking leads to an overall modular
design which emphasizes singular state.

### Git
The git module (`package bare`) handles defining the bare bones types that
drive μgHub at its' core. The main purpose of this package is provide a
'core' service for μgHub, in other words, it provides the functionality of
a Git server which is at the heart of any Git service.

### Database
The database module (`package db`) handles data management for μgHub. All
requests for data are defined in a GraphQL [API]().

### User Interface
The user interface module (`package ui`) handles presenting a UI for μgHub.

## Self Hosting
The following section describes the primary ways GoHub is intended to be
deployed.

### Mono
The mono deployment refers to simply deploying μgHub using the `gohub` cli
tool. Deploying with the cli is geared toward the same principles that drive
the design of μgHub itself. Configuration:
* All desired endpoints must be specified
* Supported databases: Dgraph, GraphQLite, Postgress and MySQL

### Distributed
A distributed deployment of μgHub refers to deploying each μgHub module as
its' own microservice.