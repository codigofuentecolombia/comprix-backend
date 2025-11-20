package utils

func GetBlockedResources() []string {
	return []string{
		// Bloquear todos los recursos
		"*.png", "*.jpg", "*.jpeg", "*.gif", "*.webp", "*.svg",
		"*.mp4", "*.avi", "*.mov", "*.webm", "*.mkv",
		"*.woff", "*.woff2", "*.ttf", "*.eot", "*.otf",

		// Bloquear descarga  de imagen jumbo
		"https://jumboargentina.vtexassets.com/arquivos/*",
		"https://jumboargentinaio.vtexassets.com/arquivos/*",

		"https://*.com/arquivos/*",

		// Bloquear anuncios
		"ads/*",
		"doubleclick.net",
		"googleadservices.com",
		"googlesyndication.com",
		"stats.g.doubleclick.net",

		"facebook.com",
		"connect.facebook.net",
		"analytics.tiktok.com",
		"analytics.google.com",
		"notifications-icommkt.com",
	}
}
