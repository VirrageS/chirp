#!/bin/bash
set -e

psql -U postgres -f /schema.sql
psql -U postgres -f /fixtures.sql
