'use client';

import { Button, type SxProps } from '@mui/material';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import RoomSelectDialog from './RoomSelectDialog';
import { useState } from 'react';
import type { ConEvent } from '$lib/types';
import { addToRoom } from './actions';
import { PoolName, RoomName } from '$lib/enums';

type Props = {
    roomCoordinates: {
        x: number;
        y: number;
    };
    events: ConEvent[];
    poolName: PoolName;
    roomName: RoomName;
};

const RoomAddButton = ({ roomCoordinates, events, poolName, roomName }: Props) => {
    const [open, setOpen] = useState<boolean>(false);
    const [selectedValue, setSelectedValue] = useState<string>('');

    const handleAddClick = () => {
        setOpen(true);
    };
    const handleClose = (value: string) => {
        setOpen(false);
        setSelectedValue(value);
        addToRoom(value, roomName, poolName);
    };
    const roomCoordinatesSx: SxProps = {
        position: 'absolute',
        left: roomCoordinates.x,
        top: roomCoordinates.y,
    };

    return (
        <>
            <Button sx={{ fontSize: '90px', color: 'lightgray', ...roomCoordinatesSx }} onClick={handleAddClick}>
                <AddCircleIcon sx={{ fontSize: '90px' }} />
            </Button>

            <RoomSelectDialog
                open={open}
                selectedValue={selectedValue}
                onClose={handleClose}
                events={events}
            ></RoomSelectDialog>
        </>
    );
};

export default RoomAddButton;
