import { Box } from '@mui/material';
import { ThemeProvider } from '@mui/material';
import { AuthProvider } from '@/components/AuthProvider';
import DayTab from '@/components/dayTab';
import EventList from '@/components/eventList';
import MainNavigator from '@/components/mainNavigator';
import { muiLight } from '@/lib/muiTheme';

export default function Home() {
    return (
        <main className="">
            <ThemeProvider theme={muiLight}>
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
