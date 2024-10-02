import { Box } from '@mui/material';
import { getAllPoolEvents } from './components/lib/serverAction';
import DaysHeader from './components/ui/DaysHeader';
import EventList from './components/EventList';
import Logo from './components/ui/Logo';
import RealtimePoolEvents from '$lib/components/RealtimePoolEvents';

export default async function Home() {
    const poolEvents = await getAllPoolEvents();

    return (
        <>
            <Box>
                <Logo />
                <DaysHeader />
                <EventList events={poolEvents} />
            </Box>
            <RealtimePoolEvents />
        </>
    );
}
