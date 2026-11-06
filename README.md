# github-activity

Простая CLI программа для получения событий пользователя GitHub через обращение к API "https://api.github.com/users/username/events".

## Запуск программы

go run github-activity.go username

## Особенности реализации

- Таймаут запроса 5 секунд
- Дубликаты событий, идущие подряд, сворачиваются в одно событие с указанием количества повторов
- Обрабатываются следующие уникальные события:
- - PushEvent
- - WatchEvent
- - IssuesEvent
- - CreateEvent
- - DeleteEvent
- - ReleaseEvent
- - ForkEvent
- - PublicEvent

Остальные события проходят через стандартный обработчик
