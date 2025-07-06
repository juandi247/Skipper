Diseño Conceptual de skipper-tunnel
Este documento describe la arquitectura y el ciclo de vida del cliente skipper-tunnel, diseñado como una Máquina de Estados Finita (FSM) para garantizar la máxima robustez, resiliencia y un comportamiento predecible.
Filosofía de Diseño
El túnel es un agente persistente y resiliente. Su objetivo es mantener una conexión estable con el proxy Skipper y reintentar de forma inteligente ante fallos de red. Todo el sistema está diseñado para ser cancelable y apagarse de forma ordenada en cualquier momento.
El Ciclo de Vida: Una Máquina de Estados
El túnel opera a través de una serie de estados discretos. Solo puede estar en un estado a la vez, y cada estado es responsable de una tarea específica antes de hacer la transición al siguiente. Un contexto de Go se pasa a través de todos los estados para manejar la cancelación global (por ejemplo, cuando el usuario presiona Ctrl+C).
Fases y Estados del Túnel
1. Fase de Inicialización (State: Initializing)
Objetivo: Validar la configuración y el entorno local antes de cualquier operación de red.
Pasos:
Leer Configuración: Parsea los flags de la línea de comandos (-subdomain, -port, etc.) y las variables de entorno (-token).
Validar Sintaxis: Comprueba que el subdominio y el puerto son válidos. Si no, el programa termina con un error claro.
Verificar Servicio Local: Intenta establecer una conexión TCP a localhost en el puerto especificado. Esto confirma que hay un servicio activo escuchando. Si falla, el programa termina, informando al usuario que debe iniciar su aplicación local primero.
Transición: Si todo es correcto, pasa al estado de Connecting.
2. Fase de Conexión y Autenticación (State: Connecting & Authenticating)
Objetivo: Establecer una conexión segura y autenticada con el proxy Skipper.
Pasos:
Conectar al Proxy: Intenta establecer una conexión TCP con el servidor de control del proxy. Esta operación respeta la señal de cancelación del contexto.
Enviar Solicitud de Túnel: Una vez conectado, envía un Frame de Control_AuthRequest con el token y el subdominio.
Esperar Reconocimiento (ACK): Espera (con un timeout) una respuesta del proxy.
Transición:
Si se recibe un ACK (success: true): La conexión es válida. Se muestra un mensaje de éxito al usuario (ej. "Túnel activo en mi-app.skipper.lat") y se transiciona al estado Running.
Si se recibe un NACK (success: false): Se muestra el error específico al usuario ("Subdominio ya en uso", "Token inválido") y se transiciona a un estado de error terminal.
Si la conexión falla o hay un timeout: Se transiciona al estado Reconnecting.
3. Fase Operacional (State: Running)
Objetivo: El túnel está activo, procesando datos y manteniéndose saludable.
Acciones (Lanzadas como Gorutinas Concurrentes): Al entrar en este estado, se inician varias tareas de fondo, todas ellas "escuchando" la señal de cancelación del contexto:
El Reactor: Una gorutina dedicada a leer eficientemente los Frames de peticiones HTTP que llegan del proxy.
El Worker Pool: Un número fijo de gorutinas que reciben las peticiones del Reactor. Cada worker es responsable de:
Reconstruir la petición HTTP.
Enviarla al servicio localhost.
Recibir la respuesta de localhost.
Empaquetarla en un Frame de respuesta y enviarla de vuelta al proxy.
El Heartbeat (Ping): Una gorutina que envía periódicamente Frames de Control_Ping al proxy para mantener la conexión viva y detectar desconexiones silenciosas.
4. Fase de Resiliencia (State: Reconnecting)
Objetivo: Manejar fallos de red sin terminar el programa.
Acción: Si la conexión con el proxy se pierde en cualquier momento durante la fase Running o Connecting, la máquina de estados entra aquí. Espera un tiempo (usando una estrategia de backoff exponencial) antes de volver a intentar la conexión.
Transición: Después de la espera, transiciona de nuevo al estado Connecting.
Robustez y Apagado (Graceful Shutdown)
Cancelación Global: Todo el ciclo de vida de la FSM es consciente del contexto global. Si el usuario presiona Ctrl+C, el contexto se cancela.
Efecto de la Cancelación:
Si el túnel está intentando conectarse, la operación se abortará inmediatamente.
Si el túnel está en estado Running, todas las gorutinas (Reactor, Workers, Heartbeat) detectarán la cancelación y terminarán limpiamente.
La máquina de estados transicionará a un estado de apagado para cerrar los recursos restantes y salir.
Esto garantiza que el programa nunca se quede en un estado "zombi" y siempre termine de forma predecible.
Funcionalidades Adicionales (A Futuro)
Dashboard Local: Se podría levantar un pequeño servidor web local (ej. en http://localhost:4040) que muestre en tiempo real el estado de la conexión, el log de peticiones procesadas y estadísticas de tráfico, proporcionando visibilidad al usuario.