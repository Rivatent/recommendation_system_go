<!-- Improved compatibility of наверх link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![project_license][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/Rivatent/private-go-test-task">
    <img src="images/gopher3.png" alt="Logo" width="150" height="80">
  </a>

<h3 align="center">Микросервисная система рекомендаций</h3>

  <p align="center">
    Тестовое задание на позицию Go-разработчика стажера
    <br />
    <a href="https://github.com/Rivatent/private-go-test-task"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/Rivatent/private-go-test-task">View Demo</a>
    ·
    <a href="https://github.com/Rivatent/private-go-test-task/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/Rivatent/private-go-test-task/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#о-проекте">О проекте</a>
      <ul>
        <li><a href="#использованные-технологии">Использованные технологии</a></li>
      </ul>
    </li>
    <li>
      <a href="#начало-работы">Начало работы</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## О проекте

[![Product Name Screen Shot][product-screenshot]](https://example.com)

Система рекомендаций для пользователей состоит из нескольких микросервисов, которые взаимодействуют между собой с использованием Apache Kafka, реляционной базы данных и кэша для улучшения производительности. Все компоненты системы запускаются с помощью Docker Compose.

<p align="right">(<a href="#readme-top">наверх</a>)</p>



### Использованные технологии

- [Go](https://golang.org/)
- [Apache Kafka](https://kafka.apache.org/)
- [Zookeeper](https://zookeeper.apache.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Prometheus](https://prometheus.io)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Gin](https://github.com/gin-gonic/gin)
- [Validator V10](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Zap Logger](https://github.com/uber-go/zap)

<p align="right">(<a href="#readme-top">наверх</a>)</p>



<!-- GETTING STARTED -->
## Начало работы

### Перед началом наботы

Убедитесь, что на вашей системе установлены следующие компоненты:
- **Docker** (версия 19.03.0+)
- **Docker Compose** (версия  1.25.0+)
- Доступ к сети интернет для загрузки образов


### Использование

- Клонируйте репозиторий

```sh
git clone https://gitverse.ru/sbertech_hr/Go-internship-PavelRalchenkov
cd Go-internship-PavelRalchenkov
```
Каждый микросервис использует .env файл для конфигурации. 
Измените соответствующие файлы для всех микросервисов (user-
service/.env, product-service/.env, analytics-service/.env, 
recommendation-service/.env), в соответствии с желаемыми настройками.

Для запуска всех компонентов выполните следующую команду:
    ```sh
    make run
    ```
  или
    ```sh
    docker-compose up --build -d
    ```
Для остановки всех компонентов выполните:
    ```sh
    make stop
    ```
    или
    ```sh
    docker-compose down
    ```

<p align="right">(<a href="#readme-top">наверх</a>)</p>



<!-- USAGE EXAMPLES -->
## Тестирование
Для тестирования API выполните запросы к эндпоинтам микросервисов с помощью любого HTTP-клиента, например Postman или curl.

Пример запроса для добавления пользователя:
```sh
curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{"username": "New User","email": "new_user@example.com", "password": "new_password"}'
```

Проект покрыт юнит-тестами. Для их запуска выполните команду:
```sh
make tests
```


<!-- ROADMAP -->
## Roadmap

- [ ] Feature 1
- [ ] Feature 2
- [ ] Feature 3
    - [ ] Nested Feature

See the [open issues](https://github.com/Rivatent/private-go-test-task/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">наверх</a>)</p>




<!-- CONTACT -->
## Contact

Павел Ральченков - [@paulralchenkov](https://t.me/paulralchenkov) - paulralchenkov@gmail.com

Project Link: [https://github.com/Rivatent/private-go-test-task](https://github.com/Rivatent/private-go-test-task)

<p align="right">(<a href="#readme-top">наверх</a>)</p>




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/Rivatent/private-go-test-task.svg?style=for-the-badge
[contributors-url]: https://github.com/Rivatent/private-go-test-task/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Rivatent/private-go-test-task.svg?style=for-the-badge
[forks-url]: https://github.com/Rivatent/private-go-test-task/network/members
[stars-shield]: https://img.shields.io/github/stars/Rivatent/private-go-test-task.svg?style=for-the-badge
[stars-url]: https://github.com/Rivatent/private-go-test-task/stargazers
[issues-shield]: https://img.shields.io/github/issues/Rivatent/private-go-test-task.svg?style=for-the-badge
[issues-url]: https://github.com/Rivatent/private-go-test-task/issues
[license-shield]: https://img.shields.io/github/license/Rivatent/private-go-test-task.svg?style=for-the-badge
[license-url]: https://github.com/Rivatent/private-go-test-task/blob/master/LICENSE.txt
[product-screenshot]: images/usrs.png

