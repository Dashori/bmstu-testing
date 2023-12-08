## End to end!

Нужно в директории до сделать:  
```
docker-compose up -d 
```
убедиться, что postgres и backend поднялись 
выполнить ```go run main.go```

если все успешно, в консоли будет:
```
New client login: aLgevuG
------- 1/4 Successfully create client -------

Ваш логин: aLgevuG
Ваш Id: 530

------- 2/4 Successfully get client info -------

New pet:  LsKqUYc cat 2 8
------- 3/4 Successfully add new pet -------

Ваши питомцы:

 №              Id питомца      Кличка          Тип             Возраст         Уровень здоровья

 1              580             LsKqUYc         cat             2               8

Конец!

------- 4/4 Successfully get pets -------
```

Чтобы захватить трафик во время выполнения:
```
sudo tcpdump -i lo0 port 8080 -w log.pcap
```

и далее я открывала файл log.pcap с помощью wireshark