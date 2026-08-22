# Validaciones de estructuras con `httpx.BodyParser`

Esta guía explica cómo funciona hoy la validación de estructuras en `uker` cuando se usa `httpx.BodyParser`, qué formato de payload espera, qué valida exactamente la tag `uker:"required"` y qué limitaciones conviene tener presentes al implementarlo en backends.

La documentación está basada en la implementación actual de [`request.go`](D:\Code\uker\uker\httpx\request.go) y [`validate.go`](D:\Code\uker\uker\validate\validate.go).

## Objetivo

`httpx.BodyParser` resuelve dos pasos en una sola llamada:

1. Lee y deserializa el body HTTP hacia una estructura Go.
2. Ejecuta una validación mínima de campos obligatorios usando tags `uker:"required"`.

Esto permite centralizar el parseo básico del request antes de aplicar validaciones de negocio más específicas.

## Firma

```go
func BodyParser(r *http.Request, target any, opts ...ParserOption) error
```

Uso típico:

```go
type CreateUserRequestDTO struct {
    Email string `json:"email" uker:"required"`
    Name  string `json:"name"`
    Age   int    `json:"age" uker:"required"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequestDTO

    if err := httpx.BodyParser(r, &req); err != nil {
        httpx.ErrorOutput(w, http.StatusBadRequest, httpx.Response{
            Status: httpx.ResponseStatus{
                Type:        httpx.Error,
                Code:        "bad_request",
                Description: err.Error(),
            },
        })
        return
    }

    // Validaciones de negocio adicionales.
}
```

## Requisito principal: `target` debe ser un puntero a struct

`BodyParser` espera recibir un puntero. Si se pasa una estructura por valor, hace `panic`.

Ejemplo correcto:

```go
var req CreateUserRequestDTO
err := httpx.BodyParser(r, &req)
```

Ejemplo incorrecto:

```go
var req CreateUserRequestDTO
err := httpx.BodyParser(r, req)
```

Comportamiento actual:

- Si `target` no es puntero, `BodyParser` hace `panic`.
- Si `target` es puntero pero no apunta a un `struct`, la validación posterior devuelve error.

## Flujo interno

El flujo real de `BodyParser` es este:

1. Verifica que `target` sea puntero.
2. Lee todo el body con `io.ReadAll`.
3. Intenta deserializar el body completo a `map[string]any`.
4. Si existe la clave `data`, procesa ese contenido.
5. Si no existe `data`, deserializa el body completo directamente sobre `target`.
6. Ejecuta `validate.RequiredFields(target, dataFields)`.

En otras palabras, la validación no se hace sobre bytes crudos ni sobre tags JSON solamente: se hace comparando la estructura ya deserializada con el mapa de campos que llegó en el request.

## Formatos de payload soportados

Hoy `BodyParser` soporta dos formas de entrada.

### 1. Payload plano en la raíz

Es el formato más directo.

```json
{
	"email": "john@example.com",
	"name": "John",
	"age": 30
}
```

Con este formato:

- `BodyParser` deserializa el body completo directamente sobre tu struct.
- La validación toma como fuente el mapa raíz del JSON.

Ejemplo:

```go
type CreateUserRequestDTO struct {
    Email string `json:"email" uker:"required"`
    Name  string `json:"name"`
    Age   int    `json:"age" uker:"required"`
}
```

### 2. Payload anidado en `data`

Si el body contiene una clave `data`, `BodyParser` intenta parsear ese valor como la carga real.

Ejemplo con `data` como objeto JSON:

```json
{
	"data": {
		"email": "john@example.com",
		"name": "John",
		"age": 30
	}
}
```

Ejemplo con `data` como string JSON:

```json
{
	"data": "{\"email\":\"john@example.com\",\"name\":\"John\",\"age\":30}"
}
```

En ambos casos:

- `BodyParser` toma solo el contenido de `data`.
- Ese contenido se deserializa sobre `target`.
- La validación obligatoria se hace usando los campos internos de `data`, no el body externo.

## Soporte para `data` en base64

Si el backend recibe `data` en base64, se debe usar la opción `httpx.WithBase64Data()`.

Ejemplo:

```json
{
	"data": "eyJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJuYW1lIjoiSm9obiIsImFnZSI6MzB9"
}
```

Uso:

```go
var req CreateUserRequestDTO
err := httpx.BodyParser(r, &req, httpx.WithBase64Data())
```

Comportamiento:

- `BodyParser` lee `data`.
- Decodifica base64.
- Interpreta el resultado como JSON.
- Deserializa el JSON sobre la estructura destino.
- Ejecuta la validación `required`.

Si el base64 es inválido, retorna:

```text
malformed base64 on data field
```

## Cómo funciona `uker:"required"`

La validación de obligatorios se implementa en `validate.RequiredFields`.

```go
type CreateUserRequestDTO struct {
    Email string `json:"email" uker:"required"`
    Name  string `json:"name"`
    Age   int    `json:"age" uker:"required"`
}
```

La regla actual es simple:

- Solo se validan los campos que tengan `uker:"required"`.
- Para encontrar la clave en el body, primero usa la tag `json`.
- Si no existe tag `json`, usa el nombre del campo Go.
- Si el campo no está presente en el mapa o su valor es `nil`, retorna error.
- Si el campo es `string` y quedó vacío luego del unmarshal, también retorna error.

## Resolución del nombre del campo

La validación busca la clave usando este orden:

1. `json:"..."` si existe.
2. Nombre del campo Go si no hay tag `json`.

Ejemplo:

```go
type RequestDTO struct {
    UserID string `json:"user_id" uker:"required"`
    Name   string `uker:"required"`
}
```

JSON esperado:

```json
{
	"user_id": "abc",
	"Name": "John"
}
```

Punto importante:

- Si un campo no tiene tag `json`, el nombre esperado en el body será exactamente el nombre del struct, por ejemplo `Name`.
- Para APIs públicas conviene definir siempre la tag `json`.

## Qué valida exactamente y qué no

La validación actual cubre esto:

- Presencia del campo requerido en el payload.
- Que un `string` requerido no quede vacío.

La validación actual no cubre esto:

- Rangos numéricos.
- Longitud mínima o máxima automática sobre fields del struct.
- Formato de email, UUID, fecha, etc.
- Validación recursiva de structs anidados.
- Validación automática de slices, maps o punteros internos.
- Diferenciación entre un entero faltante y un entero con valor cero cuando la clave sí vino en el request.

## Comportamiento por tipo de dato

### `string`

Si el campo requerido:

- no viene en el JSON, falla.
- viene en `null`, falla.
- viene como `""`, falla.

Ejemplo que falla:

```json
{
	"email": ""
}
```

### `int`, `bool`, `float`, etc.

Si el campo requerido:

- no viene en el JSON, falla.
- viene en `null`, falla.
- viene con su valor cero válido, pasa.

Ejemplos que pasan:

```json
{
	"age": 0
}
```

```json
{
	"enabled": false
}
```

Esto ocurre porque la validación revisa valor vacío solo para `string`.

## Ejemplos prácticos

### Caso válido con body raíz

```go
type CreateBusinessRequestDTO struct {
    Name   string `json:"name" uker:"required"`
    Active bool   `json:"active" uker:"required"`
    Seats  int    `json:"seats" uker:"required"`
}
```

```json
{
	"name": "Acme",
	"active": false,
	"seats": 0
}
```

Resultado:

- `name` pasa porque no está vacío.
- `active` pasa aunque sea `false`.
- `seats` pasa aunque sea `0`.

### Caso inválido por campo faltante

```json
{
	"active": true,
	"seats": 10
}
```

Resultado:

```text
missing required parameter: Name
```

El mensaje usa el nombre del campo Go, no la key JSON.

### Caso inválido por string vacío

```json
{
	"name": "",
	"active": true,
	"seats": 10
}
```

Resultado:

```text
missing required parameter: Name
```

## Errores que puede devolver `BodyParser`

Errores frecuentes:

- `cannot read request body`
- `error happend on json unmarshal of request body`
- `missing field 'data' inside of the request`
- `malformed base64 on data field`
- `error unmarshalling JSON: ...`
- `missing required parameter: <FieldName>`

Notas:

- Hay mensajes con typos como `happend`. La guía refleja el comportamiento actual.
- Si el request no tiene JSON válido, falla antes de validar obligatorios.

## Orden real de validación

El orden importa para entender ciertos errores:

1. Primero se hace el unmarshal.
2. Después se valida `required`.

Eso implica:

- Si el tipo JSON no coincide con la estructura Go, el error ocurre antes de la validación.
- Si el JSON parsea correctamente pero falta un campo requerido, el error viene de `RequiredFields`.

Ejemplo:

```json
{
	"age": "treinta"
}
```

Si `age` es `int`, el error será de unmarshal, no de `required`.

## Recomendación de uso en backends

Patrón recomendado:

1. Definir un DTO por endpoint.
2. Taggear con `json:"..."` todos los campos expuestos.
3. Marcar con `uker:"required"` solo los obligatorios de transporte.
4. Ejecutar `httpx.BodyParser`.
5. Aplicar después validaciones de negocio con `validate` u otra capa propia.

Ejemplo:

```go
type CreateMemberRequestDTO struct {
    BusinessID string `json:"business_id" uker:"required"`
    Email      string `json:"email" uker:"required"`
    Name       string `json:"name" uker:"required"`
    Role       string `json:"role"`
}

func CreateMemberHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateMemberRequestDTO

    if err := httpx.BodyParser(r, &req); err != nil {
        httpx.ErrorOutput(w, http.StatusBadRequest, httpx.Response{
            Status: httpx.ResponseStatus{
                Type:        httpx.Error,
                Code:        "invalid_body",
                Description: err.Error(),
            },
        })
        return
    }

    if err := validate.NotEmpty(req.Email); err != nil {
        httpx.ErrorOutput(w, http.StatusBadRequest, httpx.Response{
            Status: httpx.ResponseStatus{
                Type:        httpx.Error,
                Code:        "invalid_email",
                Description: err.Error(),
            },
        })
        return
    }
}
```

## Buenas prácticas

- Usar siempre structs específicos por endpoint.
- Declarar siempre tags `json`.
- Reservar `uker:"required"` para obligatoriedad del contrato HTTP.
- Dejar validaciones de dominio en una segunda capa.
- Usar `WithBase64Data()` solo en endpoints que realmente lo necesiten.
- Mantener consistencia: o body plano o `data`, salvo que el contrato histórico del servicio ya obligue a soportar ambos.

## Limitaciones actuales a tener en cuenta

Estas limitaciones son importantes si la documentación se va a usar como guía para varios backends:

- `BodyParser` hace `panic` si `target` no es puntero.
- El error de campo faltante informa el nombre del campo Go, no la key JSON.
- `required` considera vacío solo a `string`.
- No hay validación profunda de estructuras anidadas.
- No hay validación declarativa adicional más allá de `required`.
- Si llega `data`, el parser ignora el resto del body para poblar el struct.

## Cuándo conviene complementar con validaciones propias

Conviene sumar otra capa de validación cuando necesites:

- Email válido.
- Longitud mínima o máxima.
- Enumeraciones.
- Reglas entre campos.
- Validación de arrays o estructuras internas.
- Rechazar explícitamente `0` o `false` como inválidos.

Ejemplo:

```go
if req.Seats <= 0 {
    return errors.New("seats must be greater than zero")
}
```

## Resumen operativo

Para usar bien las validaciones de estructuras en `uker` con `BodyParser`:

1. Pasa siempre un puntero a struct.
2. Define tags `json`.
3. Marca obligatorios con `uker:"required"`.
4. Usa `WithBase64Data()` solo si el payload viene dentro de `data` en base64.
5. Asume que la validación automática cubre presencia y strings no vacíos, no reglas de negocio.

Con ese criterio, `BodyParser` funciona bien como primera barrera de validación de entrada y mantiene el handler limpio, pero no reemplaza una capa de validación de dominio más completa.
