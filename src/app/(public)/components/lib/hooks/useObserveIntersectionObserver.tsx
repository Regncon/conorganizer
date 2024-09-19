import { useEffect } from 'react';
import { useIntersectionObserver } from './useIntersectionObserver';

export const useObserveIntersectionObserver = (ref: React.RefObject<HTMLDivElement>) => {
    const { intersectionObserver } = useIntersectionObserver();
    useEffect(() => {
        if (ref.current) {
            intersectionObserver?.observe(ref.current);
        }
        return () => {
            if (ref.current) {
                intersectionObserver?.unobserve(ref.current);
            }
        };
    }, [ref, ref.current]);
};
