(() => {
    'use strict';

    const settings = document.querySelector('[data-varsler-settings]');
    if (!settings) return;

    const status = settings.querySelector('[data-varsler-status]');
    const help = settings.querySelector('[data-varsler-help]');
    const enableButton = settings.querySelector('[data-varsler-enable]');
    const disableButton = settings.querySelector('[data-varsler-disable]');
    const logoutLink = document.querySelector('[data-varsler-logout]');
    const configURL = settings.dataset.varslerConfigUrl || '/api/varsler/config';

    const setStatus = (message, {kind = '', helpText = ''} = {}) => {
        if (status) {
            status.textContent = message;
            status.dataset.state = kind;
        }
        if (help) {
            help.textContent = helpText;
            help.hidden = !helpText;
        }
    };

    const hideControls = () => {
        if (enableButton) enableButton.hidden = true;
        if (disableButton) disableButton.hidden = true;
    };

    const showEnable = () => {
        if (enableButton) enableButton.hidden = false;
        if (disableButton) disableButton.hidden = true;
    };

    const showDisable = () => {
        if (enableButton) enableButton.hidden = true;
        if (disableButton) disableButton.hidden = false;
    };

    const supported = () => 'Notification' in window && 'serviceWorker' in navigator && 'PushManager' in window;

    const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent) || (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1);
    const isStandalone = window.matchMedia?.('(display-mode: standalone)').matches || navigator.standalone === true;

    const toApplicationServerKey = (key) => {
        const padding = '='.repeat((4 - key.length % 4) % 4);
        const base64 = (key + padding).replace(/-/g, '+').replace(/_/g, '/');
        const raw = window.atob(base64);
        return Uint8Array.from(raw, character => character.charCodeAt(0));
    };

    const subscriptionURL = (endpoint) => `/api/varsler/subscriptions?endpoint=${encodeURIComponent(endpoint)}`;

    const requestJSON = async (url, options = {}) => {
        const response = await fetch(url, {credentials: 'same-origin', ...options});
        if (!response.ok) throw new Error(`Varsler svarte med ${response.status}`);
        return response.status === 204 ? null : response.json();
    };

    const setEnabled = () => {
        settings.dataset.varslerState = 'enabled';
        showDisable();
        setStatus('Varsler er på denne enheten.', {kind: 'enabled'});
    };

    const inspectSubscription = async () => {
        const registration = await navigator.serviceWorker.register('/varsler-sw.js', {scope: '/'});
        const subscription = await registration.pushManager.getSubscription();
        if (!subscription) {
            showEnable();
            setStatus('Du kan få beskjed når festivalprogrammet endres.');
            return {registration, subscription: null};
        }
        const result = await requestJSON(subscriptionURL(subscription.endpoint));
        if (result?.subscribed === true) {
            setEnabled();
            return {registration, subscription};
        }
        await subscription.unsubscribe();
        showEnable();
        setStatus('Du kan få beskjed når festivalprogrammet endres.');
        return {registration, subscription: null};
    };

    const enable = async () => {
        enableButton.disabled = true;
        setStatus('Ber om tillatelse til varsler …');
        try {
            const permission = await Notification.requestPermission();
            if (permission !== 'granted') {
                showEnable();
                if (permission === 'denied') {
                    setStatus('Varsler er blokkert i nettleseren.', {kind: 'denied', helpText: 'Åpne nettleserens innstillinger for dette nettstedet og tillat varsler.'});
                } else {
                    setStatus('Varsler ble ikke slått på.');
                }
                return;
            }

            const config = await requestJSON(configURL);
            if (config?.enabled !== true || typeof config.publicKey !== 'string' || !config.publicKey) {
                hideControls();
                setStatus('Varsler er ikke tilgjengelige akkurat nå.');
                return;
            }
            const registration = await navigator.serviceWorker.register('/varsler-sw.js', {scope: '/'});
            const subscription = await registration.pushManager.subscribe({userVisibleOnly: true, applicationServerKey: toApplicationServerKey(config.publicKey)});
            await requestJSON('/api/varsler/subscriptions', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(subscription.toJSON()),
            });
            setEnabled();
        } catch (error) {
            showEnable();
            setStatus('Vi klarte ikke å slå på varsler. Prøv igjen senere.', {kind: 'error'});
        } finally {
            enableButton.disabled = false;
        }
    };

    const disable = async () => {
        disableButton.disabled = true;
        try {
            const registration = await navigator.serviceWorker.getRegistration('/');
            const subscription = await registration?.pushManager.getSubscription();
            if (subscription) {
                await requestJSON('/api/varsler/subscriptions', {
                    method: 'DELETE',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({endpoint: subscription.endpoint}),
                });
                await subscription.unsubscribe();
            }
            showEnable();
            setStatus('Varsler er slått av på denne enheten.');
        } catch (error) {
            setStatus('Vi klarte ikke å slå av varsler. Prøv igjen senere.', {kind: 'error'});
        } finally {
            disableButton.disabled = false;
        }
    };

    let logoutInProgress = false;
    const logout = async (event) => {
        if (logoutInProgress) return;
        logoutInProgress = true;
        event.preventDefault();
        logoutLink.setAttribute('aria-disabled', 'true');

        try {
            const registration = await navigator.serviceWorker.getRegistration('/');
            const subscription = await registration?.pushManager.getSubscription();
            if (subscription) {
                try {
                    await requestJSON('/api/varsler/subscriptions', {
                        method: 'DELETE',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({endpoint: subscription.endpoint}),
                        keepalive: true,
                    });
                } catch {
                }
                try {
                    await subscription.unsubscribe();
                } catch {
                }
            }
        } finally {
            window.location.assign(logoutLink.href);
        }
    };

    enableButton?.addEventListener('click', enable);
    disableButton?.addEventListener('click', disable);
    logoutLink?.addEventListener('click', logout);

    if (!supported()) {
        hideControls();
        setStatus('Denne nettleseren støtter ikke varsler.');
    } else if (isIOS && !isStandalone) {
        hideControls();
        setStatus('Legg Regncon til på Hjem-skjerm for å bruke varsler.', {helpText: 'Trykk Del, velg «Legg til på Hjem-skjerm», og åpne Regncon derfra.'});
    } else {
        requestJSON(configURL).then(config => {
            if (config?.enabled !== true || typeof config.publicKey !== 'string' || !config.publicKey) {
                hideControls();
                setStatus('Varsler er ikke tilgjengelige akkurat nå.');
                return;
            }
            if (Notification.permission === 'denied') {
                showEnable();
                setStatus('Varsler er blokkert i nettleseren.', {kind: 'denied', helpText: 'Åpne nettleserens innstillinger for dette nettstedet og tillat varsler.'});
                return;
            }
            return inspectSubscription();
        }).catch(() => {
            hideControls();
            setStatus('Varsler er ikke tilgjengelige akkurat nå.');
        });
    }
})();
