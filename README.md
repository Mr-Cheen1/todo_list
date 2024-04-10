# ToDo List App

<p align="center">
  <img src="todo\static\checklist-1024.webp" alt="ToDo List App" width="300">
</p>

<table>
  <tr>
    <td>
      <a href="https://github.com/Mr-Cheen1/todo/releases">
        <img src="https://img.shields.io/github/release/Mr-Cheen1/todo.svg" alt="GitHub Release">
      </a>
    </td>
    <td>
      <a href="https://goreportcard.com/report/github.com/Mr-Cheen1/todo">
        <img src="https://goreportcard.com/badge/github.com/Mr-Cheen1/todo" alt="Go Report Card">
      </a>
    </td>
    <td>
      <a href="https://coveralls.io/github/Mr-Cheen1/todo?branch=main">
        <img src="https://coveralls.io/repos/github/Mr-Cheen1/todo/badge.svg?branch=main" alt="Coverage Status">
      </a>
    </td>
    <td>
      <a href="https://codecov.io/gh/Mr-Cheen1/todo">
        <img src="https://codecov.io/gh/Mr-Cheen1/todo/branch/main/graph/badge.svg" alt="codecov">
      </a>
    </td>
  </tr>
</table>

Это приложение для управления списком задач (ToDo List). Оно позволяет:

1) Создавать новые задачи (Присутствует валидация пустого ввода, а так же кол-во допустимых символов в поле ввода 255);
2) Редактировать существующие задачи (Присутствует валидация пустого ввода, а так же кол-во допустимых символов в поле ввода 255);
3) Удалять задачи;
4) Изменять статус задач (в процессе/завершено);
5) Просматривать задачи по критериям: 
  - Только задачи со статусом "в процессе";
  - Только задачи со статусом "завершено";
  - Сортировка задач по дате убывания;
  - Сортировка задач по дате возрастания.


## Технологии

- Backend: Go.
- Frontend: HTML, CSS, JavaScript.
- База данных: PostgreSQL.


## Запуск приложения

1. Клонируйте репозиторий 
2. Установите зависимости
3. Настройте подключение к БД PostgreSQL
4. Запустите backend: `go run server/main.go`
5. Откройте frontend: `todo/static/index.html`

