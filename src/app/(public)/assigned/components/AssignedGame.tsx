import { Box, Link, Typography } from '@mui/material';
import { getAssignedGameByDay } from './lib/helper';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { cookies } from 'next/headers';
import type { ParticipantCookie } from '$lib/types';
import NextLink from 'next/link';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { getAllPoolEvents, getPoolEventById } from '$app/(public)/components/lib/serverAction';
import LoadingAssignedGameWrapper from './ui/LoadingAssignedGameWrapper';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';
import EventCardBig from '$app/(public)/components/components/EventCardBig';

type Props = {};

const AssignedGame = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();

    const cookie = cookies();
    console.log(cookie);

    let myParticipants;
    try {
        myParticipants = JSON.parse(cookie.get('myParticipants')?.value ?? '{}') as ParticipantCookie[] | undefined;
    } catch (error) {
        myParticipants = undefined;
    }
    if (!myParticipants || !user) {
        console.warn('fant ikke deltakere i cookie for bruker: ', user?.uid ?? 'ingen bruker logget inn');
        return (
            <>
                {user ?
                    <ParticipantSelector />
                :   null}
                <Typography sx={{ display: 'inline-block' }}>
                    {user ?
                        'Du må ha billett for og se denne sida'
                    :   'Du må være logget inn for og se dine på meldinger.'}
                </Typography>{' '}
                {user ?
                    <Link component={NextLink} href={'/my-profile/my-tickets'}>
                        mer info her
                    </Link>
                :   null}
            </>
        );
    }
    const currentGamesForParticipant = await getAssignedGameByDay(
        myParticipants.find((participant) => participant.isSelected)?.id ?? ''
    );
    console.log('myParticipants', myParticipants);

    if (!currentGamesForParticipant) {
        return (
            <>
                <ParticipantSelector />
                <Typography>Ingen påmeldinger</Typography>;
            </>
        );
    }
    const poolEvents = await getAllPoolEvents();
    return (
        <LoadingAssignedGameWrapper>
            <Box sx={{ display: 'grid', placeContent: 'center' }}>
                <ParticipantSelector />
            </Box>
            {currentGamesForParticipant ?
                [...poolEvents.entries()].map(([poolName, poolEvents]) => {
                    return (
                        <Box sx={{ maxWidth: '24.7143rem' }} key={poolName}>
                            <Typography variant="h1">{translatedDays.get(poolName)}</Typography>
                            {poolEvents
                                .filter((poolEvent) =>
                                    currentGamesForParticipant.some((game) => game.poolEventId === poolEvent.id)
                                )
                                .map((poolEvent) => {
                                    return (
                                        <EventCardBig
                                            gameMaster={poolEvent?.gameMaster}
                                            shortDescription={poolEvent?.shortDescription}
                                            title={poolEvent?.title}
                                            system={poolEvent?.system}
                                            backgroundImage={poolEvent.smallImageURL}
                                            icons={poolEvent?.icons ?? []}
                                        />
                                    );
                                })}
                            {/* {JSON.stringify(currentGamesForParticipant)} */}
                        </Box>
                    );
                })
            : currentGamesForParticipant.length === 0 ?
                <Typography>Ingen påmeldinger for denne puljen</Typography>
            :   <Typography>Ingen påmeldinger</Typography>}
        </LoadingAssignedGameWrapper>
    );
};

export default AssignedGame;
