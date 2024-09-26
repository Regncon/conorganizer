import { Box } from '@mui/material';
import { getAllPoolEvents } from './components/lib/serverAction';
import DaysHeader from './components/ui/DaysHeader';
import EventList from './components/EventList';
import Logo from './components/ui/Logo';
import RealtimePoolEvents from '$lib/components/RealtimePoolEvents';
import type { IconName } from '$lib/types';
type Props = {
    searchParams: {
        [key in IconName]: string;
    };
};
export default async function Home({ searchParams }: Props) {
    const poolEvents = await getAllPoolEvents();

    return (
        <>
            <Box>
                <Logo />
                <DaysHeader />
                <EventList events={poolEvents} searchParams={searchParams} />
            </Box>
            <RealtimePoolEvents />
        </>
    );
}
