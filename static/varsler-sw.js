const profileNotificationURL = (value) => {
    if (typeof value !== 'string' || !value) return null;
    try {
        const url = new URL(value, self.location.origin);
        if (url.origin !== self.location.origin) return null;
        if (url.username || url.password) return null;
        if (url.pathname !== '/profile' && url.pathname !== '/profile/') return null;
        return url.href;
    } catch {
        return null;
    }
};

self.addEventListener('push', event => {
    let payload = {};
    try {
        payload = event.data?.json() ?? {};
    } catch {
        payload = {body: event.data?.text() ?? ''};
    }
    const title = typeof payload.title === 'string' && payload.title ? payload.title : 'Regncon';
    const options = {
        body: typeof payload.body === 'string' ? payload.body : '',
        tag: typeof payload.tag === 'string' ? payload.tag : undefined,
        data: {url: profileNotificationURL(payload.data?.url ?? payload.url)},
        icon: '/static/favicon.ico',
    };
    event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', event => {
    event.notification.close();
    const destination = profileNotificationURL(event.notification.data?.url);
    if (!destination) return;
    event.waitUntil(clients.matchAll({type: 'window', includeUncontrolled: true}).then(windows => {
        const existing = windows.find(client => {
            try {
                const url = new URL(client.url);
                return url.origin === self.location.origin && (url.pathname === '/profile' || url.pathname === '/profile/');
            } catch {
                return false;
            }
        });
        if (existing) return existing.focus().then(() => existing.navigate(destination));
        return clients.openWindow(destination);
    }));
});
