'use client';
import { useSetCustomCssVariable } from '$lib/hooks/useSetCustomCssVariable';
import MainAppBarContent from './MainAppBarContent';

type Props = {
    isLoggedIn: boolean;
    admin: boolean;
};

const MainAppBarDesktop = ({ isLoggedIn, admin }: Props) => {
    const ref = useSetCustomCssVariable({ '--app-bar-height': 'height' });
    return (
        <MainAppBarContent
            isLoggedIn={isLoggedIn}
            mobile={false}
            admin={admin}
            sx={{ display: { xs: 'none', sm: 'flex' } }}
            ref={ref}
        />
    );
};

export default MainAppBarDesktop;
