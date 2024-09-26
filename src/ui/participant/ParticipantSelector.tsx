'use client';
import { ParticipantCookie } from '$lib/types';
import { useState, useEffect } from 'react';
import { Menu, MenuItem, Button, Box } from '@mui/material';
import ParticipantAvatar from './ParticipantAvatar';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { updateEventById } from '$app/(public)/components/lib/serverAction';

type Props = {
    poolEventId?: string;
};

const ParticipantSelector = ({ poolEventId }: Props) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const [participants, setParticipants] = useState<ParticipantCookie[]>([]);

    useEffect(() => {
        const myParticipantsCookie = document.cookie
            .split('; ')
            .find((row) => row.startsWith('myParticipants='))
            ?.split('=')[1];

        const myParticipants: ParticipantCookie[] = JSON.parse(myParticipantsCookie || '[]');
        if (myParticipants) {
            setParticipants(myParticipants);
        }
    }, []);

    const selectedParticipant = participants.find((participant) => participant.isSelected);

    if (!participants || participants.length === 0 || selectedParticipant === undefined) {
        return (
            <Box>
                <Button variant="contained" href="/my-profile/my-tickets">
                    Hent billett
                </Button>
            </Box>
        );
    }

    const handleParticipantSelect = (id: string | undefined) => {
        if (id === undefined) {
            return;
        }
        const updatedParticipants = participants.map((participant) =>
            participant.id === id ? { ...participant, isSelected: true } : { ...participant, isSelected: false }
        );
        setParticipants(updatedParticipants);
        document.cookie = `myParticipants=${JSON.stringify(updatedParticipants)}; path=/;`;
        handleClose();
        updateEventById(poolEventId ?? '');
    };

    return (
        <Box>
            <Button
                aria-controls={open ? 'participant-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
                onClick={handleClick}
                variant="text"
                sx={{ textDecoration: 'none', textTransform: 'none' }}
            >
                <ParticipantAvatar
                    firstName={selectedParticipant.firstName}
                    lastName={selectedParticipant.lastName}
                    small
                />
                <ExpandMoreIcon />
            </Button>
            <Menu
                id="participant-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                    'aria-labelledby': 'participant-button',
                }}
            >
                {participants.map((participant) => (
                    <MenuItem key={participant.id} onClick={() => handleParticipantSelect(participant.id)}>
                        <ParticipantAvatar firstName={participant.firstName} lastName={participant.lastName} />
                    </MenuItem>
                ))}
            </Menu>
        </Box>
    );
};

export default ParticipantSelector;
