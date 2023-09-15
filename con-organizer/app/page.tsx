import { Box } from '@mui/material';
import { AuthProvider } from '@/components/auth';
import { Login } from '@mui/icons-material';
import EventList from '@/components/eventList';
import DayTab from '@/components/dayTab';
import MainNavigator from '@/components/mainNavigator';
import Dialog from '@mui/material';

export default function Home() {
    return (
        <main className=''>
            <DayTab />

            <Box className='flex flex-row flex-wrap justify-center gap-4'>
                <AuthProvider>
                    <Login />
                    <EventList />
                </AuthProvider>
            </Box>
            <MainNavigator />
        </main>
    );
}
