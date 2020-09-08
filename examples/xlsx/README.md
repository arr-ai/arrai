# XLSX

For code review purposes, `foo.xlsx` is an Excel workbook with two sheets. Sheet 0 is empty, and sheet 1 is roughly as follows:

| Customer ID | AccountType |         | Account, Balance |                    |
| ----------- | ----------- | ------- | ---------------- | ------------------ |
| foo         | Checking    | empty   | 100              | Commas ignored     |
|             |             |         |                  | Empty rows ignored |
| bar         | Checking    | columns | 200              | No header skipped  |
|             | Savings     | ignored | 300              |                    |
|             |             |         |                  |                    |

It touches a variety of edge cases that XLSX decoding is designed to handle, including:

- Columns with empty headers are ignored.
- Rows with no values are ignored.
- Values after the last header column are ignored.
- Formatting is ignored altogether.
- Column names are snaked_cased, with capital letters used as WorkBreaks (`=> word_breaks`).
- Merged cells are treated as though the merged value was copied to every merged cell.

The output is the relation:

```arrai
{
    |@row, customer_id, account_type, account_balance|
    (1   , 'foo'      , 'Checking'  , '100'          ), 
    (3   , 'bar'      , 'Checking'  , '200'          ),
    (4   , 'bar'      , 'Savings'   , '300'          ),
}
```
