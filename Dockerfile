FROM postgres:11.3
from mdillon/postgis

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 8080

