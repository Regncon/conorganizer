import type { ReactNode } from 'react';

type Props = {
    children: ReactNode;
};

const Layout = async ({ children }: Props) => {
    return <>{children}</>;
};

export default Layout;
