import { Box } from '@mui/material';
import EventList from '@/components/EventList';
import MainNavigator from '@/components/MainNavigator';

export default function Home() {
    return (
        <main>
            <Box className="flex flex-row flex-wrap justify-center gap-4">
                <EventList />
            </Box>
            <MainNavigator />
        </main>
    );
}
