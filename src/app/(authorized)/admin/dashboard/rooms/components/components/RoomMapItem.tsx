'use client';
import { Box, Button, Dialog, DialogActions, DialogTitle, Typography } from '@mui/material';
import { PoolName, RoomName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { useState } from 'react';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import DeleteIcon from '@mui/icons-material/Delete';
import RoomAddButton from './components/RoomAddButton';
import { convertToPoolEvent, removeFromRoom, removeFromPool } from './lib/actions';
import RoomCard from './ui/RoomCard';
import RoomSelectDialog from './ui/RoomSelectDialog';

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
        if (roomName !== RoomName.NotSet) {
            removeFromRoom(selectedDeleteEvent?.id ?? '', roomName, poolName);
        } else removeFromPool(selectedDeleteEvent?.id ?? '', poolName);
    };

    const filteredEvents = events.filter(poolFilters[poolName]);
    let eventsInRoom: ConEvent[] = [];
    filteredEvents.forEach((event) => {
        const roomIds = event.roomIds;
        if (roomIds) {
            roomIds.forEach((roomId) => {
                if (roomId.poolName === poolName && roomId.roomName === roomName) {
                    eventsInRoom.push(event);
                }
            });
        }
        if (roomName === RoomName.NotSet) {
            if (!roomIds || !roomIds.find((roomId) => roomId.poolName === poolName)) {
                eventsInRoom.push(event);
            }
        }
    });

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
                    backgroundColor: 'white',
                }}
            >
                <Typography variant="h2" sx={{ fontWeight: 'bold', margin: '0', color: 'black' }}>
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
                <DialogTitle>Slett?</DialogTitle>
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
