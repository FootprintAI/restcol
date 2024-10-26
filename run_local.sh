#!/usr/bin/env bash
#
go build -o restcol && \
    ./restcol --restcol_auto_migrate \
    --restcol_db_endpoint=localhost:5432 \
    --restcol_db_name=unittest \
    --restcol_db_user=postgres \
    --restcol_db_password=password


# open default project url
# open http://localhost:50091/v1/projects/1001/apidoc
# open swaggerui http://localhost:50091/swaggerui/
