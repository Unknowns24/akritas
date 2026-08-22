# Flutter Mobile Permissions Policy

Esta policy aplica a apps Flutter cuando una feature requiere permisos del dispositivo o cambios de configuración nativa.

## Objetivo

Evitar permisos innecesarios, configuraciones incompletas en Android/iOS y errores de publicación por falta de descripciones de privacidad.

Toda feature que use capacidades sensibles del dispositivo debe declarar, justificar y manejar correctamente sus permisos.

## Permisos comunes

Revisar esta policy si la feature usa:

- ubicación;
- cámara;
- fotos, galería o media library;
- notificaciones push/locales;
- micrófono;
- contactos;
- archivos/storage;
- bluetooth;
- biometría;
- background location;
- background fetch/tasks;
- sensores del dispositivo.

## Reglas generales

- No pedir permisos antes de que exista una intención clara del usuario.
- No pedir permisos en el arranque de la app salvo que sea indispensable para el funcionamiento principal.
- Explicar en UI por qué se necesita el permiso antes de invocar el request nativo cuando el contexto no sea obvio.
- Pedir el permiso más limitado posible.
- No agregar permisos amplios si una alternativa más acotada resuelve el caso de uso.
- Si el usuario deniega el permiso, mostrar un estado recuperable.
- Si el permiso queda denegado permanentemente, ofrecer navegación a settings cuando tenga sentido.
- No bloquear toda la app si solo una feature necesita el permiso.
- No simular datos reales cuando el permiso fue rechazado.

## Android

Cuando se agregue un permiso, revisar:

```txt
android/app/src/main/AndroidManifest.xml
android/app/build.gradle
android/build.gradle
```

Reglas:

- Agregar solo los permisos requeridos por la feature.
- Revisar cambios de permisos por versión de Android.
- Para Android 13+, revisar permisos de notificaciones y media según corresponda.
- Para ubicación en background, justificar explícitamente la necesidad y evitarla si la app puede funcionar con foreground location.
- Para cámara o galería, preferir permisos modernos y acotados cuando el plugin lo permita.
- Si se usan services en background, declarar el service y su foreground service type cuando corresponda.

## iOS

Cuando se agregue un permiso, revisar:

```txt
ios/Runner/Info.plist
ios/Podfile
```

Reglas:

- Toda API sensible debe tener su usage description correspondiente en `Info.plist`.
- Los textos de propósito deben ser claros para el usuario final.
- No usar textos genéricos como “Necesitamos este permiso”.
- Si una librería referencia una API sensible aunque la app la use indirectamente, agregar la purpose string requerida.
- Revisar permisos de ubicación `WhenInUse` vs `Always` y usar el más limitado posible.
- No habilitar background modes salvo que la feature lo necesite realmente.

Ejemplos de keys frecuentes:

```txt
NSCameraUsageDescription
NSPhotoLibraryUsageDescription
NSPhotoLibraryAddUsageDescription
NSLocationWhenInUseUsageDescription
NSLocationAlwaysAndWhenInUseUsageDescription
NSMicrophoneUsageDescription
NSUserTrackingUsageDescription
NSContactsUsageDescription
NSFaceIDUsageDescription
```

## Manejo en UI

Toda pantalla que dependa de permisos debe contemplar:

- permiso no solicitado;
- permiso concedido;
- permiso denegado;
- permiso denegado permanentemente;
- servicio apagado, por ejemplo ubicación desactivada;
- error técnico del plugin o plataforma.

Reglas:

- Mostrar mensajes en español y accionables.
- Incluir botón de reintento cuando corresponda.
- Incluir acceso a configuración del sistema cuando corresponda.
- Mantener la UI existente salvo que sea necesario agregar estados.
- No mezclar lógica de permisos compleja dentro del widget; moverla a `application` o a un service/gateway adecuado.

## Arquitectura

Los permisos deben respetar la arquitectura del proyecto.

Permitido:

```txt
presentation → application → domain
application → gateway contract
data/infrastructure → plugin concreto
```

Prohibido:

```txt
presentation → plugin nativo directamente, si la lógica es compleja o reusable
presentation → data
presentation → ApiClient para resolver permisos
```

Si la integración es simple y puramente visual, puede vivir cerca de la pantalla. Si se reutiliza o tiene reglas, debe encapsularse.

## Testing

Cuando una feature dependa de permisos, el plan TDD debe cubrir al menos:

- permiso concedido;
- permiso denegado;
- permiso denegado permanentemente;
- error del gateway/plugin;
- acción esperada de retry o settings;
- que no se llame al backend si falta un permiso indispensable.

Usar fakes/mocks de gateways de permisos. No testear contra plugins reales en unit tests.

## Checklist

Antes de cerrar una feature con permisos:

- [ ] Se pidió el permiso solo cuando había intención clara del usuario.
- [ ] AndroidManifest fue actualizado si correspondía.
- [ ] Info.plist fue actualizado si correspondía.
- [ ] Los textos de propósito son claros y específicos.
- [ ] La UI maneja denegado y denegado permanentemente.
- [ ] La feature tiene retry o alternativa cuando corresponde.
- [ ] No se agregaron permisos innecesarios.
- [ ] No se habilitaron background modes sin justificación.
- [ ] El plan TDD cubre estados de permisos.
