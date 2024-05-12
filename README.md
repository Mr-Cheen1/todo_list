# ToDo List App

<p align="center">
  <img src="todo\static\checklist-1024.webp" alt="ToDo List App" width="300">
</p>



<p align="center">
  <a href="https://github.com/Mr-Cheen1/todo_list/actions/workflows/lint.yml"><img src="https://github.com/Mr-Cheen1/todo_list/actions/workflows/lint.yml/badge.svg" alt="Lint Status"/></a>
  <a href="https://github.com/Mr-Cheen1/todo_list/actions/workflows/test.yml"><img src="https://github.com/Mr-Cheen1/todo_list/actions/workflows/test.yml/badge.svg" alt="Test Status"/></a>
  <a href="https://github.com/Mr-Cheen1/todo_list/actions/workflows/build.yml"><img src="https://github.com/Mr-Cheen1/todo_list/actions/workflows/build.yml/badge.svg" alt="Build Status"/></a>
</p>

# Презентация проекта ToDo List App

## Цели и задачи

Основная цель проекта - создать удобное и функциональное веб-приложение для управления списком задач. 

Задачи:
1. Реализовать базовые функции для работы с задачами: создание, редактирование, обновление, удаление.
2. Обеспечить возможность фильтрации и сортировки задач для удобства пользователей.
3. Разработать интуитивно понятный пользовательский интерфейс.
4. Применить современные технологии и подходы в разработке, такие как Go, PostgreSQL, Docker.
5. Провести тестирование приложения для обеспечения его надежности и стабильности.

## Актуальность

Управление временем и задачами - важный аспект повседневной жизни и работы. Существует множество приложений для управления списками дел, но создание собственного приложения позволяет:

1. **Получить опыт разработки:** Разработка ToDo List App - отличная возможность улучшить навыки программирования и изучить новые технологии.

2. **Кастомизация:** Создавая свое приложение, можно реализовать именно те функции, которые необходимы, и настроить интерфейс по своему вкусу.

3. **Контроль над данными:** При использовании собственного приложения, данные пользователя хранятся на личном сервере, что обеспечивает конфиденциальность и безопасность.

4. **Расширяемость:** Имея доступ к исходному коду, можно легко добавлять новые функции и интегрировать приложение с другими сервисами.

Таким образом, разработка ToDo List App является актуальной задачей, как для личного использования, так и для улучшения навыков разработки веб-приложений.

## Технологии и подход

При разработке ToDo List App были использованы следующие технологии:

1. **Go:** Мощный и эффективный язык программирования для backend разработки.
2. **PostgreSQL:** Надежная реляционная база данных для хранения задач.
3. **HTML, CSS, JavaScript:** Стандартные веб-технологии для создания frontend части приложения.
4. **Rest API:** Архитектурный стиль для построения веб-сервиса.
5. **Docker:** Платформа контейнеризации для упрощения развертывания и масштабирования приложения.

Применена многоуровневая архитектура, разделяющая приложение на frontend, backend и слой для работы с базой данных. Это обеспечивает модульность, масштабируемость и упрощает разработку и поддержку приложения.

Для backend части использовали подход REST API, что позволяет frontend части взаимодействовать с сервером через HTTP запросы. REST API - это архитектурный стиль для построения веб-сервисов. Он основан на использование HTTP методов (GET, POST, PUT, DELETE) для определения операций над ресурсами. Так же были реализованы юнит-тесты и интеграционные тесты для обеспечения стабильности и надежности backend части.

Для frontend части использовали простой и интуитивно понятный дизайн, ориентированный на удобство использования. Применили асинхронные запросы к API для динамического обновления страницы без перезагрузки. Это реализовано с помощью функций `createTask`, `updateTask`, `deleteTask` в файле `todo/static/script.js`, которые отправляют соответствующие HTTP запросы и обновляют список задач после успешного выполнения.

## Обзор

Это приложение для управления списком задач (ToDo List App). Оно позволяет:

1) Создавать новые задачи;
2) Редактировать существующие задачи;
3) При создании и редактировании задачи присутствует валидация пустого ввода, количества допустимых символов в поле ввода (255), а также проверка, что дата предполагаемого завершения задачи не может быть раньше даты создания;
4) Удалять задачи;
5) Изменять статус задач (в процессе/завершено/тестирование/возвращено);
6) Просматривать задачи по критериям: 
   - Только задачи со статусом "в процессе";
   - Только задачи со статусом "завершено";
   - Только задачи со статусом "тестирование";
   - Только задачи со статусом "возвращено";
   - Сортировка задач по дате убывания;
   - Сортировка задач по дате возрастания.
7) Указывать дату предполагаемого завершения задачи при создании и редактировании.

## Презентация проекта
[Файл презентации в PDF](todo/static/Go%20(Golang)%20Developer%20Basic.pdf).

## Запуск приложения с помощью Docker

Прежде чем начать, убедитесь, что у вас установлен Docker. Если он не установлен, посетите официальный сайт Docker для получения инструкций по установке: [Установка Docker](https://docs.docker.com/get-docker/).

1. Клонируйте репозиторий с помощью следующей команды:
`git clone https://github.com/Mr-Cheen1/todo_list.git`

2. Перейдите в каталог проекта где находится файл `docker-compose.yml`:
`cd todo_list/todo`

3. Запустите команду для сборки и запуска контейнеров:
`docker-compose up --build`
(Эта команда соберет образы для серверной части приложения и базы данных, а затем запустит их.)

4. После успешного запуска, приложение будет доступно по адресу http://localhost:8081/ в вашем браузере.

5. Для остановки и удаления контейнеров используйте команду:
`docker-compose down`


## Основные этапы и задачи по разработке приложения для управления списком задач (ToDo List App)

## Backend (Go)
- [x] Настройка подключения к БД PostgreSQL;
- [x] Создание структуры Task для представления задачи;
- [x] Реализация сериализации и десериализации данных:
  - [x] Реализация методов MarshalJSON и UnmarshalJSON для структуры Task;
  - [x] Обработка форматов дат при сериализации и десериализации.
- [x] Реализация HTTP обработчиков:
  - [x] GetTasks - получение списка задач;
  - [x] CreateTask - создание новой задачи;
  - [x] UpdateTask - обновление задачи;
  - [x] DeleteTask - удаление задачи.
- [x] Реализация обработчиков для работы с базой данных:
  - [x] GetAllTasks - получение списка задач из БД;
  - [x] CreateTask - создание задачи в БД;
  - [x] UpdateTask - обновление задачи в БД;
  - [x] DeleteTask - удаление задачи из БД.
- [x] Реализация валидации данных:
  - [x] Валидация пустого ввода при создании и редактировании задачи;
  - [x] Ограничение количества символов в поле ввода при создании и редактировании задачи (255 символов);
  - [x] Проверка, что дата предполагаемого завершения задачи не может быть раньше даты создания.

## Frontend (HTML, CSS, JavaScript)
- [x] Верстка страницы со списком задач;
- [x] Форма для создания новой задачи;
- [x] Отображение списка задач;
- [x] Редактирование задачи;
- [x] Удаление задачи;
- [x] Изменение статуса задачи (в процессе/завершено/тестирование/возвращено);
- [x] Фильтрация задач по статусу;
- [x] Сортировка задач по дате.

## База данных (PostgreSQL)
- [x] Создание таблицы tasks для хранения задач;
- [x] Настройка подключения к БД из приложения.

## Тестирование
- [x] Модульное тестирование обработчиков;
- [x] Интеграционное тестирование приложения.

## Развертывание
- [x] Dockerfile для сборки образа приложения;
- [x] Docker Compose для запуска приложения и БД.

## Cтруктура проекта ToDo List App

- .github/ - Директория с конфигурациями для GitHub Actions.
  - workflows/ - Директория с файлами определения рабочих процессов для GitHub Actions.
    - lint.yml - Файл с настройками для линтинга кода.
    - test.yml - Файл с настройками для запуска тестов.
    - build.yml - Файл с настройками для сборки проекта.
- todo/
  - .golangci.yml - Файл конфигурации для GolangCI Lint.
  - init.d/ - Директория с скриптами инициализации базы данных.
    - create_tables.sql - SQL скрипт для создания таблиц.
  - server/ - Директория с серверной частью приложения на Go.
    - db/ - Директория с файлами для работы с базой данных PostgreSQL.
      - db.go - Файл с функциями для работы с базой данных PostgreSQL.
      - db_handlers_test.go - Файл с тестами для функций работы с базой данных PostgreSQL.
      - db_handlers.go - Файл с обработчиками для операций с базой данных PostgreSQL.
    - handlers/ - Директория с обработчиками HTTP-запросов.
      - task_handlers.go - Файл с обработчиками для операций с задачами.
      - task_handlers_test.go - Файл с тестами для обработчиков задач.
    - main.go - Главный файл серверного приложения.
    - main_test.go - Файл с интеграционными тестами серверного приложения.
    - Dockerfile - Dockerfile для сборки образа серверного приложения.
  - static/ - Директория с клиентской частью приложения (HTML, CSS, JavaScript).
    - index.html - Главная страница приложения.
    - script.js - Файл с JavaScript-кодом для взаимодействия с сервером и управления задачами.
    - styles.css - Файл с CSS-стилями для оформления страницы.
    - checklist-1024.webp - Изображение логотипа приложения.
    - favicon.ico - Иконка загрузки страницы (favicon).
    - Go (Golang) Developer Basic.pdf - Файл с PDF-презентацией проекта.
  - go.mod - Файл с зависимостями Go-проекта.
  - go.sum - Файл с контрольными суммами зависимостей Go-проекта.
  - docker-compose.yml - Файл конфигурации Docker Compose для запуска приложения и базы данных.
  - create_todo_db.bat - Пакетный файл для создания базы данных todo_db в PostgreSQL.
- README.md - Файл с описанием проекта и инструкциями по запуску.