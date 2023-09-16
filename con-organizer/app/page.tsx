import { Box } from '@mui/material';
import { AuthProvider } from '@/components/AuthProvider';
import DayTab from '@/components/dayTab';
import EventList from '@/components/eventList';
import MainNavigator from '@/components/mainNavigator';
import { Theme } from '@/components/ThemeProvider';

export default function Home() {
    return (
        <main className="">
            <Theme>
                <DayTab />
                <Box className="flex flex-row flex-wrap justify-center gap-4">
                    <AuthProvider>
                        <EventList />
                    </AuthProvider>
                </Box>
                <MainNavigator />
            </Theme>
        </main>
    );
}
