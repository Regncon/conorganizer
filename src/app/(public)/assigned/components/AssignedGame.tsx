import { Box, Link, Typography } from '@mui/material';
import { getAssignedGameByDay } from './lib/helper';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { cookies } from 'next/headers';
import type { ParticipantCookie } from '$lib/types';
import NextLink from 'next/link';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { getPoolEventById } from '$app/(public)/components/lib/serverAction';
import LoadingAssignedGameWrapper from './ui/LoadingAssignedGameWrapper';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';
import EventCardBig from '$app/(public)/components/components/EventCardBig';

type Props = {};

const AssignedGame = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();

    const cookie = cookies();
    const myParticipants = JSON.parse(cookie.get('myParticipants')?.value ?? '') as ParticipantCookie[] | undefined;
    if (!myParticipants || !user) {
        console.warn('fant ikke deltakere i cookie for bruker: ', user?.uid ?? 'ingen bruker logget inn');
        return (
            <>
                <Typography sx={{ display: 'inline-block' }}>
                    {user ?
                        'Du må ha billett for og se denne sida'
                    :   'Du må være logget inn for og se dine interesser.'}
                </Typography>{' '}
                <Link component={NextLink} href={'/my-profile/my-tickets'}>
                    mer info her
                </Link>
            </>
        );
    }
    const { currentGamesForParticipant, poolName } = await getAssignedGameByDay(
        myParticipants.find((participant) => participant.isSelected)?.id ?? ''
    );
    const poolEvent = await getPoolEventById(currentGamesForParticipant?.poolEventId ?? '');
    console.log(poolEvent, 'poolEvent');

    return (
        <LoadingAssignedGameWrapper>
            <Box sx={{ display: 'grid', placeContent: 'center' }}>
                <ParticipantSelector />
            </Box>
            {currentGamesForParticipant ?
                <Box sx={{ maxWidth: '24.7143rem' }}>
                    <EventCardBig
                        gameMaster={poolEvent?.gameMaster}
                        shortDescription={poolEvent?.shortDescription}
                        title={poolEvent?.title}
                        system={poolEvent?.system}
                        backgroundImage={poolEvent.smallImageURL}
                        icons={poolEvent?.icons ?? []}
                    />
                    {JSON.stringify(currentGamesForParticipant)}
                </Box>
            : poolName && currentGamesForParticipant ?
                <Typography>Ingen på meldinger for {translatedDays.get(poolName)}</Typography>
            :   <Typography>Ingen påmeldinger</Typography>}
        </LoadingAssignedGameWrapper>
    );
};

export default AssignedGame;
