# Tarea1-SD

repo: https://github.com/DanielATV/Tarea1-SD

### Integrantes

Gonzalo Larrain 201673516-K
Daniel Toro 201673595-K

### Instrucciones de uso

Para ejecutar la tarea es necesario encontrarse en la carpeta src y ejecutar make build en cada una de 
las maquinas vituales. Luego, en la VM dist01 ejecutar make finanza, en la VM dist04 ejecutar make servidor. Finalmente, make camion en la VM dist02 y make cliente en la VM dist03.

#### Cliente

El cliente pregunta por dos parámetros, el primero de ellos si es cliente tipo pyme (0) o retail(1).
Luego, pregunta por la cantidad de tiempo(en segundos) que espera entre enviar ordernes, la cual es un número entero.

Cuando se escoge la opción de cliente cada 3 ordenes pregunta por el estado de la primera de estas. Por ejemplo, se envian la ordenes A1, A2, A3. Al enviar A3 pregunta por el estado de A1 y este ciclo vuelve a repetirse con A4, A5, A6, en el cual se pregunta por A4 al enviar A6.

#### Logistica(server)

Logística  crea un registo llamado RegistroPedidos el cual lleva el registro histórico de las ordenes que recibe. La estructura timestamp, id-paquete, tipo, nombre, valor, origen, destino, seguimiento. Seguimiento es 0 en caso de ser retail.

#### Camion

Camion pregunta por 2 parámetros el primero de estos es la cantidad de tiempo(en segundos) que se demora en hacer un entrega. Luego, pregunta la cantidad de tiempo(en segundos) que espera por un segundo paquete. Ambos valores deben ser enteros.

El camion 1 y 2 son de retail y el camion 3 es normal. Cada camion lleva un registro de su funcionamiento con la siguiente estructura id-paquete, tipo, valor, origen, destino, intentos, fecha entrega.

#### Finanzas

Finanzas lleva un registro de llamado RegistroFinanciero con la siguiente estructura: Id del paquete, tipo, estado, intentos, ganancia del paquete, perdida del paquete.


#### Consideraciones
- El firewall debe estar desactivado en todas las maquinas virtuales
- Los archivos que lee el cliente deben llamarse pyme.csv o retail.csv
- En la VM dist01 asegurar que el servidor rabbit este activo: systemctl start rabbitmq-server
- La perdida de dignipesos por paquete solo considera el costo de los re intentos




