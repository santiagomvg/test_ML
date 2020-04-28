# test_ML
###MercadoLibre Ejercicio
Servicio para obtener diversa informacion dada un direccion IP. <br/>
Con el uso va cacheando en un server Redis los response de algunas APIs.
Dichos caches se usan para evitar llamar a las APIs externas todo el tiempo,
y tambien persiste informacion en Redis para ejecutar las APIs estadisticas
eficientemente por medio de Sets y hashes incrementales.

###Ejecucion

cd %GOPATH%/src<br/>
git clone https://github.com/santiagomvg/test_ML <br/>
cd test_ML<br/>

editar config.json para apuntar a un servidor de Redis

go get -t<br/>
go run *.go<br/>

Abrir browser en http://localhost:5000

###Arquitectura de desarrollo
Linux Mint 19.3 - 64bits <br/>
GO 1.14.2 <br/>
Redis 4.0.9

###Dependencias
No utilice vendors. <br/>
el comando "go get" las resolverá. <br/>
El servicio solo depende de las librerias gin-gonic y redigo. <br/>
* No contiene clases de testeo