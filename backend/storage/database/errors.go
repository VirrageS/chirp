package database

// TODO this is only for PostgreSQL

// UniqueConstraintViolationCode is error returned by PostgreSQL instance when
// violation of uniqueness in table happens.
const UniqueConstraintViolationCode = "23505"
