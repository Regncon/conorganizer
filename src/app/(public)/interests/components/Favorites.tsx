import { getAllPoolEvents, getUsersInterestById } from '$app/(public)/components/lib/serverAction';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Box, Link, Typography } from '@mui/material';

import EventCardBig from '$app/(public)/components/components/EventCardBig';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';

import NextLink from 'next/link';
import type { Route } from 'next';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { cookies } from 'next/headers';
import type { ParticipantCookie } from '$lib/types';

import {
    interestLevelToImage,
    InterestLevelToLabel,
} from '$app/(authorized)/admin/dashboard/events/event-dashboard/[id]/[tab]/interest/components/lib/helpers/InterestHelper';
import Image from 'next/image';
import { buildParticipantPoolEventsMap } from './lib/helpers/helpers';
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
        console.warn('fant ikke deltaker i cookie for bruker: ', user.uid);
        return (
            <>
                <Typography sx={{ display: 'inline-block' }}>Du må ha billett for og se denne sida</Typography>{' '}
                <Link component={NextLink} href={'/my-profile/my-tickets'}>
                    mer info her
                </Link>
            </>
        );
    }

    const activeParticipant = myParticipants?.filter((participant) => participant.isSelected)[0];
    if (!activeParticipant) {
        console.warn('fant ikke deltaker i myParticipants for bruker: ', user.uid);
        return (
            <>
                <Typography sx={{ display: 'inline-block' }}>Du må ha billett for og se denne sida</Typography>{' '}
                <Link component={NextLink} href={'/my-profile/my-tickets'}>
                    mer info her
                </Link>
            </>
        );
    }

    const participantName = `${activeParticipant.firstName} ${activeParticipant.lastName}`;

    const usersInterests = await getUsersInterestById(user.uid);
    const poolEvents = await getAllPoolEvents();
    try {
        const participantMap = await buildParticipantPoolEventsMap(usersInterests, poolEvents);
        const currentParticipant = participantMap.get(participantName);

        if (!currentParticipant) {
            throw new Error('Fant ikkje participant i participantMap');
        }

        const myInterests = myParticipants.length === 1 ? true : false;
        const participantsInterests = myParticipants.length > 1 ? true : false;
        const headerText =
            myInterests ? `Her kan du sjå dine interesser.`
            : participantsInterests ? `Her kan du sjå interessene til alle deltakarane dine.`
            : null;

        return (
            <>
                <Typography variant="h1">{headerText}</Typography>
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
                                            <Box key={poolEvent.id}>
                                                <Box sx={{ display: 'flex', gap: '1rem' }}>
                                                    <Image
                                                        src={interestLevelToImage[poolEvent.interestLevel]}
                                                        alt={InterestLevelToLabel[poolEvent.interestLevel]}
                                                        width={100}
                                                        height={50}
                                                    />
                                                    <Typography variant="h3">
                                                        {InterestLevelToLabel[poolEvent.interestLevel]}
                                                    </Typography>
                                                </Box>
                                                <Box
                                                    component={NextLink}
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
                                            </Box>
                                        );
                                    })}
                                </Box>
                            </Box>
                        );
                    })}
                </LoadingParticipantWrapper>
            </>
        );
    } catch (error) {
        const err = error as Error;
        console.log(err.message);
        if (err.message === 'Fant ikkje participant i participantMap') {
            return (
                <Box sx={{ display: 'grid', marginBlockStart: '1rem' }}>
                    <Box sx={{ placeSelf: 'center' }}>
                        <ParticipantSelector />
                    </Box>
                    <Typography variant="h1">Denne deltakaren har ikkje valt nokre interesser.</Typography>
                </Box>
            );
        }

        return (
            <Typography>
                Noe gjekk gale kontakt styret <Link href="mailto:styret@regncon.no">styret@regncon.no</Link>
            </Typography>
        );
    }
};

export default Favorites;
