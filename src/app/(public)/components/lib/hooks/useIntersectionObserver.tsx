import { useState, useEffect } from 'react';

let lastTimeout: NodeJS.Timeout;
export const useIntersectionObserver = () => {
    const [intersectionObserver, setIntersectionObserver] = useState<IntersectionObserver | null>(null);
    useEffect(() => {
        if (typeof window !== 'undefined' && typeof document !== 'undefined') {
            const handleIntersectionObserver: (entries: IntersectionObserverEntry[]) => void = (entries) => {
                entries.forEach(async (entry) => {
                    const allLinks = [...document.querySelectorAll('.links a')];
                    const currentLink = allLinks.find(
                        (link) => link.textContent === entry.target.querySelector('h1')?.textContent
                    );
                    if (entry.isIntersecting) {
                        if (lastTimeout) clearTimeout(lastTimeout);
                        lastTimeout = setTimeout(function () {
                            currentLink?.classList.toggle('active', true);
                        }, 200);
                    }
                    currentLink?.classList.toggle('active', false);
                });
            };

            setIntersectionObserver(
                new IntersectionObserver(handleIntersectionObserver, {
                    root: null,
                    threshold: 0.15,
                })
            );
        }
    }, []);
    return { intersectionObserver };
};
