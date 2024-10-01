import { Box } from '@mui/material';
import { getAllPoolEvents, migrateInterestsToParticipantInterests } from './components/lib/serverAction';
import DaysHeader from './components/ui/DaysHeader';
import EventList from './components/EventList';
import Logo from './components/ui/Logo';
import RealtimePoolEvents from '$lib/components/RealtimePoolEvents';
import type { IconName } from '$lib/types';
import { migrateParticipantAndInterest } from '$lib/serverActions/Migration';
import Test from './Test';

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
            <Test />
        </>
    );
}
