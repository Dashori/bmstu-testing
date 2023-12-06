import matplotlib.pyplot as plot

with open('result.txt', 'r') as file:
    data = file.read()

# Удаляем лишние символы и разделяем массив на подмассивы
data = data.replace('[','').replace(']','').replace('"','').replace(',',' ').split('[')
numbers = list(map(int, data[0].split()))

arrays = [[] for _ in range(12)]

# Записываем числа в соответствующие массивы
for i, num in enumerate(numbers):
    arrays[i % 12].append(num)

# Выводим результат
# for i in range(12):
    # print(f'Массив {i+1}: {arrays[i]}')

# 1 -- gorm add NsPerOp
# 2 -- gorm add AllocsPerOp
# 3 -- gorm add AllocedBytesPerOp

# 4 -- gorm get NsPerOp
# 5 -- gorm get AllocsPerOp
# 6 -- gorm get AllocedBytesPerOp

# 7 -- sqlx add NsPerOp
# 8 -- sqlx add AllocsPerOp
# 9 -- sqlx add AllocedBytesPerOp

# 10 -- sqlx get NsPerOp
# 11 -- sqlx get AllocsPerOp
# 12 -- sqlx get AllocedBytesPerOp

size = []
for i in range(100):
    size.append(i + 1)

# add

plot.ylabel("NsPerOp")
plot.xlabel("Номер бенчмарка")
plot.grid(True)
plot.plot(size, arrays[0], color = "darkmagenta", label='gorm', marker='^')
plot.plot(size, arrays[6], color = "gold", label='sqlx', marker='*')
plot.legend(["gorm", "sqlx"])
plot.title('Add client NsPerOp')
plot.savefig('addClientNsPerOp.pdf')
plot.show()

plot.ylabel("AllocsPerOp")
plot.xlabel("Номер бенчмарка")
plot.grid(True)
plot.plot(size, arrays[1], color = "darkmagenta", label='gorm', marker='^')
plot.plot(size, arrays[7], color = "gold", label='sqlx', marker='*')
plot.legend(["gorm", "sqlx"])
plot.title('Add client AllocsPerOp')
plot.savefig('addClientAllocsPerOp.pdf')
plot.show()

plot.ylabel("AllocedBytesPerOp")
plot.xlabel("Номер бенчмарка")
plot.grid(True)
plot.plot(size, arrays[2], color = "darkmagenta", label='gorm', marker='^')
plot.plot(size, arrays[8], color = "gold", label='sqlx', marker='*')
plot.legend(["gorm", "sqlx"])
plot.title('Add client AllocedBytesPerOp')
plot.savefig('addClientAllocedBytesPerOp.pdf')
plot.show()

# get

plot.ylabel("NsPerOp")
plot.xlabel("Номер бенчмарка")
plot.grid(True)
plot.plot(size, arrays[3], color = "darkmagenta", label='gorm', marker='^')
plot.plot(size, arrays[9], color = "gold", label='sqlx', marker='*')
plot.legend(["gorm", "sqlx"])
plot.title('Get client NsPerOp')
plot.savefig('getClientNsPerOp.pdf')
plot.show()

plot.ylabel("AllocsPerOp")
plot.xlabel("Номер бенчмарка")
plot.grid(True)
plot.plot(size, arrays[4], color = "darkmagenta", label='gorm', marker='^')
plot.plot(size, arrays[10], color = "gold", label='sqlx', marker='*')
plot.legend(["gorm", "sqlx"])
plot.title('Get client AllocsPerOp')
plot.savefig('getClientAllocsPerOp.pdf')
plot.show()

plot.ylabel("AllocedBytesPerOp")
plot.xlabel("Номер бенчмарка")
plot.grid(True)
plot.plot(size, arrays[5], color = "darkmagenta", label='gorm', marker='^')
plot.plot(size, arrays[11], color = "gold", label='sqlx', marker='*')
plot.legend(["gorm", "sqlx"])
plot.title('Get client AllocedBytesPerOp')
plot.savefig('getClientAllocedBytesPerOp.pdf')
plot.show()
