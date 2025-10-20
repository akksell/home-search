# ADR - Use nix to create reproducible dev environments
Author: Axel Ramone
Creation Date: 2025-10-17

## Problem
This project requires multiple languages, each with their own tooling. If another developer wants to come and work on this project, there will be a significant barrier and productivity loss as the developer will need to spend time installing the required tooling and environment configurations.

## Analysis
Nix is a powerful tool that allows for reproducible environments. It allows for developer environments to be created by specifying the exact packages a developer wants to use in their shell.

There is a steep learning learning curve to Nix and their documentation is not the most the most helpful. However, I can reference existing usage of this tool in other projects that I've worked on.

## Solution
Nix will be configured and incorporated into this project. It's definitely overkill considering the likelyhood of another developer working on this project is low but it will allow me to at least use consistent versions of the languages used in this project across each service. I also won't need to install different version managers just for a specific language.