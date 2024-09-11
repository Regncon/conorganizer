'use client';

import { Button, type SxProps } from '@mui/material';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import RoomSelectDialog from '../RoomSelectDialog';
import { useState } from 'react';
import type { ConEvent } from '$lib/types';
import { addToRoom } from '../lib/actions';
import { PoolName, RoomName } from '$lib/enums';

type Props = {
    events: ConEvent[];
    poolName: PoolName;
    roomName: RoomName;
};

const RoomAddButton = ({ events, poolName, roomName }: Props) => {
    const [open, setOpen] = useState<boolean>(false);
    const handleAddClick = () => {
        setOpen(true);
    };
    const handleClose = async (value: string) => {
        setOpen(false);
        if (!value) {
            return;
        }
        addToRoom(value, roomName, poolName);
    };

    const dayFilteredEvents = events.filter((event) => event.poolIds?.some((pool) => pool.poolName === poolName));

    return (
        <>
            <Button sx={{ fontSize: '90px', color: 'lightgray' }} onClick={handleAddClick}>
                <AddCircleIcon sx={{ fontSize: '90px' }} />
            </Button>

            <RoomSelectDialog open={open} selectedValue={''} onClose={handleClose} events={events}></RoomSelectDialog>
        </>
    );
};

export default RoomAddButton;
