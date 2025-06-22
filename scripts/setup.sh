#!/bin/bash

COUCH_URL="http://admin:password@localhost:5984/"
DB_NAME=zendo

# Check if database exists
if curl -s -o /dev/null -w "%{http_code}" "$COUCH_URL/$DB_NAME" | grep -q 404; then
  echo "Creating database: $DB_NAME"
  curl -X PUT "$COUCH_URL/$DB_NAME"
else
  echo "Database $DB_NAME already exists"
fi

# Setup non admin API user
curl -X PUT http://admin:password@localhost:5984/_users/org.couchdb.user:api \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api",
    "password": "password",
    "roles": [],
    "type": "user"
  }'

# Give API user read write access to the zendo databse

curl -X PUT "http://admin:password@localhost:5984/${DB_NAME}/_security" \
  -H "Content-Type: application/json" \
  -d '{
    "admins": {
      "names": [],
      "roles": []
    },
    "members": {
      "names": ["api"],
      "roles": []
    }
  }'

# Add the filter design document

curl -X PUT "http://admin:password@localhost:5984/${DB_NAME}/_design/filters" \
  -H "Content-Type: application/json" \
  --data "@couchdb/filters_design_doc.json"

# Add the view design document

curl -X PUT "http://admin:password@localhost:5984/${DB_NAME}/_design/views" \
  -H "Content-Type: application/json" \
  --data "@couchdb/views_design_doc.json"
