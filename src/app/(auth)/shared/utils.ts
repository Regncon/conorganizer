import type { Route } from 'next';
import type { AppRouterInstance } from 'next/dist/shared/lib/app-router-context.shared-runtime';
import type { ComponentProps, FormEvent } from 'react';

export const emailRegExp = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+.[a-zA-Z]{2,4}$/;

export const updateSearchParamsWithEmail = async (
    e: FormEvent<HTMLFormElement>,
    router: AppRouterInstance,
    href: Route
) => {
    const { value, name } = e.target as HTMLInputElement;
    if (name === 'email') {
        router.replace(`${href}?email=${value}`);
    }
};
