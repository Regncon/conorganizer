'use client';
import { Box, Typography } from '@mui/material';
import { PoolName, RoomName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { useState } from 'react';
import RoomSelectDialog from './RoomSelectDialog';
import RoomCard from './RoomCard';

type Props = {
    poolName: PoolName;
    roomName: RoomName;
    title: string;
    gameMaster: string;
    system: string;
    imageUri: string;
    events: ConEvent[];
};
const RoomMapItem = ({ roomName, title, gameMaster, system, imageUri, events }: Props) => {
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
                    position: 'absolute',
                    top: roomCoordinates().y,
                    left: roomCoordinates().x,
                }}
            >
                <RoomCard title={title} gameMaster={gameMaster} system={system} imageUri={imageUri} />
            </Box>
            <RoomSelectDialog selectedValue={selectedValue} open={open} onClose={handleClose} events={events} />
        </>
    );
};
export default RoomMapItem;
