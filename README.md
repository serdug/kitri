# Kitri: Trivially Easy Bookkeeping

An accounting aggregator that bridges the gap between spreadsheets and bookkeeping software

Gear up slick spreadsheet office suites for next-level efficiency!

* For non-accountants. Anyone who learns about bookkeeping principles or has basic understanding how accounting works can find it ideal for trial & error accounting.

* A no-brainer tool designed to reconcile trial accounts with spreadsheets of transactions. It enables to create a trial Balance Sheet and a Profit & Loss Statement locally with almost unlimited flexibility in virtually no time: no need to worry about connectivity, upload transactions and store any data in databases of third parties. 

* Kitri aggregates amounts of recorded transactions by accounts in accordance with a user-defined structure

* Kitri streamlines the part that is tricky in Microsoft Excel, LibreOffice Calc or Google Sheets 

* Kitri in the blink of an eye aggregates transactions and returns a worksheet with balance and profit-loss amounts split by user-defined categories. It takes tables of categories (i.e. the Chart of Accounts) and records of transactions as input.

![Example 1: include record](https://github.com/serdug/kitri/blob/master/examples/kitri_example_include.png)
![Example 1: recalculate](https://github.com/serdug/kitri/blob/master/examples/kitri_example_recalc.png)

***

## Output

Kitri returns results in the form convenient for use in Microsoft Excel, LibreOffice Calc or Google Sheets
![Example 1: output](https://github.com/serdug/kitri/blob/master/examples/kitri_example_output.png)


## Usage

Categories and records are accepted in CSV files (Comma Separated Values). Any single spreadsheet from MS Excel, Google Spreadsheets or LibreOffice Calc may be saved as a CSV file. The CSV format preserves cell values and the structure of columns and rows. Formulas are omitted, although the number formatting remains as is. So please make sure that the number format is set to General / Automatic before saving data as CSV.


#### Categories

The following column order must be respected:

* `Category` - character string, a category identificator (ID) that must be unique 
* `Name` - character string, a descriptive category name
* `Balance` - general number (no thousand separators!), a starting balance per category

The other columns are ignored by the calculator. 

It is assumed that the first row of data contains column titles. The first row is ignored by the calculator. So, all columns may be given any names.
![Example 1: output](https://github.com/serdug/kitri/blob/master/examples/kitri_example_input-assets.png)


#### Records

The following column order must be respected:

* `Amount` - general number (no thousand separators!), the monetary value of transaction 
* `Source` - character string, category ID; it is where the money goes from, or is debited from
* `Purpose` - character string, category ID; it is where the money goes to, or is credited to 

The other columns may contain any comments, notes or explanations. They are ignored by the calculator.

It is assumed that the first row of data contains column titles. The first row is ignored by the calculator. So, all columns may be given any names.
![Example 1: output](https://github.com/serdug/kitri/blob/master/examples/kitri_example_input-records.png)


## Examples

The structure of input files and configuration templates can be considered on examples

* [Micro company accounts](https://github.com/serdug/kitri/blob/master/examples/small-no-vat)
* [Config file](https://github.com/serdug/kitri/blob/master/examples/template-ex1.yaml)


## Compilation and Running

* Clone the repository

```
git clone https://github.com/serdug/kitri.git
```

* Compile from source

```
go install github.com/serdug/kitri
```

* Create folders for settings, input and output files

* Copy example files in the folders

* Update the working directory path in the config file accordingly

* Run a compiled package

```
$GOPATH/bin/kitri
```


## Dependencies

#### Prerequisites

You will need a C compiler to compile the application from source and an up-to-date graphics driver to run it.

#### Packages

* [Go](https://go.googlesource.com/go) 1.10+
* [Fyne](https://github.com/fyne-io/fyne) 1.3+ for UI
* [golang.org/x/text](https://github.com/golang/text) 0.3+
