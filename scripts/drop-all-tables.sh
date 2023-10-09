#!/bin/bash

if [ "${GERARDOUS_DB}" == "" ] ; then
  echo "The env var GERARDOUS_DB not set"
  exit 1
fi

if [ ! -f "${GERARDOUS_DB}" ] ; then
  echo "The filepath '${GERARDOUS_DB}' DOES NOT EXIST; filepath from env var GERARDOUS_DB."
  exit 2
fi

# Get the list of all tables in the SQLite database
tables=$(sqlite3 "${GERARDOUS_DB}" \
  "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE '%sqlite_%';")

if [ "${tables}" == "" ]; then
  echo "No tables to drop."
  exit 0
fi

# Loop through each table name and drop it
for table in $tables; do
  echo "Dropping ${table}..."
  sqlite3 "${GERARDOUS_DB}" "DROP TABLE ${table};"
done
echo "Done."
