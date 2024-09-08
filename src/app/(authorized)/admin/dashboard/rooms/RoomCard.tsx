'use client';
import { Box, Typography } from '@mui/material';
import { PoolName, RoomName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { useState } from 'react';
import RoomSelectDialog from './RoomSelectDialog';

type Props = {
    poolName: PoolName;
    roomName: RoomName;
    title: string;
    gameMaster: string;
    system: string;
    imageUri: string;
    events: ConEvent[];
};
const RoomCard = ({ roomName, title, gameMaster, system, imageUri, events }: Props) => {
    const [open, setOpen] = useState(false);
    const [selectedValue, setSelectedValue] = useState('');
    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = (value: string) => {
        setOpen(false);
        setSelectedValue(value);
        console.log(value);
    };

    const smallRoomRowX = 2460;
    const styreRomRowX = 1000;

    const roomCoordinates = () => {
        switch (roomName) {
            case RoomName.Styreromm1:
                return { x: styreRomRowX, y: 350 };
            case RoomName.Styreromm2:
            case RoomName.Styreromm3:
            case RoomName.Styreromm4:
            case RoomName.Styreromm5:
            case RoomName.Styreromm6:
            case RoomName.Klang:
            case RoomName.Sonate:
            case RoomName.Klang:
                return { x: smallRoomRowX, y: 450 };
            case RoomName.Sonate:
                return { x: smallRoomRowX, y: 640 };
            default:
                return { x: 0, y: 0 };
        }
    };
    return (
        <>
            <Box
                onClick={handleClickOpen}
                sx={{
                    padding: '1rem',
                    width: '306px',
                    opacity: 0.9,
                    background: `linear-gradient( rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5) ), url(${imageUri})`,
                    backgroundSize: 'cover',
                    borderRadius: '1.75rem',
                    position: 'absolute',
                    left: roomCoordinates().x,
                    top: roomCoordinates().y,
                }}
            >
                <Typography>{title}</Typography>
                <Typography sx={{ fontWeight: 'bold', fontSize: '1.1rem' }}> {gameMaster} </Typography>
                <Typography sx={{ color: 'primary.main' }}>{system}</Typography>
            </Box>

            <RoomSelectDialog selectedValue={selectedValue} open={open} onClose={handleClose} events={events} />
        </>
    );
};
export default RoomCard;
