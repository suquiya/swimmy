swimmy
====

## Overview
Swimmy is the tool that pull meta info from url and write info to html or json.

It is a package that fetch URL Info and process it. It is for embedding external site information as card or outputting as JSON format.

suquiya is a beginner of go programming, so pull requests and issues are appropriated.
Please help suquiya......

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
go get -u github.com/suquiya/swimmy/cmd/swimmy
```

### Command
```
swimmy url <output> <flags>
```
Flags:
    + -i, --IfOutputExist: define behavitor in case that output file specified by user is already exist
    + -f, --format: user can choose html or json
    + -l, --list: If this flag is true, url is interpreted as path to file listing urls.

More details, see swimmy's help (swimmy --help).

### Package
 User can choose this as library.

Example:
```
url, ctype, content, _ := swimmy.Fetch(url)
pd := swimmy.BuildPageData(url, ctype, string(content))
jbyte, _ := pd.ToJSON()
cb := swimmy.DefaultCardBuilder
/*
user can also set like below.
cb := NewCardBuilder(cardTemplate, classNames)
*/
cb.Execute(pd, w)   //w is io.Writer

```

## To Do:

+ Implement Test
+ Add Function about favicon and images
+ improve readme
