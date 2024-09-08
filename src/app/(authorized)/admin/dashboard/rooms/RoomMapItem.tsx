'use client';
import { Box, Button, Dialog, DialogActions, DialogTitle, Paper } from '@mui/material';
import { PoolName, RoomName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { useState } from 'react';
import RoomSelectDialog from './RoomSelectDialog';
import RoomCard from './RoomCard';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import DeleteIcon from '@mui/icons-material/Delete';
import { convertToPoolEvent, removeFromPool } from './actions';

type Props = {
    poolName: PoolName;
    eventId: string;
    events: ConEvent[];
};
const RoomMapItem = ({ eventId, poolName, events }: Props) => {
    const [open, setOpen] = useState(false);
    const [selectedValue, setSelectedValue] = useState('');
    const [roomName, setRoomName] = useState<RoomName>(RoomName.NotSet);

    const poolFilters = {
        [PoolName.fridayEvening]: (event: ConEvent) => event.puljeFridayEvening,
        [PoolName.saturdayMorning]: (event: ConEvent) => event.puljeSaturdayMorning,
        [PoolName.saturdayEvening]: (event: ConEvent) => event.puljeSaturdayEvening,
        [PoolName.sundayMorning]: (event: ConEvent) => event.puljeSundayMorning,
    };
    const filteredEvents = events.filter(poolFilters[poolName]);

    const handleClose = (value: string) => {
        setOpen(false);
        setSelectedValue(value);
        if (!value) return;
        convertToPoolEvent(value, poolName);
    };
    const [openDeletDialog, setOpenDeleteDialog] = useState(false);
    const handleClickDeleteOpen = (id: string | undefined) => {
        setSelectedDeleteEvent(events.find((event) => event.id === id));
        setOpenDeleteDialog(true);
    };

    const [selectedDeleteEvent, setSelectedDeleteEvent] = useState<ConEvent>();

    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleCancel = () => {
        setOpenDeleteDialog(false);
    };

    const handleOkDeleteDialog = () => {
        setOpenDeleteDialog(false);
        console.log('selectedDeleteEvent', selectedDeleteEvent);
        removeFromPool(selectedDeleteEvent?.id ?? '', poolName);
    };

    const smallRoomRowX = 2460;
    const styreRomRowX = 1500;

    const roomCoordinates = () => {
        switch (roomName) {
            case RoomName.NotSet:
                return { x: 1000, y: 350 };
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

    // const eventsInRoom = events.filter((event) => event.room === roomName);
    const eventsInRoom = filteredEvents;

    return (
        <>
            {eventsInRoom && (
                <>
                    <Box
                        sx={{
                            position: 'absolute',
                            top: roomCoordinates().y,
                            left: roomCoordinates().x,
                            padding: '1rem',
                            display: 'flex',
                            flexDirection: 'column',
                            gap: '1rem',
                            borderRadius: '1rem',
                            border: '4px solid black',
                            maxHeight: '650px',
                            overflowY: 'auto',
                        }}
                    >
                        {eventsInRoom.map((roomEvent) => {
                            return (
                                <Box
                                    key={roomEvent.id}
                                    sx={{ display: 'flex', backgroundColor: 'inherit', border: 'none' }}
                                >
                                    <RoomCard
                                        title={roomEvent?.title ?? 'Not set'}
                                        gameMaster={roomEvent?.gameMaster ?? 'Not set'}
                                        system={roomEvent?.system ?? 'Not set'}
                                        imageUri={'/blekksprut2.jpg'}
                                    />
                                    <Button onClick={handleClickOpen} sx={{ fontSize: '90px', color: 'lightgray' }}>
                                        <AddCircleIcon sx={{ fontSize: '90px' }} />
                                    </Button>

                                    <Button onClick={() => handleClickDeleteOpen(roomEvent.id)}>
                                        <DeleteIcon sx={{ fontSize: '90px' }} />
                                    </Button>
                                </Box>
                            );
                        })}
                    </Box>
                    <RoomSelectDialog selectedValue={selectedValue} open={open} onClose={handleClose} events={events} />
                    <Dialog open={openDeletDialog}>
                        <DialogTitle>Fjern fra pulje</DialogTitle>
                        <DialogActions>
                            <Button autoFocus onClick={handleCancel}>
                                Cancel
                            </Button>
                            <Button onClick={handleOkDeleteDialog}>Ok</Button>
                        </DialogActions>
                    </Dialog>
                </>
            )}
        </>
    );
};
export default RoomMapItem;
