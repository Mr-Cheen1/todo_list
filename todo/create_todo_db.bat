@echo off
set PGPASSWORD=4217
psql -U postgres -p 8080 -c "CREATE DATABASE todo_db;"