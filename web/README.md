# Purpose

This directory holds resources used for the web portion of an application

# Organization

> [!WARNING]  
> If any pathing is updated, make sure to update the paths across scripts in the `package.json` and `Taskfile.yml` files

## Components

This directory holds reusable component templates that can be composed across different page layouts

## Layouts

This directory is responsible for different layouts across web pages of your site

## Libs

This directory serves as an entrypoint to be used for storage of any JS/TS libraries needed for the project to run

Currently it is being used to hold the following:

- Datastar as a vendor dependency from the `npm` registry - [NOTE](./libs/datastar/README.md)
- Web Components powered by Lit-html

## Pages

This directory is responsible for the different pages accessible through your site

## Styles

This directory is responsible for styling, it is currently setup to use [TailwindCSS](https://tailwindcss.com/)
