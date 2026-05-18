FROM migrate/migrate:v4.18.3

COPY pet-service/migrations /migrations/pet-service
COPY user-service/migrations /migrations/user-service
COPY order-service/migrations /migrations/order-service

COPY migrations-runner/run.sh /run.sh

ENTRYPOINT ["/bin/sh", "/run.sh"]
