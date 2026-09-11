const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const workerSource = fs.readFileSync(path.join(__dirname, 'varsler-sw.js'), 'utf8');

function loadWorker() {
    const handlers = {};
    const shownNotifications = [];
    const serviceWorker = {
        showNotification: async (title, options) => shownNotifications.push({title, options}),
    };
    const self = {location: {origin: 'https://regncon.example'}, registration: serviceWorker, addEventListener: (name, handler) => { handlers[name] = handler; }};
    const clients = {};
    const context = vm.createContext({self, URL, clients});
    vm.runInContext(workerSource, context);
    return {handlers, shownNotifications, self, clients};
}

test('push notifications render trusted profile destination data', async () => {
    const worker = loadWorker();
    const waits = [];
    await worker.handlers.push({
        data: {json: () => ({title: 'Programoppdatering', body: 'Nytt tidspunkt', tag: 'program', data: {url: '/profile?pulje=FredagKveld'}})},
        waitUntil: promise => waits.push(promise),
    });
    await Promise.all(waits);

    assert.equal(worker.shownNotifications.length, 1);
    assert.equal(worker.shownNotifications[0].title, 'Programoppdatering');
    assert.equal(worker.shownNotifications[0].options.body, 'Nytt tidspunkt');
    assert.equal(worker.shownNotifications[0].options.data.url, 'https://regncon.example/profile?pulje=FredagKveld');
});

test('notification clicks navigate an existing profile window and reject external URLs', async () => {
    const worker = loadWorker();
    const existing = {url: 'https://regncon.example/profile', focus: async () => {}, navigate: async url => { existing.navigatedTo = url; }};
    const originalMatchAll = worker.clients.matchAll;
    worker.clients.matchAll = async () => [existing];

    const waits = [];
    await worker.handlers.notificationclick({
        notification: {data: {url: '/profile#mitt-program'}, close: () => {}},
        waitUntil: promise => waits.push(promise),
    });
    await Promise.all(waits);
    assert.equal(existing.navigatedTo, 'https://regncon.example/profile#mitt-program');

    let rejectedWait = false;
    worker.handlers.notificationclick({
        notification: {data: {url: 'https://evil.example/profile'}, close: () => {}},
        waitUntil: () => { rejectedWait = true; },
    });
    assert.equal(rejectedWait, false);

	worker.handlers.notificationclick({
		notification: {data: {url: 'https://user:password@regncon.example/profile'}, close: () => {}},
		waitUntil: () => { rejectedWait = true; },
	});
	assert.equal(rejectedWait, false);
    worker.clients.matchAll = originalMatchAll;
});
