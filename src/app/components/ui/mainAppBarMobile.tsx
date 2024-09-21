'use client';
import { useSetCustomCssVariable } from '$lib/hooks/useSetCustomCssVariable';
import { useMediaQuery } from '@mui/material';
import MainAppBarContent from './MainAppBarContent';

type Props = {
    isLoggedIn: boolean;
    admin: boolean;
};

const MainAppBarMobile = ({ isLoggedIn, admin }: Props) => {
    const isMobile = useMediaQuery('(max-width: 600px)');
    const ref = useSetCustomCssVariable({ '--app-bar-height': 'height' }, isMobile);
    return (
        <MainAppBarContent
            isLoggedIn={isLoggedIn}
            mobile={true}
            admin={admin}
            sx={{
                top: 'auto',
                bottom: 0,
                display: { xs: 'flex', sm: 'none' },
                minHeight: '5.3571rem',
                placeContent: 'center',
            }}
            ref={ref}
        />
    );
};

export default MainAppBarMobile;
