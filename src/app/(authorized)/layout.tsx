import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { redirect } from 'next/navigation';
import type { ReactNode } from 'react';

type Props = {
    children: ReactNode;
};

const Layout = async ({ children }: Props) => {
    const { auth } = await getAuthorizedAuth();
    if (auth) {
        return children;
    }
    redirect('/login');
};

export default Layout;
