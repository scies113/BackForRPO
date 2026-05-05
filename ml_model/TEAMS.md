# Команды, которые распознаёт ML-модель

ML-модель обучена на данных матчей **English Premier League** с сезона 2000/01 по 2024/25.

> **Важно:** При создании матча в системе названия команд должны **точно совпадать** с написанием ниже (на английском). Если ввести название на русском или с другим написанием — ML-модель вернёт ошибку.

## Список команд (39 команд)

| # | Команда | # | Команда |
|---|---------|---|---------|
| 1 | Arsenal | 21 | Man City |
| 2 | Aston Villa | 22 | Man United |
| 3 | Birmingham | 23 | Middlesbrough |
| 4 | Blackburn | 24 | Newcastle |
| 5 | Bolton | 25 | Nott'm Forest |
| 6 | Bournemouth | 26 | Norwich |
| 7 | Bradford | 27 | QPR |
| 8 | Brentford | 28 | Reading |
| 9 | Brighton | 29 | Sheffield United |
| 10 | Burnley | 30 | Southampton |
| 11 | Cardiff | 31 | Stoke |
| 12 | Charlton | 32 | Sunderland |
| 13 | Chelsea | 33 | Swansea |
| 14 | Coventry | 34 | Tottenham |
| 15 | Crystal Palace | 35 | Watford |
| 16 | Derby | 36 | West Brom |
| 17 | Everton | 37 | West Ham |
| 18 | Fulham | 38 | Wigan |
| 19 | Hull | 39 | Wolves |
| 20 | Ipswich | | |
| | Leeds | | |
| | Leicester | | |
| | Liverpool | | |

## Примеры правильного ввода

```
Home Team: Arsenal
Away Team: Chelsea

Home Team: Man United
Away Team: Liverpool

Home Team: Nott'm Forest
Away Team: Brighton
```

## Сокращения команд

| Полное название | В ML-модели |
|-----------------|-------------|
| Manchester City | **Man City** |
| Manchester United | **Man United** |
| Nottingham Forest | **Nott'm Forest** |
| West Bromwich Albion | **West Brom** |
| Sheffield United | **Sheffield United** |
| Queens Park Rangers | **QPR** |
| Crystal Palace | **Crystal Palace** |
| West Ham United | **West Ham** |

## Данные для обучения

- **Источник:** EPL Final Dataset (`epl_final.csv`)
- **Период:** сезоны 2000/01 — 2024/25
- **Всего матчей:** ~9 381
- **Статистика:** голы, удары, удары в створ, угловые, фолы, жёлтые/красные карточки
- **Модель:** XGBoost (multi:softprob, 3 класса: победа хозяев / ничья / победа гостей)
