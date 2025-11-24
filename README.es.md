# HomeKitGenQRCode

**Idioma / Language:** [Español](README.es.md) | [English](README.md)

Una aplicación en Go para generar etiquetas con códigos QR de HomeKit con información del dispositivo.

![WEB](HomeKitGenQRCode.png)

![Ejemplo de Etiqueta con Código QR de HomeKit](ejemplo.png)

## Descripción

¡Genera fácilmente etiquetas con códigos QR de HomeKit para tus accesorios ESP32! Esta herramienta en Go crea automáticamente etiquetas listas para imprimir que incluyen:

- **Código de configuración de HomeKit** (con código QR) - Escanea con iPhone para emparejar tu dispositivo
- **Código de dispositivo** - Identificador único para tu accesorio
- **Dirección MAC** - Identificador de red con código de barras
- **Número de serie** - Serie única con código de barras
- **CSN** (Número de serie del componente) - Identificador adicional con código de barras
- **Diseño ordenado y estético** - Apariencia profesional que coincide con los estándares de HomeKit de Apple

Perfecto para etiquetar profesionalmente tus proyectos DIY de HomeKit construidos con [HomeSpan](https://github.com/HomeSpan/HomeSpan/).

Esta herramienta fue creada para el proyecto [HomeSpan](https://github.com/HomeSpan/HomeSpan/), tomando inspiración de la implementación original en Python de [AchimPieters/esp32-homekit-qrcode](https://github.com/AchimPieters/esp32-homekit-qrcode), pero reescrita en Go para mejor rendimiento, distribución más fácil y compatibilidad multiplataforma mejorada.

## Características

- Genera etiquetas completas con códigos QR de HomeKit con toda la información requerida
- Soporte para todas las categorías de dispositivos HomeKit (Luz, Interruptor, Termostato, etc.)
- Generación automática de códigos de configuración, IDs de configuración y direcciones MAC
- Formato de etiqueta profesional que coincide con los estándares de HomeKit de Apple
- Códigos QR de alta calidad optimizados para escaneo
- Generación de códigos de barras para direcciones MAC, números de serie y CSNs
- Interfaz de línea de comandos con múltiples subcomandos
- Ejecutable binario único - no se requieren dependencias en tiempo de ejecución

## Instalación

### Desde el código fuente

```bash
git clone https://github.com/lordbasex/HomeKitGenQRCode.git
cd HomeKitGenQRCode
go build ./cmd/homekitgenqrcode
```

### Usando Go Install

```bash
go install github.com/lordbasex/HomeKitGenQRCode/cmd/homekitgenqrcode@latest
```

### Versión Web (Sin Instalación Requerida)

¿Prefieres usar la herramienta directamente en tu navegador? Otra opción es usar nuestra versión web alojada desde CNSoluciones - no necesitas descargar nada, solo genera y listo!

🌐 **[Usar Versión Web →](https://cnsoluciones.com/esp32-homekit-qrcode-generator-web/)**

La versión web proporciona la misma funcionalidad que la herramienta CLI, ejecutándose completamente en tu navegador sin necesidad de instalación ni descargas. ¡Perfecto para generar etiquetas rápidamente!

*Un regalo de Federico Pereira de navidad.*

## Uso

### Inicio rápido (Recomendado)

Genera una etiqueta con código QR con código de configuración generado automáticamente:

```bash
homekitgenqrcode code -c 5 -o ejemplo.png
```

### Generar con todos los parámetros

```bash
homekitgenqrcode generate --category 5 --password "613-80-755" --setup-id "ABCD" --mac "AABBCCDDEEFF" --output ejemplo.png
```

### Listar categorías disponibles

```bash
homekitgenqrcode list-categories
```

## Comandos

### `code` - Auto-generar código de configuración (Más fácil)

Genera automáticamente código de configuración, ID de configuración y dirección MAC:

```bash
homekitgenqrcode code -c <categoría> -o <salida.png>
```

Opciones:
- `-c, --category`: ID de categoría HomeKit (requerido)
- `-o, --output`: Ruta del archivo de imagen de salida (requerido)
- `-s, --setup-id`: ID de configuración personalizado (opcional, se genera automáticamente si no se proporciona)
- `-m, --mac`: Dirección MAC personalizada (opcional, se genera automáticamente si no se proporciona)

### `generate` - Generación manual

Genera con todos los parámetros especificados manualmente:

```bash
homekitgenqrcode generate -c <categoría> -p <contraseña> -s <setup-id> -m <mac> -o <salida.png>
```

Opciones:
- `-c, --category`: ID de categoría HomeKit (requerido)
- `-p, --password`: Contraseña de configuración en formato XXX-XX-XXX (requerido)
- `-s, --setup-id`: ID de configuración: 4 caracteres alfanuméricos (0-9, A-Z) (requerido)
- `-m, --mac`: Dirección MAC: 12 caracteres hexadecimales (requerido)
- `-o, --output`: Ruta del archivo de imagen de salida (requerido)

### `list-categories` - Listar categorías disponibles

Muestra todas las categorías de dispositivos HomeKit disponibles:

```bash
homekitgenqrcode list-categories
```

## Categorías de HomeKit

La siguiente tabla lista todas las categorías de dispositivos HomeKit soportadas con sus IDs:

| ID | Nombre de Categoría |
|----|---------------------|
| 1 | Otro |
| 2 | Puente |
| 3 | Ventilador |
| 4 | Abridor de puerta de garaje |
| 5 | Luz |
| 6 | Cerradura |
| 7 | Enchufe |
| 8 | Interruptor |
| 9 | Termostato |
| 10 | Sensor |
| 11 | Sistema de seguridad |
| 12 | Puerta |
| 13 | Ventana |
| 14 | Cobertura de ventana |
| 15 | Interruptor programable |
| 16 | Extensor de rango |
| 17 | Cámara IP |
| 18 | Timbre de video |
| 19 | Purificador de aire |
| 20 | Calentador |
| 21 | Aire acondicionado |
| 22 | Humidificador |
| 23 | Deshumidificador |
| 24 | Apple TV |
| 26 | Altavoz |
| 27 | Airport |
| 28 | Aspersor |
| 29 | Grifo |
| 30 | Cabeza de ducha |
| 31 | Televisión |
| 32 | Control remoto objetivo |

**Nota:** El ID de categoría 25 no está definido en la especificación de HomeKit.

## Ejemplos

```bash
# Generar con valores completamente automáticos
homekitgenqrcode code -c 5 -o ejemplo.png

# Generar con ID de configuración y MAC personalizados
homekitgenqrcode code -c 5 -o ejemplo.png -s ABCD -m AABBCCDDEEFF

# Generar en un directorio específico (se creará automáticamente)
homekitgenqrcode code -c 5 -o salida/ejemplo.png

# Usando flags largos
homekitgenqrcode generate --category 5 --password "613-80-755" --setup-id "ABCD" --mac "AABBCCDDEEFF" --output ejemplo.png
```

## Cómo funciona

1. **Carga la plantilla de etiqueta** (`assets/qrcode_ext.png`)
2. **Genera un código de configuración de HomeKit** (formato: XXX-XX-XXX) o usa uno proporcionado
3. **Crea un código QR** siguiendo los estándares de HomeKit de Apple con corrección de errores adecuada
4. **Genera información del dispositivo**:
   - Código de dispositivo (formato basado en categoría)
   - Dirección MAC (12 caracteres hexadecimales)
   - Número de serie (patrón alfanumérico único)
   - CSN (Número de serie del componente)
5. **Posiciona todos los elementos** estéticamente en la plantilla
6. **Exporta una etiqueta terminada** como archivo PNG de alta resolución (300 DPI)

Cada ejecución crea una etiqueta única con:
- Un código de configuración de HomeKit válido
- Número de serie y CSN únicos
- Códigos de barras generados automáticamente (formato Code 39)
- Un código QR siguiendo los estándares de HomeKit de Apple

## Impresión

El archivo PNG generado está listo para imprimir en etiquetas adhesivas blancas. La salida está optimizada a 300 DPI para impresión de alta calidad.

**Para mejores resultados:**
- Usa una impresora láser o una impresora de inyección de tinta de alta resolución
- Imprime en papel de etiquetas adhesivas blanco
- Asegúrate de que la configuración de la impresora coincida con el tamaño de la etiqueta

## Requisitos

- Go 1.24.0 o posterior (solo necesario para compilar desde el código fuente)
- Carpeta de assets con fuentes y plantilla de imagen requeridas (incluidas en el repositorio):
  - `SF-Pro-Text-Regular.otf` - Fuente de texto principal
  - `barcode39.ttf` - Fuente de código de barras
  - `qrcode_ext.png` - Plantilla de etiqueta

**Nota:** ¡Cuando uses el binario precompilado, no se requieren dependencias adicionales!

## Licencia

Este proyecto es de código abierto y está disponible para su uso.

## Contribuciones

¡Las contribuciones son bienvenidas! Por favor, siéntete libre de enviar un Pull Request.

## Créditos

Este proyecto fue creado para [HomeSpan](https://github.com/HomeSpan/HomeSpan/) - una biblioteca robusta de Arduino para crear dispositivos HomeKit basados en ESP32.

La idea y concepto original provienen de [AchimPieters/esp32-homekit-qrcode](https://github.com/AchimPieters/esp32-homekit-qrcode), un generador de códigos QR basado en Python. Esta implementación en Go proporciona:

- **Mejor rendimiento**: Los binarios compilados en Go son más rápidos y eficientes que los scripts interpretados en Python
- **Distribución más fácil**: Ejecutable binario único, no se requiere tiempo de ejecución de Python ni dependencias
- **Multiplataforma**: Funciona perfectamente en Windows, macOS y Linux sin configuración adicional
- **CLI mejorado**: Interfaz de línea de comandos moderna usando Cobra con soporte de autocompletado

## Integración

Esta herramienta está diseñada para funcionar perfectamente con [HomeSpan](https://github.com/HomeSpan/HomeSpan/), una biblioteca robusta de Arduino para crear dispositivos HomeKit basados en ESP32.

**Flujo de trabajo típico:**
1. Genera tu etiqueta con código QR de HomeKit usando esta herramienta
2. Imprime la etiqueta y adjúntala a tu dispositivo ESP32
3. Usa el código de configuración y la información en tu sketch de HomeSpan
4. Empareja tu dispositivo con HomeKit usando el código QR

## Proyectos relacionados

- [HomeSpan](https://github.com/HomeSpan/HomeSpan/) - Biblioteca HomeKit para Arduino-ESP32
- [esp32-homekit-qrcode](https://github.com/AchimPieters/esp32-homekit-qrcode) - Implementación original en Python

## Autor

**Federico Pereira**

Creado con ❤️ por Federico Pereira. Este proyecto fue desarrollado para hacer la generación de etiquetas con códigos QR de HomeKit más fácil y accesible para la comunidad DIY de HomeKit.

Para preguntas, sugerencias o contribuciones, por favor visita el [repositorio de GitHub](https://github.com/lordbasex/HomeKitGenQRCode).

