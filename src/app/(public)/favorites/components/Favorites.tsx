import { getAllPoolEvents, getUsersInterestById } from '$app/(public)/components/lib/serverAction';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Box, Typography } from '@mui/material';
import { buildParticipantPoolEventsMap } from './lib/helpers/helpers';
import EventCardBig from '$app/(public)/components/components/EventCardBig';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';
import { Fragment } from 'react';
import Link from 'next/link';
import type { Route } from 'next';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { cookies } from 'next/headers';
import type { ParticipantCookie } from '$lib/types';
import LoadingParticipantWrapper from './ui/LoadingParticipantWrapper';

type Props = {};

const Favorites = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    if (!user) {
        return null;
    }
    const cookie = cookies();
    const myParticipants = JSON.parse(cookie.get('myParticipants')?.value ?? '') as ParticipantCookie[] | undefined;
    if (!myParticipants) {
        throw new Error('Fant ikkje participant i cookie');
    }

    const activeParticipant = myParticipants?.filter((participant) => participant.isSelected)[0];
    const participantName = `${activeParticipant.firstName} ${activeParticipant.lastName}`;

    const usersInterests = await getUsersInterestById(user.uid);
    const poolEvents = await getAllPoolEvents();
    const participantMap = buildParticipantPoolEventsMap(usersInterests, poolEvents);
    const currentParticipant = participantMap.get(participantName);
    if (!currentParticipant) {
        throw new Error('Fant ikkje participant i participantMap');
    }

    return (
        <>
            <Typography variant="h1">Her kan du sjå favorittane til alle deltakarane dine.</Typography>
            <Box sx={{ display: 'grid', placeContent: 'center' }}>
                <ParticipantSelector />
            </Box>
            <LoadingParticipantWrapper>
                {[...currentParticipant.entries()].map(([poolName, poolEvents]) => {
                    return (
                        <Box key={poolName}>
                            <Typography variant="h1">{translatedDays.get(poolName)}</Typography>
                            <Box
                                sx={{
                                    display: 'grid',
                                    gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 345px))',
                                    gap: '1rem',
                                }}
                            >
                                {poolEvents.map((poolEvent) => {
                                    return (
                                        <Fragment key={poolEvent.id}>
                                            <Box
                                                component={Link}
                                                key={poolEvent.id}
                                                sx={{ textDecoration: 'none' }}
                                                prefetch
                                                href={`/event/${poolEvent.id}` as Route}
                                            >
                                                <EventCardBig
                                                    title={poolEvent.title}
                                                    gameMaster={poolEvent.gameMaster}
                                                    shortDescription={poolEvent.shortDescription}
                                                    system={poolEvent.system}
                                                    backgroundImage={poolEvent.smallImageURL}
                                                />
                                            </Box>
                                        </Fragment>
                                    );
                                })}
                            </Box>
                        </Box>
                    );
                })}
            </LoadingParticipantWrapper>
        </>
    );
};

export default Favorites;
