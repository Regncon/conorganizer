import type { Route } from 'next';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';

export const useSafeRouterReplace = () => {
    const router = useRouter();
    const pathname = usePathname();
    const searchParams = useSearchParams();

    const setQuery = (queryParams: { key: string; value?: string }[]) => {
        const current = new URLSearchParams(Array.from(searchParams.entries()));

        queryParams.forEach(({ key, value }) => {
            if (!value) {
                current.delete(key);
            } else {
                current.set(key, value);
            }
        });

        const search = current.toString();
        const query = search ? `?${search}` : '';
        console.log(query);
        router.replace(`${pathname}${query}` as Route);
    };
    return { setQuery };
};

export default useSafeRouterReplace;
