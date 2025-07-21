// Simple service worker to prevent 404 errors
// This is a minimal service worker that doesn't do anything special

self.addEventListener('install', function(event) {
  // Skip waiting to activate immediately
  self.skipWaiting();
});

self.addEventListener('activate', function(event) {
  // Claim all clients immediately
  event.waitUntil(self.clients.claim());
});

// Optional: Add basic fetch handling if needed
self.addEventListener('fetch', function(event) {
  // Let the browser handle all fetch requests normally
  // This is just here to prevent any potential issues
});
