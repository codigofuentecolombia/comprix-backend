package pages

import (
	"comprix/app/constants"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func (scrapper *ScrapPage) createCacheDir() {
	cacheDir := constants.LinuxChromeCacheFolder
	// Verificar el sistema operativo y establecer la ruta para la caché
	if runtime.GOOS == "windows" {
		cacheDir = constants.WindowsChromeCacheFolder
	}
	// Verificar si sirve
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			scrapper.log.Error(fmt.Sprintf("Error al crear el directorio de caché: %s", err.Error()))
		}
	}
}

func (scrapper *ScrapPage) removeCacheDir() {
	cacheDir := constants.LinuxChromeCacheFolder
	// Verificar el sistema operativo y establecer la ruta para la caché
	if runtime.GOOS == "windows" {
		cacheDir = constants.WindowsChromeCacheFolder
	}
	// Verificar si el directorio existe y eliminarlo
	if _, err := os.Stat(cacheDir); err == nil {
		err := os.RemoveAll(cacheDir) // Elimina el directorio y todo su contenido
		if err != nil {
			scrapper.log.Error(fmt.Sprintf("Error al eliminar el directorio de caché: %s", err.Error()))
		}
	}
}

func (scrapper *ScrapPage) removeChromeProcesses() {
	// Comprobar el sistema operativo y ejecutar el comando correspondiente
	if runtime.GOOS == "windows" {
		// En Windows, terminamos el proceso de Chrome utilizando "taskkill"
		if err := exec.Command("taskkill", "/F", "/IM", "chrome.exe").Run(); err != nil {
			scrapper.log.Error(fmt.Sprintf("Error al terminar los procesos de Chrome: %s", err.Error()))
		}
	} else if runtime.GOOS == "linux" {
		commands := [][]string{
			{"pkill", "-9", "chrome"},
			{"pkill", "-9", "chromium"},
		}
		// Iterar comandos
		for _, cmdArgs := range commands {
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			// Verificar si ocurrio error
			if err := cmd.Run(); err != nil {
				scrapper.log.Error(fmt.Sprintf("Error al ejecutar el comando %v: %s", cmdArgs, err.Error()))
			}
		}
	}
	// Esperar
	time.Sleep(15 * time.Second)
}
