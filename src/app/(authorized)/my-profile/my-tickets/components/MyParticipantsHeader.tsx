'use client';

import { Participant, ParticipantCookie } from '$lib/types';
import { Box, Typography } from '@mui/material';
import { useEffect } from 'react';

type Props = {
    participants: Participant[];
};

const GenerateNewParticipantStorage = (participants: Participant[]): ParticipantCookie[] => {
    const newParticipants = participants.map((participant, i) => {
        return {
            id: participant.id,
            firstName: participant.firstName,
            lastName: participant.lastName,
            isSelected: i === 0 ? true : false,
        };
    });

    return newParticipants;
};

const MyParticipantsHeader = ({ participants }: Props) => {
    useEffect(() => {
        const newParticipants = GenerateNewParticipantStorage(participants);
        console.log(newParticipants, 'newParticipants');
        const twoWeekExpire = 14 * 24 * 60 * 60 * 1000;
        const expirationDate = Date.now() + twoWeekExpire;
        document.cookie = `myParticipants=${JSON.stringify(newParticipants)}; expires=${new Date(expirationDate).toUTCString()}; path=/`;
    }, [participants]);

    return (
        <Box>
            <Typography variant="h1" sx={{ marginBlockStart: '0' }}>
                Mine billetter
            </Typography>{' '}
            <Typography variant="h4">
                Vi fann følgande billettar på di bestilling. Du kan legga til eigne epostadresser for kvar billett
                nedanfor, slik at kvar deltakar kan melda seg på arrangement på eiga hand.
            </Typography>
        </Box>
    );
};

export default MyParticipantsHeader;
