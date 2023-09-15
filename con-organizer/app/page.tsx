import { Box } from '@mui/material';
import { AuthProvider } from '@/components/AuthProvider';
import DayTab from '@/components/dayTab';
import EventList from '@/components/eventList';
import MainNavigator from '@/components/mainNavigator';

export default function Home() {
    return (
        <main className="">
            <DayTab />
            <Box className="flex flex-row flex-wrap justify-center gap-4">
                <AuthProvider>
                    <EventList />
                </AuthProvider>
            </Box>
            <MainNavigator />
        </main>
    );
}