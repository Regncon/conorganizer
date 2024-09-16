import { Metadata } from 'next';
import MainEvent from './components/MainEvent';
import {
    getAdjacentPoolEventsById,
    getAllEvents,
    getAllPoolEvents,
    getEventById,
    getPoolEventById,
} from '$app/(public)/components/lib/serverAction';
import MainEventBig from './components/MainEventBig';
import BigMediaQueryWrapper from './components/ui/BigMediaQueryWrapper';
import SmallMediaQueryWrapper from './components/ui/SmallMediaQueryWrapper';
import { Box } from '@mui/material';
import { PoolName } from '$lib/enums';

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
    const poolEvent = await getPoolEventById(id);
    const { prevNavigationId, nextNavigationId } = await getAdjacentPoolEventsById(id, poolEvent.poolName);

    return (
        <>
            <SmallMediaQueryWrapper>
                <MainEvent id={id} prevNavigationId={prevNavigationId} nextNavigationId={nextNavigationId} />
            </SmallMediaQueryWrapper>

            <BigMediaQueryWrapper>
                <Box sx={{ display: 'grid', placeContent: ' center' }}>
                    <MainEventBig
                        poolEvent={poolEvent}
                        prevNavigationId={prevNavigationId}
                        nextNavigationId={nextNavigationId}
                    />
                </Box>
            </BigMediaQueryWrapper>
        </>
    );
};
export default EventPage;
