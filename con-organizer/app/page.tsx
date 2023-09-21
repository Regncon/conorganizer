import { Box } from '@mui/material';
import EventList from '@/components/eventList';
import MainNavigator from '@/components/mainNavigator';

export default function Home() {

    return (
        <main className="">
            <Box className="flex flex-row flex-wrap justify-center gap-4">
                <EventList />
            </Box>
            <MainNavigator />
        </main>
    );
}
