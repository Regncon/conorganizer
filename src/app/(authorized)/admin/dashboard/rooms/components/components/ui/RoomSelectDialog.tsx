'use client';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import DialogTitle from '@mui/material/DialogTitle';
import Dialog from '@mui/material/Dialog';
import { ConEvent } from '$lib/types';
import { Box, DialogContent, Typography } from '@mui/material';
import RoomCard from './RoomCard';

type Props = {
    open: boolean;
    selectedValue: string;
    events: ConEvent[];
    onClose: (value: string) => void;
};

const RoomSelectDialog = ({ open, selectedValue, onClose, events }: Props) => {
    const handleClose = () => {
        onClose(selectedValue);
    };

    const handleListItemClick = (value: string) => {
        onClose(value);
    };

    return (
        <Dialog onClose={handleClose} open={open} fullWidth={true} maxWidth={'md'}>
            <DialogTitle>Velg arrangement</DialogTitle>
            <DialogContent>
                <List sx={{ pt: 0 }}>
                    {events.map((event) => {
                        return (
                            <ListItem disableGutters key={event.id}>
                                <ListItemButton onClick={() => handleListItemClick(event.id ?? '')}>
                                    <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                                        <Box
                                            sx={{
                                                width: '30rem',
                                                display: 'flex',
                                                flexDirection: 'column',
                                                alignItems: 'flex-start',
                                            }}
                                        >
                                            <RoomCard
                                                title={event.title}
                                                gameMaster={event.gameMaster}
                                                system={event.system}
                                                imageUri="/blekksprut2.jpg"
                                            />
                                        </Box>
                                        <Box
                                            sx={{
                                                width: '20rem',
                                                display: 'flex',
                                                flexDirection: 'column',
                                                alignItems: 'flex-start',
                                            }}
                                        >
                                            {event.poolIds?.map((poolId) => {
                                                return (
                                                    <Box key={poolId.id}>
                                                        <Typography>Puje: </Typography>
                                                        <Typography>Rom: Styrerom 1 Gerhard</Typography>
                                                        <Typography>Rom: Symfoni</Typography>
                                                    </Box>
                                                );
                                            })}
                                        </Box>
                                    </Box>
                                </ListItemButton>
                            </ListItem>
                        );
                    })}
                </List>
            </DialogContent>
        </Dialog>
    );
};
export default RoomSelectDialog;
