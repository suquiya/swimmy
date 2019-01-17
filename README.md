swimmy(WIP)
====

## Overview
Swimmy is the tool that pull meta info from url and write info to html or json.

It is a package that fetch URL Info and process it. It is for embedding external site information as card or outputting as JSON format.

## Description

This package contain cli tool and library.

Swimmy can create embed card as below:

+ First: swimmy fetch information of website specified in given URL(html document and text document)
+ Second: swimmy sanitize html contents and parse it in order to get the information of webpage
+ Third: After Parse, swimmy get PageData of the website and create Card or output information as JSON format.

## Usage

To install cli tool, use

```
go get -u github.com/suquiya/swimmy/swimmy
```

