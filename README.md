# PasswordBot
Данные хранятся в базе PostgreSQL, развёрнутой в Docker. Для обеспечения уникальности пространства пользователей в таблице хранится их telegram id, и все пользовательские операции осуществляются только на 
соответствующих им строках. Очистка таблицы происходит в кроне раз в день.