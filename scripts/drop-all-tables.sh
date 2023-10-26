#!/bin/bash

if [ "${GERARDUS_DB}" == "" ] ; then
  echo "The env var GERARDUS_DB not set"
  exit 1
fi

if [ ! -f "${GERARDUS_DB}" ] ; then
  echo "The filepath '${GERARDUS_DB}' DOES NOT EXIST; filepath from env var GERARDUS_DB."
  exit 2
fi

# Get the list of all tables in the SQLite database
tables=$(sqlite3 "${GERARDUS_DB}" \
  "SELECT name FROM sqlite_master WHERE type='table' AND ( name NOT LIKE '%sqlite_%' AND name<>'project' AND name<>'codebase');")

if [ "${tables}" == "" ]; then
  echo "No tables to drop."
  exit 0
fi

# Loop through each table name and drop it
for table in $tables; do
  echo "Dropping ${table}..."
  sqlite3 "${GERARDUS_DB}" "DROP TABLE ${table};"
done
echo "Done."
