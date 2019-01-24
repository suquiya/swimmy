swimmy(WIP)
====

## Overview
Swimmy is the tool that pull meta info from url and write info to html or json.

It is a package that fetch URL Info and process it. It is for embedding external site information as card or outputting as JSON format.

suquiya is a beginner of go programming, so pull requests and issues are appropriated.
To be honest, please help suquiya...

## Description

This package contain cli tool and library.

Swimmy can create embed card as below:

+ First: swimmy fetch information of website specified in given URL(html document and text document)
+ Second: swimmy sanitize html contents and parse it in order to get the information of webpage
+ Third: After Parse, swimmy get PageData of the website and create Card or output information as JSON format.

Of cource, Swimmy is not only cli tool but also library. Using as library(package), you can use function of parts of swimmy.

## Usage

To install,

```
go get -u github.com/suquiya/swimmy/swimmy
```

## To Do:

+ Implement Test
+ Add Function about favicon and images
