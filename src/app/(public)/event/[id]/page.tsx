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
import { PoolName, type InterestLevel } from '$lib/enums';
import RealtimePoolEvent from './components/components/RealtimePoolEvent';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { cookies } from 'next/headers';
import type { ParticipantCookie } from '$lib/types';
import { getInterest } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';

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
    const { user } = await getAuthorizedAuth();
    const claims = (await user?.getIdTokenResult())?.claims;
    const isAdmin = claims?.admin ?? false;

    const cookie = cookies();
    const activeParticipantsString = cookie.get('myParticipants');
    const activeParticipants: ParticipantCookie[] = JSON.parse(activeParticipantsString?.value ?? '');
    const activeParticipantId = activeParticipants?.find((participant) => participant.isSelected)?.id;
    const interestLevel = await getInterest(activeParticipantId, poolEvent.id);

    return (
        <>
            <SmallMediaQueryWrapper>
                <MainEvent
                    id={id}
                    prevNavigationId={prevNavigationId}
                    nextNavigationId={nextNavigationId}
                    isAdmin={isAdmin}
                    interestLevel={interestLevel}
                />
            </SmallMediaQueryWrapper>

            <BigMediaQueryWrapper>
                <Box sx={{ display: 'grid', placeItems: ' center' }}>
                    <MainEventBig
                        poolEvent={poolEvent}
                        prevNavigationId={prevNavigationId}
                        nextNavigationId={nextNavigationId}
                        isAdmin={isAdmin}
                        interestLevel={interestLevel}
                    />
                </Box>
            </BigMediaQueryWrapper>
            <RealtimePoolEvent id={id} />
        </>
    );
};
export default EventPage;
