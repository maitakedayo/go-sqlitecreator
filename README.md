# go-sqlitecreator

This Go program provides functionality for handling a SQLite database and performing various operations such as data insertion, querying, and calculations.

## Overview

The application encapsulates several essential operations within its methods. Here's a summary of the available methods:

- `OpenDatabase`: Establishes a connection to the SQLite database.
- `CreateTable`: Creates the 'entries' table if it doesn't exist.
- `ProcessEmbeddedFile`: Reads data from a 'money.txt' file and inserts it into the database.
- `IsDateFormatValid`: Validates the date format for data in the file.
- `OutputEntries`: Queries the database and displays the first five entries in the 'entries' table.
- `CalculateTotalAmount`: Computes the total amount of all entries in the 'entries' table.
- `CalculateTotalAmountPerMonth`: Calculates the total amount for each month from the 'entries' table.
- `CalculateAverageTotalAmountPerMonth`: Computes the average total amount per month.
- `CalculateTotalAmountPerMonthExcludingYatin`: Retrieves the total amount per month while excluding 'yatin' category entries.
- `CalculateAverageTotalAmountPerMonthExcludingYatin`: Calculates the average total amount per month, excluding 'yatin' category entries.

## Usage

1. Establish a connection to the database using `OpenDatabase`.
2. Create the 'entries' table using `CreateTable`.
3. Read data from 'money.txt' and insert it into the database with `ProcessEmbeddedFile`.
4. Utilize various methods to perform queries and calculations on the database.

## Installation

```bash
$ go get github.com/maitakedayo/go-sqlitecreator
```

## License

MIT

## Author

maitakedayo

## ライセンス

このプロジェクトは [MIT ライセンス](LICENSE) のもとで公開されています。