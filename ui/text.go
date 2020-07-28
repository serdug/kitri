// Package ui for GUI (front end)
package ui

const (
	txtAboutProject = `Juggle accounts with Kitri

MADE FOR NON-ACCOUNTANTS

Kitri is an accounting calculator complementing spreadsheets

TRIVIALLY EASY BOOKKEEPING

Kitri bridges the gap between spreadsheets and bookkeeping software

Kitri streamlines the part that is tricky in Microsoft Excel, LibreOffice Calc or Google Sheets 

Kitri takes spreadsheets as input and in the blink of an eye makes double-entry accounts in the form convenient for use in MS Excel

PARAMOUNT DATA SAFETY

Kitri provides users with absolute control over their data, there is no need to upload or share anything

Kitri is an open-source program running on macOS, Linux and Microsoft Windows desktops

Kitri can be run offline with no internet connection whatsoever
`

	txtGuidanceOpen = `Kitri aggregates amounts of recorded transactions by categories in accordance with a user-defined structure of accounts.

CONFIGURATION TEMPLATES

Click 'Next' at the bottom of this page to create and edit a configuration template. Configurations can be saved as .yaml files and are readable from files in the YAML format.

Use 'Select file' to load a saved configuration. The template can be edited and saved on the next page. See an example distributed with the application.

DATA INPUT FORMAT

Kitri takes files with user-defined tables of categories (i.e. the Chart of Accounts) and records of transactions as input. The input files should be prepared beforehand by users in spreadsheets. The files should be saved in the CSV (Comma Separated Values) format. Any single spreadsheet may be saved as a .csv file from MS Excel, LibreOffice Calc or Google Sheets. The CSV format preserves cell values and the structure of rows and columns. Formulas are omitted, although the number formatting remains. So please make sure that the number format is set to General / Automatic before saving data as CSV.

Kitri allows to customize the accounts according to user needs and vary the level of detail appropriately.

CATEGORIES

The following column order must be respected:
• Category – taken as a character string, the category identificator (ID) that must be unique 
• Name - taken as a character string, descriptive category name
• Balance - taken as a number (no thousand separators!), starting balance per category, carried from previous periods

The other columns are ignored, their data does not take part in calculations. It is assumed that the first row of data contains column titles. The first row is ignored by the calculator. So, all columns may be given any names.

RECORDS

The following column order must be respected:
• Amount - taken as a number (no thousand separators!), the monetary value of transaction 
• Source - taken as a character string, category ID; it is where the money goes from
• Purpose - taken as a character string, category ID; it is where the money goes to 

The other columns may contain any comments, notes or explanations. They are ignored by the calculator. It is assumed that the first row of data contains column titles. The first row is ignored. So, all columns may be given any names. 
`

)