## Тестирование!

Мой любимый предмет 7 семестра, но скорее потому что тут было много девопсерских приколов. В какой-то момент у меня был тройной dind, ибо я не хотела делать shell-раннер, но к 4 лабе пришлось.


Я брала проект с курса по ППО 5 семестра на Go, бд Postgres, интерфейс консольный 💅
На каждую лабу отдельная ветка, лучше так и сдавать. Только у меня вроде был какой-то прикол с ветками для первых двух лаб, а точнее из-за одного теста, но это не точно (возможно просто этого теста нет в main). Я сдавала все одной из первой, а вы сами знаете, что это максимальная лояльность + много времени уходило на обсуждение dind.



##### Если не знаете как и что запускать, то все команды можно посмотреть в [.gitlab-ci.yml](https://github.com/Dashori/bmstu-testing/blob/main/.gitlab-ci.yml)


#### Лабораторные работы
1. Unit тесты в internal/services. Там было условие про классический и Лондонский вариант. У меня все было сделано еще на курсе, кроме классического теста. Его я делала в моменте и сейчас я его нахожу только в [коммите](https://github.com/Dashori/bmstu-testing/commit/add3e6abb2231b2c384a9ba72ec1470074004523#diff-670949da573334494a923566f54ce892c39fd4ed04f885bcae07c509bffdf14aR231). И кажется он еще нормально не работает аххаха... С allure там тоже какая-то самодельная история от ребят с потока. В [.gitlab-ci.yml](https://github.com/Dashori/bmstu-testing/blob/main/.gitlab-ci.yml) все команды для запуска.


2. Интеграционные тесты были написаны на курсе ранее, там вроде без проблем, все находится в internal/repository Под E2E тесты отдельная папочка, собственно нужно поднять docker-compose.yml и потом запустить этот тест. Для локального запуска нужно будет сделать пару приседаний, вероятно поменять адрес на localhost в конфиге. Для захвата трафика я просто делала tcpdump и открыла в Wireshark.


3. Тут вроде все норм, есть отдельный [docker-compose-bemch.yml](https://github.com/Dashori/bmstu-testing/blob/main/docker-compose-bench.yaml) со сбором метрик с помощью Prometheus и отрисовкой в Grafana. Я еще нашла отдельный докерфайл для этого, не помню как в итоге сдавала, но думаю там несложно будет разобраться.


4. Моя любимая лаба! Сделала 2FA с помощью почты, использовала личный почтовый сервак Вовы, спасибо ему, ибо с другими почтами были проблемы. Создал мне две кринжовые почты dashylya@huds.su и dashori@huds.su, одна так сказать системная и с нее я отправляла всем письма для входа, а другая пользовательская. 
Логика теста следующая: отправляем запрос с почтой => на эту почту высылается OTP код (но на самом деле это просто [хэш](https://github.com/Dashori/bmstu-testing/blob/main/backend/internal/services/implementation/client.go#L86) от этой же почты) => читаем это (последнее) письмо и берем из него OTP => повторяем запрос с этим OTP. Для хранения паролей от почт я использовала [Vault](https://www.vaultproject.io/), поднимала его локально в docker-контейнере и так как shell-раннер работает на этом же ноуте, то можно получить до него доступ. Настройка есть в [папочке](https://github.com/Dashori/bmstu-testing/tree/main/vault). Перед началом теста я просто клала пароли из волта в [переменные окружения](https://github.com/Dashori/bmstu-testing/blob/main/.gitlab-ci.yml#L67) и дальше в коде их оттуда и [получала](https://github.com/Dashori/bmstu-testing/blob/lab_04/e2e/client_controller/client_otp_test.go#L37). BDD-фреймворк для Go — [ginkgo](github.com/onsi/ginkgo).


5. Я сдавала ее вместе с 4 лабой, так что ее особо не смотрели. Я использовала jaeger, точно помню, что у меня не было красивых вложенных трейсов и все спаны были 1 уровня. Но из прикольного забирала инфу для визуализации [курлом](https://github.com/Dashori/bmstu-testing/blob/lab_05/.gitlab-ci.yml#L53) и клала в артефакты.



----


##### оффтопы


У меня еще были попытки запускать тесты не в gitlab, а в teamcity (решение от JetBrains). Но думаю, что локально поднять gitlab-раннер, а не сервак teamcity и отдельно агент для него — лучшее решение. Но для опыта попробовать настроить teamcity было прикольно, ибо там вся конфигурация задается через ui (если вы не любитель kotlin dsl).


По поводу отчетов, есть [allure](https://github.com/ozontech/allure-go) для Go от ozontech. А условие, что отчет генерируется если тесты падают, легко сделать с помощью дополнительных переменных, но я этим особо не занималась, ибо не спрашивали.

Первые две лабы странные, потому что я не хотела ничего менять с уже существующего проекта и пыталась сдать их с минимальными усилиями, но дальше вроде получше. Надеюсь, оправдалась.
