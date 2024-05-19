package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Declaración de variables para almacenar los valores de las banderas de la línea de comandos
var (
	url         = flag.String("url", "", "URL to test")               // URL a probar
	concurrency = flag.Int("c", 1, "Concurrency level")               // Nivel de concurrencia (número de gorutinas)
	requests    = flag.Int("n", 1, "Number of requests")              // Número de peticiones a realizar por cada gorutina
	duration    = flag.Duration("d", 10*time.Second, "Test duration") // Duración del test
)

func main() {
	// Parsear las banderas de la línea de comandos
	flag.Parse()

	// Verificar si la URL es vacía; si es así, imprimir un mensaje de error y salir del programa
	if *url == "" {
		fmt.Println("Error: URL is required")
		return
	}

	// Crear un grupo de espera para sincronizar las gorutinas
	var wg sync.WaitGroup
	wg.Add(*concurrency)

	// Iniciar el tiempo de ejecución del test
	startTime := time.Now()

	// Bucle de concurrencia: crear un número de gorutinas igual al nivel de concurrencia especificado
	for i := 0; i < *concurrency; i++ {
		go func() {

			// Indicar que la gorutina ha finalizado cuando termine
			defer wg.Done()

			// Bucle para realizar el número de peticiones especificado
			for j := 0; j < *requests; j++ {

				// Crear una solicitud HTTP con el método GET y la URL especificada
				req, err := http.NewRequest("GET", *url, nil)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Realizar la solicitud utilizando el cliente HTTP predeterminado
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Cerrar el cuerpo de la respuesta
				defer resp.Body.Close()
			}
		}()
	}

	// Esperar a que todas las gorutinas terminen
	wg.Wait()

	// Calcular el tiempo transcurrido desde el inicio del test
	elapsedTime := time.Since(startTime)

	// Imprimir los resultados, incluyendo la tasa de peticiones por segundo
	fmt.Printf("Requests: %d, Concurrency: %d, Duration: %v, Req/Sec: %.2f\n", *requests, *concurrency, elapsedTime, float64(*requests)/elapsedTime.Seconds())
}
