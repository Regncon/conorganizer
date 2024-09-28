import { getAllPoolEvents, getUsersInterestById } from '$app/(public)/components/lib/serverAction';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Box, Typography } from '@mui/material';
import { buildParticipantPoolEventsMap } from './lib/helpers/helpers';
import EventCardBig from '$app/(public)/components/components/EventCardBig';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';
import { Fragment } from 'react';
import Link from 'next/link';
import type { Route } from 'next';

type Props = {};

const Favorites = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    if (!user) {
        return null;
    }
    const usersInterests = await getUsersInterestById(user.uid);
    const poolEvents = await getAllPoolEvents();
    const participantMap = buildParticipantPoolEventsMap(usersInterests, poolEvents);

    return (
        <>
            <Typography variant="h1">Her kan du se dine interesser </Typography>
            {[...participantMap.entries()].map(([name, poolEvents]) => {
                return (
                    <Box key={name}>
                        <Typography variant="h1">{name}</Typography>
                        {[...poolEvents.entries()].map(([poolName, events]) => {
                            return (
                                <Fragment key={poolName}>
                                    <Typography variant="h2" sx={{ fontSize: '2.5rem' }}>
                                        {translatedDays.get(poolName)}
                                    </Typography>
                                    <Box
                                        sx={{
                                            display: 'grid',
                                            gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 345px))',
                                            gap: '1rem',
                                        }}
                                    >
                                        {events.map((event) => {
                                            return (
                                                <Box
                                                    component={Link}
                                                    key={event.id}
                                                    sx={{ textDecoration: 'none' }}
                                                    prefetch
                                                    href={`/event/${event.id}` as Route}
                                                >
                                                    <EventCardBig
                                                        title={event.title}
                                                        gameMaster={event.gameMaster}
                                                        shortDescription={event.shortDescription}
                                                        system={event.system}
                                                        backgroundImage={event.smallImageURL}
                                                    />
                                                </Box>
                                            );
                                        })}
                                    </Box>
                                </Fragment>
                            );
                        })}
                    </Box>
                );
            })}
        </>
    );
};

export default Favorites;
