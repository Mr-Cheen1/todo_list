@echo off

REM. Установка пароля для подключения к базе данных PostgreSQL.
set PGPASSWORD=4217

REM. Создание базы данных todo_db с использованием psql.
psql -U postgres -p 8080 -c "CREATE DATABASE todo_db;"