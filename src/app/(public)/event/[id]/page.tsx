import { Metadata } from 'next';
import MainEvent from './components/MainEvent';
import { getAllEvents, getAllPoolEventsSortedByDay, getEventById } from '$app/(public)/components/lib/serverAction';
import MainEventBig from './components/MainEventBig';
import BigMediaQueryWrapper from './components/ui/BigMediaQueryWrapper';
import SmallMediaQueryWrapper from './components/ui/SmallMediaQueryWrapper';
import { Box } from '@mui/material';

type Props = {
    params: { id: string };
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
    const id = params.id;
    const event = await getEventById(id);

    return {
        title: event.title,
        description: event.shortDescription,
    };
}

const EventPage = async ({ params: { id } }: Props) => {
    const events = await getAllEvents();
    const eventIndex = events.findIndex((event) => event.id === id);
    const prevNavigationId = events[eventIndex - 1]?.id;
    const nextNavigationId = events[eventIndex + 1]?.id;
    const event = await getEventById(id);
    const poolEvents = await getAllPoolEventsSortedByDay();
    // console.log(poolEvents[0].poolEvents, poolEvents[0].day);
    // console.log(poolEvents[1].poolEvents, poolEvents[1].day);
    // console.log(poolEvents[2].poolEvents, poolEvents[2].day);
    // console.log(poolEvents[3].poolEvents, poolEvents[3].day);

    return (
        <>
            <SmallMediaQueryWrapper>
                <MainEvent id={id} prevNavigationId={prevNavigationId} nextNavigationId={nextNavigationId} />
            </SmallMediaQueryWrapper>

            <BigMediaQueryWrapper>
                <Box sx={{ display: 'grid', placeContent: ' center' }}>
                    <MainEventBig
                        event={event}
                        prevNavigationId={prevNavigationId}
                        nextNavigationId={nextNavigationId}
                    />
                </Box>
            </BigMediaQueryWrapper>
        </>
    );
};
export default EventPage;
