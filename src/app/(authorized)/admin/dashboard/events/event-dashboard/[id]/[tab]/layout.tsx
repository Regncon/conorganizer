import { Paper } from '@mui/material';
import EventDashboardTabs from './components/EventDashboardTabs';
import type { PropsWithChildren } from 'react';
import type { TabNames } from './lib/types/types';

type Props = {
    params: {
        id: string;
        tab: TabNames;
    };
};

const Layout = async ({ params: { id, tab }, children }: PropsWithChildren<Props>) => {
    let currentTabValue = 4;
    switch (tab) {
        case 'edit':
            currentTabValue = 4;
            break;
        case 'room':
            currentTabValue = 3;
            break;
        case 'settings':
            currentTabValue = 2;
            break;
        default:
            currentTabValue = 4;
            break;
    }
    return (
        <>
            <Paper>
                <EventDashboardTabs id={id} value={currentTabValue} />
            </Paper>
            {children}
        </>
    );
};

export default Layout;
