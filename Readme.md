# Notes API

**Notes API** — это RESTful веб-сервис для создания заметок и управления ими. Он построен на языке Go и использует PostgreSQL в качестве базы данных. API поддерживает аутентификацию пользователей, создание заметок с форматированием текста, вложенные чек-листы и пользовательские таблицы.

## ✨ Основные возможности

*   **Аутентификация пользователей:** Регистрация и вход с использованием JWT (JSON Web Tokens).
*   **CRUD для заметок:** Полный набор операций (Create, Read, Update, Delete) для управления заметками.
*   **Вложенные сущности:**
    *   **Чек-листы:** Добавляйте пункты чек-листа к любой заметке.
    *   **Таблицы:** Создавайте структурированные таблицы с кастомными колонками и строками внутри заметок.
*   **Стилизация текста:** Применяйте к тексту заметок и чек-листов стили `bold` и `italic`.
*   **Документация API:** Автоматически генерируемая документация с помощью Swagger.

## 🛠️ Стек технологий

*   **Язык:** [Go](https://golang.org/) (версия 1.23+)
*   **Веб-фреймворк:** [gorilla/mux](https://github.com/gorilla/mux) для роутинга.
*   **База данных:** [PostgreSQL](https://www.postgresql.org/)
*   **Драйвер БД:** [lib/pq](https://github.com/lib/pq)
*   **Аутентификация:** [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
*   **Документация:** [Swaggo](https://github.com/swaggo/swag)
