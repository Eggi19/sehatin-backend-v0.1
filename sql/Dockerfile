FROM postgis/postgis:16-3.4

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres

COPY ./ddl.sql /docker-entrypoint-initdb.d/001.sql
COPY ./dml.sql /docker-entrypoint-initdb.d/002.sql
COPY ./2_location.sql /docker-entrypoint-initdb.d/003.sql
COPY ./3_products.sql /docker-entrypoint-initdb.d/004.sql
COPY ./4_pharmacies.sql /docker-entrypoint-initdb.d/005.sql
COPY ./5_pharmacy_products.sql /docker-entrypoint-initdb.d/006.sql

HEALTHCHECK CMD pg_isready -U postgres || exit 1