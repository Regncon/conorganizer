const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const source = fs.readFileSync(path.join(__dirname, 'profile-program.js'), 'utf8');

test('pulje query opens its group and focuses the pulje card ahead of the page fragment', () => {
    class DetailsElement {}
    const root = target();
    const pulje = target();
    const group = Object.assign(new DetailsElement(), {open: false});
    const document = {
        body: {},
        querySelector: selector => selector === 'details[data-program-group="FredagKveld"]' ? group : null,
        getElementById: id => ({'mitt-program': root, 'program-FredagKveld': pulje})[id] ?? null,
    };
    const window = {
        location: {href: 'https://regncon.example/profile?pulje=FredagKveld#mitt-program'},
        CSS: {escape: value => value},
        setTimeout: () => {},
    };
    class MutationObserver {
        observe() {}
        disconnect() {}
    }

    vm.runInContext(source, vm.createContext({window, document, URL, CSS: window.CSS, HTMLDetailsElement: DetailsElement, MutationObserver}));

    assert.equal(group.open, true);
    assert.equal(pulje.scrolled, true);
    assert.equal(pulje.focused, true);
    assert.equal(root.scrolled, false);
});

function target() {
    return {
        scrolled: false,
        focused: false,
        scrollIntoView() { this.scrolled = true; },
        setAttribute() {},
        focus() { this.focused = true; },
    };
}
