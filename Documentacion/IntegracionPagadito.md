APIPG - API Pagadito
Pagadito ofrece la APIPG, una API que los Pagadito Comercios podrán integrar a sus plataformas, para conectarse con Pagadito y utilizarlo como medio de cobros.

Pagadito le permite a los Pagadito Comercios realizar cobros de forma rápida y segura, a través de su plataforma de pagos. La tecnología desarrollada en Pagadito, permite la comunicación con múltiples plataformas de forma síncrona, mediante conexiones seguras, únicas y autorizadas.

1. Funciones
La APIPG contiene las siguientes funciones:

Función	PHP	Java
add_detail	X	X
calc_amount	X	X
call	X	X
change_currency_crc	X	 
change_currency_dop	X	 
change_currency_gtq	X	 
change_currency_hnl	X	 
change_currency_nio	X	 
change_currency_pab	X	 
change_currency_usd	X	 
change_format_json	X	X
change_format_php	X	X
change_format_xml	X	X
config	X	X
connect	X	X
construct	X	X
decode_response	X	X
decode_response_extended	 	X
encode_details	 	X
exec_trans	X	X
format_post_vars	X	X
get_exchange_rate_crc	X	 
get_exchange_rate_dop	X	 
get_exchange_rate_gtq	X	 
get_exchange_rate_hnl	X	 
get_exchange_rate_nio	X	 
get_exchange_rate_pab	X	 
get_rs_code	X	X
get_rs_datetime	X	X
get_rs_date_trans	X	X
get_rs_message	X	X
get_rs_reference	X	X
get_rs_status	X	X
get_rs_value	X	X
get_status	X	X
get_xml_element	 	X
get_xml_value	 	X
mode_sandbox_on	X	X
return_attr_response	X	X
return_attr_value	X	X
set_custom_param	X	 
 

1.1 add_detail
(PHP, Java)

Descripción
Agrega un detalle a la orden de cobro, previo a su ejecución.

Parámetros
quantity
(int) Define la cantidad del producto.
description
(String) Define la descripción del producto.
price
(Double) Define el precio del producto en términos de dólares americanos (USD).
url_product
(String) Define la url de referencia del producto.
 

1.2 calc_amount
(PHP, Java)

Descripción
Devuelve la sumatoria de los productos entre cantidad y precio de todos los detalles de la transacción.

Valores devueltos
(Double) La sumatoria.

 

1.3 call
(PHP, Java)

Descripción
Ejecuta una llamada a Pagadito y devuelve la respuesta.

Parámetros
params
(PHP: Array / Java: HashMap) Variables y sus valores a enviarse en la llamada.
extended (Solo Java)
(Boolean) extended Define si es una respuesta extendida la que se recibirá.
Valores devueltos
(PHP: String / Java: HashMap) La cadena devuelta por Pagadito.

 

1.4 change_currency_crc
(PHP)

Descripción
Cambia la moneda en uso a colón costarricense.

 

1.5 change_currency_dop
(PHP)

Descripción
Cambia la moneda en uso a peso dominicano.

 

1.6 change_currency_gtq
(PHP)

Descripción
Cambia la moneda en uso a quetzal.

 

1.7 change_currency_hnl
(PHP)

Descripción
Cambia la moneda en uso a lempira.

 

1.8 change_currency_nio
(PHP)

Descripción
Cambia la moneda en uso a córdoba.

 

1.9 change_currency_pab
(PHP)

Descripción
Cambia la moneda en uso a balboa.

 

1.10 change_currency_usd
(PHP)

Descripción
Cambia la moneda en uso a dólar americano.

 

1.11 change_format_json
(PHP, Java)

Descripción
Cambia el formato de retorno a JSON.

 

1.12 change_format_php
(PHP, Java)

Descripción
Cambia el formato de retorno a PHP.

 

1.13 change_format_xml
(PHP, Java)

Descripción
Cambia el formato de retorno a XML.

 

1.14 config
(PHP, Java)

Descripción
Establece los valores por defecto.

 

1.15 connect
(PHP, Java)

Descripción
Conecta con Pagadito y autentica al Pagadito Comercio.

Valores devueltos
(Boolean) Devuelve true si realizó la conexión exitosamente. De lo contrario devuelve false.

 

1.16 construct
(PHP, Java)

Descripción
Constructor de la clase, el cual inicializa los valores por defecto.

Parámetros
uid
(String) El identificador del Pagadito Comercio.
wsk
(String) La clave de acceso.
 

1.17 decode_response
(PHP, Java)

Descripción
Devuelve un objeto con los datos de la respuesta de Pagadito.

Parámetros
response
(String) Cadena contenedora de la estructura a ser decodificada.
Valores devueltos
(PHP: Object / Java: HashMap) Estructura con los datos devueltos por Pagadito.

 

1.18 decode_response_extended
(Java)

Descripción
Devuelve un objeto con los datos de la respuesta de Pagadito.

Parámetros
response
(String) Cadena contenedora de la estructura a ser decodificada.
Valores devueltos
(HashMap) Estructura con los datos devueltos por Pagadito.

 

1.19 encode_details
(Java)

Descripción
Devuelve una cadena con el formato válido de los detalles de los productos a enviar en una llamada.

Parámetros
details
(List) Estructura con los detalles de la compra.
Valores devueltos
(String) Los detalles de la compra en formato de cadena.

 

1.20 exec_trans
(PHP, Java)

Descripción
Solicita el registro de la transacción y redirecciona a la pantalla de cobros de Pagadito.

Parámetros
ern
(String) External Reference Number - Número de Referencia Externa.
Valores devueltos
(Boolean) Devuelve true si se registró la transacción correctamente. De lo contrario devuelve falase.

 

1.21 format_post_vars
(PHP, Java)

Descripción
Devuelve una cadena con el formato válido de variables y valores para enviar en una llamada.

Parámetros
vars
(PHP: Array / Java: HashMap) Variables con sus valores a ser formateados.
Valores devueltos
(String) Variables con sus valores en formato de cadena.

 

1.22 get_exchange_rate_crc
(PHP)

Descripción
Devuelve la tasa de cambio del colón costaricense.

 

1.23 get_exchange_rate_dop
(PHP)

Descripción
Devuelve la tasa de cambio del peso dominicano.

 

1.24 get_exchange_rate_gtq
(PHP)

Descripción
Devuelve la tasa de cambio del quetzal.

 

1.25 get_exchange_rate_hnl
(PHP)

Descripción
Devuelve la tasa de cambio del lempira.

 

1.26 get_exchange_rate_nio
(PHP)

Descripción
Devuelve la tasa de cambio del córdoba.

 

1.27 get_exchange_rate_pab
(PHP)

Descripción
Devuelve la tasa de cambio del balboa.

 

1.28 get_rs_code
(PHP, Java)

Descripción
Devuelve el código de la respuesta.

Valores devueltos
(String) Código de la respuesta.

 

1.29 get_rs_datetime
(PHP, Java)

Descripción
Devuelve la fecha y hora de la respuesta.

Valores devueltos
(String) Fecha y hora de la respuesta.

 

1.30 get_rs_date_trans
(PHP, Java)

Descripción
Devuelve la fecha y hora de la transacción consultada, después de un get_status().

Valores devueltos
(String) Fecha y hora de la transacción consultada.

 

1.31 get_rs_message
(PHP, Java)

Descripción
Devuelve el mensaje de la respuesta.

Valores devueltos
(String) Mensaje de la respuesta.

 

1.32 get_rs_reference
(PHP, Java)

Descripción
Devuelve la referencia de la transacción consultada, después de un get_status().

Valores devueltos
(String) Referencia de la transacción consultada.

 

1.33 get_rs_status
(PHP, Java)

Descripción
Devuelve el estado de la transacción consultada, después de un get_status().

Valores devueltos
(String) Estado de la transacción consultada.

 

1.34 get_rs_value
(PHP, Java)

Descripción
Devuelve el valor de la respuesta.

Valores devueltos
(String) Valor de la respuesta.

 

1.35 get_status
(PHP, Java)

Descripción
Solicita el estado de una transacción en base a su token.

Parámetros
token_trans
(String) El identificador de la conexión a consultar.
Valores devueltos
(Boolean) Devuelve true si consultó exitosamente. De lo contrario devuelve false.

 

1.36 get_xml_element
(Java)

Descripción
Devuelve el objeto del ítem solicitado.

Parámetros
element
(Element) Objeto con datos provenientes de un XML.
item
(String) Nombre del ítem solicitado.
Valores devueltos
(Element) Objeto hijo de un XML.

 

1.37 get_xml_value
(Java)

Descripción
Devuelve el valor del ítem solicitado.

Parámetros
element
(Element) Objeto con datos provenientes de un XML.
item
(String) Nombre del ítem solicitado.
Valores devueltos
(String) Valor de un elemento XML.

 

1.38 mode_sandbox_on
(PHP, Java)

Descripción
Habilita el modo de pruebas SandBox.

 

1.39 return_attr_response
(PHP, Java)

Descripción
Devuelve el valor del atributo solicitado.

Parámetros
attr
(String) Nombre del atributo de la respuesta.
Valores devueltos
(String) Valor de un atributo de la respuesta proveniente de Pagadito.

 

1.40 return_attr_value
(PHP, Java)

Descripción
Devuelve el valor del atributo solicitado.

Parámetros
attr
(String) Nombre del atributo del valor devuelto en la respuesta.
Valores devueltos
(String) Valor de un atributo del valor de la respuesta proveniente de Pagadito.

 

1.41 set_custom_param
(PHP)

Descripción
Establece el valor que tomará el parámetro personalizado especificado en la orden de cobro, previo a su ejecución.

Parámetros
code
(String) Código del parámetro a enviar.
value
(String) Define el valor que se asignará al parámetro.
