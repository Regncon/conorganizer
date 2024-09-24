'use client';
import { Participant } from '$lib/types';
import { useState } from 'react';
import { Menu, MenuItem, Button } from '@mui/material';
import ParticipantAvatar from './ParticipantAvatar';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

type Props = {
    participants: Participant[];
    activeParticipantId: string;
};

const ParticipantSelector = ({ participants, activeParticipantId }: Props) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    return (
        <>
            <Button
                aria-controls={open ? 'participant-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
                onClick={handleClick}
                variant="text"
            >
                <ParticipantAvatar name={'Ola Nordmann'} />
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
                    <MenuItem key={participant.id} onClick={handleClose}>
                        <ParticipantAvatar name={`${participant.firstName} ${participant.lastName}`} />
                    </MenuItem>
                ))}
            </Menu>
        </>
    );
};

export default ParticipantSelector;
