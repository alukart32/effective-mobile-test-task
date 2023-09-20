# Тестовое задание для Effective Mobile

Реализовать сервис, который будет получать поток ФИО, из открытых апи обогащать
ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в
БД. По запросу выдавать инфу о найденных людях (persons).

Cервис общается через:

- REST
- GraphQL
- Kafka

## Deployment

Используется Docker compose.

```bash
make docker-up
```

или

```bash
docker compose-up
```