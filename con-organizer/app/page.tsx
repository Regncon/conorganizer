import { useState } from 'react';
import { Box } from '@mui/material';
import { ThemeProvider } from '@mui/material';
import { AuthProvider } from '@/components/AuthProvider';
import DayTab from '@/components/dayTab';
import EventList from '@/components/eventList';
import MainNavigator from '@/components/mainNavigator';
import { pool } from '@/lib/enums';
import { muiDark } from '@/lib/muiTheme';

export default function Home() {
/*     const [activePool, setActivePool] = useState<pool>(pool.FirdayEvening);

    const handlePoolChange = (pool: pool) => {
        setActivePool(pool);
    } */

    return (
        <main className="">
            <ThemeProvider theme={muiDark}>
                <AuthProvider>
                    <DayTab />
                    <Box className="flex flex-row flex-wrap justify-center gap-4">
                        <EventList />
                    </Box>
                    <MainNavigator />
                </AuthProvider>
            </ThemeProvider>
        </main>
    );
}
