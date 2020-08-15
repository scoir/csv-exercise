# Assignment
## User Stories

### As a User, given a csv, I can produce a json file of valid csv records
* csv files will be placed into _`input-directory`_
    * once the application starts it watches _`input-directory`_ for any new files that need to be processed
    * files will be considered new if the file name has not been recorded as processed before.
    * file names will end in `.csv`
    * once a file has been processed, the system deletes it from the _`input-directory`_.
* csv columns and validation
    * csv files will have headers.
    * columns
        1. `INTERNAL_ID` : 8 digit positive integer. Cannot be empty.
        1. `FIRST_NAME` : 15 character max string. Cannot be empty.
        1. `MIDDLE_NAME` : 15 character max string. Can be empty.
        1. `LAST_NAME` : 15 character max string. Cannot be empty.
        1. `PHONE_NUM` : string that matches this pattern `###-###-####`. Cannot be empty.
* json files should be output to _`output-directory`_
    * file name should be the same name as the csv file with a `.json` extension
    * in the event of file name collision, the latest file should overwrite the earlier version.
    * The middle field should not exist if there is no middle name.
* json format:
```js
[
    {
        "id": <INTERNAL_ID>,
        "name": {
            "first": "<FIRST_NAME>",
            "middle": "<MIDDLE_NAME>",
            "last": "<LAST_NAME>"
        },
        "phone": "<PHONE_NUM>"
    }
]
```

#### Example

input of:

```
INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,Tables,555-555-5555
```

would produce:

```json
[
    {
        "id": 12345678,
        "name": {
            "first": "Bobby",
            "last": "Tables"
        },
        "phone": "555-555-5555"
    }
]
```
---

### As a User, I can produce a csv file containing validation errors
* error records should be written to a csv file in _`error-directory`_
* if errors exist, one error file is created per input file.
* processing should continue in the event of an invalid row; all errors should be collected and added to the corresponding error csv.
* an error record should contain:
    1. `LINE_NUM` : the number of the record which was invalid
    * `ERROR_MSG` : a human readable error message about what validation failed
* in the event of name collision, the latest file should overwrite the earlier version.
* the error file should match the name of the input file.

---

### As a User, I can configure input, output, and error directories
