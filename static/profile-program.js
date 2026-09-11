(() => {
    'use strict';

    let focusedTarget = null;

    const focusDeepLink = () => {
        const url = new URL(window.location.href);
        const pulje = url.searchParams.get('pulje');
        let fragment = '';
        try {
            fragment = decodeURIComponent(url.hash.slice(1));
        } catch {
            fragment = url.hash.slice(1);
        }

        const escapedPulje = pulje && window.CSS?.escape ? CSS.escape(pulje) : '';
        const group = escapedPulje ? document.querySelector(`details[data-program-group="${escapedPulje}"]`) : null;
        if (group instanceof HTMLDetailsElement) {
            group.open = true;
        }

        const target = (pulje && document.getElementById(`program-${pulje}`))
            || group
            || (fragment && document.getElementById(fragment));
        if (!target || target === focusedTarget) return target;

        target.scrollIntoView({block: 'start'});
        if (typeof target.focus === 'function') {
            target.setAttribute('tabindex', '-1');
            target.focus({preventScroll: true});
        }
        focusedTarget = target;
        return target;
    };

    const url = new URL(window.location.href);
    if (url.hash || url.searchParams.has('pulje')) {
        const observer = new MutationObserver(focusDeepLink);
        observer.observe(document.body, {childList: true, subtree: true});
        window.setTimeout(() => observer.disconnect(), 30000);
    }
    focusDeepLink();
})();
