'use client';

import { Participant, ParticipantLocalStorage } from '$lib/types';
import { Box, Typography } from '@mui/material';
import { useEffect } from 'react';

type Props = {
    participants: Participant[];
};

const GenerateNewParticipantStorage = (participants: Participant[]): ParticipantLocalStorage[] => {
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

        localStorage.setItem('myParticipants', JSON.stringify(newParticipants));
    }, [participants]);

    return (
        <Box>
            <Typography>En smart hjelpetekst skrevet av en som ikke er meg eller dyslektiker</Typography>
            <Typography variant="h1">Mine billetter</Typography>{' '}
        </Box>
    );
};

export default MyParticipantsHeader;
