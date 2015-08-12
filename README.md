# *xliffer* - a tool to work with XLIFF files

XLIFF is the XML Localisation Interchange File Format, standardized in
http://docs.oasis-open.org/xliff/xliff-core/xliff-core.html 

*xliffer* intends to make working with XLIFF files easier, especially
converting translations from and to XLIFF.

## Usage

Imagine an UI with a fair amount of UI-items. One approach to translate these
is to assign an unique ID (sometimes called "Message ID") to a "thing" and then
map a translation onto it. Some people find it easier to create that initial
mapping in a spreadsheet program such as Excel:

	$> xliffer from-xlsx -in base-tr.xlsx > base-tr.xlf

This base file will be transmitted to a translation service and comes back as
a number of translated files. To create a nice and shiney key-value file
JSON-format:

	$> xliffer to-json -in app-tr-jp.xlf > app-tr-jp.json

That file can then be used with http://formatjs.io/


### Detailed Usage

    xliffer converts to and from XLIFF files

    Usage: /home/mg/work/xliffer/xliffer [-h] <converter> [cflags]

    Available converters:

     from-xlsx      - Converts an Excel sheet to XLIFF
     to-json        - Converts XLIFF to JSON (key,value)

    Use <converter> -h to get the flags specific for the relevant converter

    from-xlsx:

      -in="": infile
      -key-column=3: column holding the key / msgid
      -note-col=0: column holding notes (0 - not used)
      -sheet=1: number of the sheet containing the translations
      -skipRows=2: number of rows to skip
      -source-col=4: column holding the source for the translation
      -source-lang="en": source language
      -target-col=5: column holding the target translation
      -target-lang="en": target language

    to-json:

      -in="": infile
      -pretty=false: pretty print the resulting json


## Building / Installing

Since *xliffer* is written in go, you need a go compiler. Consult your OS how
to get one or go to http://golang.org/dl.

Once you have a working go compiler:

	$> GOPATH=`pwd` go build -v github.com/travelping/xliffer

You should now have the *xliffer* binary in your working directory.

## Ideas for converters

* Accept OpenDocumentSpreadsheet support
* Convert XLIFF to POT (Portable Object  Template) and PO (Portable Object),
  suitable for gettext

## Related Projects

* http://toolkit.translatehouse.org/


## Limitations

* Currently a subset of XLIFF-1.2 is supported, eg. <groups> are not.

## Contributors

* Mathias Gumz <mg@travelping.com>
