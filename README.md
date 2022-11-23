# ALUMNI API REST
Antes de inicar el servidor debe asegurarse de tener el archivo .env donde esta las configuraciones tando de postgresql y Bucket de aws

# Base de datos
Esta API REST utiliza el ORM GORM y el motor de base de datos postgres, para que GORM puede trabajar con la base de datos, primero se debe crear la base de datos.
``` sql
create database alumni_prueba;
```

## Para iniciar el servidor
``` bash
$ go run main.go
```
# Técnologias utilizadas
- Go lang
- GORM
- Viper
- Fiber

## Desarrollador
- Roberto Suárez
- email: roberto.suarez.job@gmail.com