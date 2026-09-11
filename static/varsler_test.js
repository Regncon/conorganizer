const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const source = fs.readFileSync(path.join(__dirname, 'varsler.js'), 'utf8');

class FakeElement {
    constructor(href = '') {
        this.dataset = {};
        this.hidden = true;
        this.disabled = false;
        this.href = href;
        this.listeners = {};
    }

    addEventListener(name, handler) {
        this.listeners[name] = handler;
    }

    setAttribute(name, value) {
        this[name] = value;
    }

    async dispatch(name, event = {}) {
        return this.listeners[name]?.(event);
    }
}

function notificationPage({subscribed = true, configEnabled = true, permission = 'granted', deleteFails = false} = {}) {
    const status = new FakeElement();
    const help = new FakeElement();
    const enableButton = new FakeElement();
    const disableButton = new FakeElement();
    const logoutLink = new FakeElement('https://regncon.example/auth/logout');
    const settings = new FakeElement();
    settings.dataset.varslerConfigUrl = '/api/varsler/config';
    settings.querySelector = selector => ({
        '[data-varsler-status]': status,
        '[data-varsler-help]': help,
        '[data-varsler-enable]': enableButton,
        '[data-varsler-disable]': disableButton,
    })[selector] ?? null;

    const calls = [];
    const order = [];
    let unsubscribeCalls = 0;
    const subscription = {
        endpoint: 'https://fcm.googleapis.com/fcm/send/device',
        unsubscribe: async () => {
            unsubscribeCalls++;
            return true;
        },
        toJSON: () => ({endpoint: 'https://fcm.googleapis.com/fcm/send/device', keys: {p256dh: 'key', auth: 'auth'}}),
    };
    const registration = {
        pushManager: {
            getSubscription: async () => subscription,
            subscribe: async () => {
                order.push('subscribe');
                return subscription;
            },
        },
    };
    const serviceWorker = {
        register: async () => registration,
        getRegistration: async () => registration,
    };
    const fetch = async (url, options = {}) => {
        calls.push({url, options});
        if (url === '/api/varsler/config') {
            return response({enabled: configEnabled, publicKey: configEnabled ? 'AQ' : ''});
        }
        if (options.method === 'DELETE') {
            if (deleteFails) throw new Error('offline');
            return response(null, 204);
        }
        if (options.method === 'POST') return response(null, 204);
        return response({subscribed});
    };
    const assignedLocations = [];
    const window = {
        atob: value => Buffer.from(value, 'base64').toString('binary'),
        matchMedia: () => ({matches: false}),
        Notification: {
            permission,
            requestPermission: async () => {
                order.push('permission');
                return 'granted';
            },
        },
        PushManager: function () {},
        location: {assign: value => assignedLocations.push(value)},
    };
    const document = {
        querySelector: selector => {
            if (selector === '[data-varsler-settings]') return settings;
            if (selector === '[data-varsler-logout]') return logoutLink;
            return null;
        },
    };
    const navigator = {userAgent: '', platform: '', maxTouchPoints: 0, serviceWorker};
    const context = vm.createContext({window, document, navigator, Notification: window.Notification, PushManager: window.PushManager, Uint8Array, Buffer, fetch});
    vm.runInContext(source, context);

    return {
        calls,
        enableButton,
        disableButton,
        logoutLink,
        assignedLocations,
        order,
        unsubscribeCalls: () => unsubscribeCalls,
    };
}

function response(body, status = 200) {
    return {
        ok: status >= 200 && status < 300,
        status,
        json: async () => body,
    };
}

const settle = () => new Promise(resolve => setImmediate(resolve));

test('disabling sends the endpoint in the authenticated JSON request before unsubscribing locally', async () => {
    const page = notificationPage();
    await settle();

    await page.disableButton.dispatch('click');

    const deletion = page.calls.find(call => call.options.method === 'DELETE');
    assert.equal(deletion.url, '/api/varsler/subscriptions');
    assert.equal(deletion.options.headers['Content-Type'], 'application/json');
    assert.deepEqual(JSON.parse(deletion.options.body), {endpoint: 'https://fcm.googleapis.com/fcm/send/device'});
    assert.equal(page.unsubscribeCalls(), 1);
});

test('a subscription not owned by the signed-in account is removed locally before enabling', async () => {
    const page = notificationPage({subscribed: false});

    await settle();

    assert.equal(page.unsubscribeCalls(), 1);
    assert.equal(page.enableButton.hidden, false);
});

test('logout removes the current device subscription and then continues to logout', async () => {
    const page = notificationPage();
    await settle();
    let prevented = false;

    await page.logoutLink.dispatch('click', {preventDefault: () => { prevented = true; }});

    assert.equal(prevented, true);
    assert.equal(page.unsubscribeCalls(), 1);
    assert.deepEqual(page.assignedLocations, ['https://regncon.example/auth/logout']);
});

test('logout still removes local push state and navigates when server cleanup fails', async () => {
    const page = notificationPage({deleteFails: true});
    await settle();

    await page.logoutLink.dispatch('click', {preventDefault: () => {}});

    assert.equal(page.unsubscribeCalls(), 1);
    assert.deepEqual(page.assignedLocations, ['https://regncon.example/auth/logout']);
});

test('disabled configuration does not register a worker or expose controls', async () => {
    const page = notificationPage({configEnabled: false});

    await settle();

    assert.equal(page.enableButton.hidden, true);
    assert.equal(page.disableButton.hidden, true);
    assert.equal(page.calls.length, 1);
});

test('notification permission is requested directly from the enable-button gesture before subscribing', async () => {
    const page = notificationPage({subscribed: false, permission: 'default'});
    await settle();
    assert.deepEqual(page.order, []);

    await page.enableButton.dispatch('click');

    assert.deepEqual(page.order, ['permission', 'subscribe']);
});
