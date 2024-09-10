'use client';
import { Box, Button, Dialog, DialogActions, DialogTitle, Typography } from '@mui/material';
import { PoolName, RoomName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { useState } from 'react';
import RoomSelectDialog from './RoomSelectDialog';
import RoomCard from './RoomCard';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import DeleteIcon from '@mui/icons-material/Delete';
import { convertToPoolEvent, removeFromPool, removeFromRoom } from './actions';
import RoomAddButton from './RoomAddButton';

type Props = {
    poolName: PoolName;
    roomName: RoomName;
    top: number;
    left: number;
    events: ConEvent[];
};
const RoomMapItem = ({ roomName, top, left, poolName, events }: Props) => {
    const [open, setOpen] = useState(false);
    const [selectedValue, setSelectedValue] = useState('');

    const poolFilters = {
        [PoolName.fridayEvening]: (event: ConEvent) => event.puljeFridayEvening,
        [PoolName.saturdayMorning]: (event: ConEvent) => event.puljeSaturdayMorning,
        [PoolName.saturdayEvening]: (event: ConEvent) => event.puljeSaturdayEvening,
        [PoolName.sundayMorning]: (event: ConEvent) => event.puljeSundayMorning,
    };

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
        if (roomName !== RoomName.NotSet) {
            removeFromRoom(selectedDeleteEvent?.id ?? '', roomName, poolName);
        } else removeFromPool(selectedDeleteEvent?.id ?? '', poolName);
    };
    //
    // const smallRoomRowX = 2460;
    // const styreRomRowX = 1500;
    //
    // const roomCoordinates = () => {
    //     switch (roomName) {
    //         case RoomName.NotSet:
    //             return { x: 1000, y: 350 };
    //         case RoomName.Styrerom1:
    //             return { x: styreRomRowX, y: 350 };
    //         case RoomName.Styrerom2:
    //         case RoomName.Styrerom3:
    //         case RoomName.Styrerom4:
    //         case RoomName.Styrerom5:
    //         case RoomName.Styrerom6:
    //         case RoomName.Klang:
    //         case RoomName.Sonate:
    //         case RoomName.Klang:
    //             return { x: smallRoomRowX, y: 450 };
    //         case RoomName.Sonate:
    //             return { x: smallRoomRowX, y: 640 };
    //         default:
    //             return { x: 0, y: 0 };
    //     }
    // };

    const filteredEvents = events.filter(poolFilters[poolName]);
    // console.log('filteredEvents', filteredEvents);
    let eventsInRoom: ConEvent[] = [];
    filteredEvents.forEach((event) => {
        const roomIds = event.roomIds;
        if (roomIds) {
            // console.log('roomIds', roomIds);
            roomIds.forEach((roomId) => {
                if (roomId.poolName === poolName && roomId.roomName === roomName) {
                    console.log('roomId', roomId);
                    eventsInRoom.push(event);
                } else if (roomName === RoomName.NotSet) {
                    eventsInRoom.push(event);
                    // console.log(event, 'event');
                }
            });
        }
    });
    // const filteredEventsInRoom = filteredEvents.filter((event) => event.roomIds? === true);
    // console.log('filteredEventsInRoom', filteredEventsInRoom.length, filteredEventsInRoom);
    // const eventsInRoom = events.filter(poolFilters[poolName]).filter((event) => event.room === roomName);
    // const eventsInRoom = filteredEvents;

    // console.log(roomName, 'roomName', 'top:', top, 'left:', left);
    return (
        <>
            <Box
                sx={{
                    position: 'absolute',
                    top: top,
                    left: left,
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
                <Typography variant="h6" sx={{ color: 'black', display: 'none' }}>
                    {roomName}
                </Typography>
                <Box sx={{ display: 'flex' }}>
                    <Box>
                        {eventsInRoom.map((roomEvent) => {
                            return (
                                <Box
                                    key={roomEvent.id}
                                    sx={{
                                        display: 'flex',
                                        backgroundColor: 'inherit',
                                        border: 'none',
                                        paddingBottom: '1rem',
                                    }}
                                >
                                    <RoomCard
                                        title={roomEvent?.title ?? 'Not set'}
                                        gameMaster={roomEvent?.gameMaster ?? 'Not set'}
                                        system={roomEvent?.system ?? 'Not set'}
                                        imageUri={'/blekksprut2.jpg'}
                                    />

                                    <Button onClick={() => handleClickDeleteOpen(roomEvent.id)}>
                                        <DeleteIcon sx={{ fontSize: '90px' }} />
                                    </Button>
                                </Box>
                            );
                        })}
                    </Box>
                    {roomName === RoomName.NotSet ?
                        <Button onClick={handleClickOpen} sx={{ fontSize: '90px', color: 'lightgray' }}>
                            <AddCircleIcon sx={{ fontSize: '90px' }} />
                        </Button>
                        : <RoomAddButton events={events} poolName={poolName} roomName={roomName} />}
                </Box>
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
    );
};
export default RoomMapItem;
