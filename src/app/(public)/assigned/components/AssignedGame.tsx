import { Box, Link, Typography } from '@mui/material';
import { getAssignedGameByDay } from './lib/helper';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { cookies } from 'next/headers';
import type { ParticipantCookie } from '$lib/types';
import NextLink from 'next/link';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';

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
    const test = await getAssignedGameByDay(myParticipants.find((participant) => participant.isSelected)?.id ?? '');
    return (
        <>
            <Box sx={{ display: 'grid', placeContent: 'center' }}>
                <ParticipantSelector />
            </Box>
            {JSON.stringify(test)}
        </>
    );
};

export default AssignedGame;
