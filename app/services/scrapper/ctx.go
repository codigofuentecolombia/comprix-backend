package service_scrapper

import (
	"comprix/app/fails"
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func (s *Service) InitContext(optTimeout ...time.Duration) (context.Context, context.CancelFunc) {
	// Crear tiempo de expiracion
	var timeout time.Duration = 20
	// Verificar si existe dato
	if len(optTimeout) > 0 {
		timeout = optTimeout[0]
	}
	// Allocator
	ctx, cancelAllocatorFn := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/XXXXX Safari/537.36"),
			chromedp.WindowSize(1920, 1080),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
		)...,
	)
	// Chromedp
	ctx, cancelChromedpFn := chromedp.NewContext(ctx)
	// Añadir timeout
	ctx, cancelTimeout := context.WithTimeout(ctx, timeout*time.Second)
	// Nueva función de cancelación que cierra ambas
	cancel := func() {
		cancelAllocatorFn()
		cancelChromedpFn()
		cancelTimeout() // Cancela el timeout
	}
	//Regresar data
	return ctx, cancel
}

func (s *Service) InitChromedp() (context.Context, ServiceChromedpCtxCancelFn, error) {
	var cancelFns ServiceChromedpCtxCancelFn

	ctx, cancelAllocatorFn := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/XXXXX Safari/537.36"),
			chromedp.WindowSize(1920, 1080),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
		)...,
	)
	if ctx == nil {
		return nil, ServiceChromedpCtxCancelFn{}, fails.Create("failed to create allocator", nil)
	}
	cancelFns.Allocator = cancelAllocatorFn

	ctx, cancelChromedpFn := chromedp.NewContext(ctx)
	if ctx == nil {
		cancelAllocatorFn()
		return nil, ServiceChromedpCtxCancelFn{}, fails.Create("failed to create chromedp context", nil)
	}
	cancelFns.Chromedp = cancelChromedpFn

	return ctx, cancelFns, nil
}

func (s *Service) CloseChromedpCtx() {
	// Verificar si tiene
	if s.ChromedpCtx != nil {
		s.ChromedpCtx.CancelFns.Allocator()
		s.ChromedpCtx.CancelFns.Chromedp()
	}
}

// Añadir timeout
// cacheDir := constants.LinuxChromeCacheFolder

// Verificar el sistema operativo y establecer la ruta para la caché
// if runtime.GOOS == "windows" {
// 	cacheDir = constants.WindowsChromeCacheFolder
// }
// Opciones optimizadas para Chrome
// opts := []func(*chromedp.ExecAllocator){
// chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),

// chromedp.Flag("headless", false),
// chromedp.Flag("incognito", true),
// chromedp.Flag("mute-audio", true),
// chromedp.Flag("disable-extensions", true),

// chromedp.Flag("disk-cache-dir", cacheDir), // Usa RAM para cache

// chromedp.Flag("enable-low-res-tiling", true),               // Reduce uso de CPU en el render
// chromedp.Flag("disable-background-timer-throttling", true), // Menos uso de CPU

// chromedp.Flag("enable-begin-frame-scheduling", true),       // Optimiza los frames renderizados
// chromedp.Flag("disable-background-networking", true),       // Evita tareas en segundo plano
// chromedp.Flag("disable-renderer-backgrounding", true),      // Evita que Chrome pause procesos
// }
// Crear contexto de Chrome con opciones
// allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
