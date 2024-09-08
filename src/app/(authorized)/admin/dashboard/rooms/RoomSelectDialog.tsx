'use client';
import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import DialogTitle from '@mui/material/DialogTitle';
import Dialog from '@mui/material/Dialog';
import PersonIcon from '@mui/icons-material/Person';
import { blue } from '@mui/material/colors';
import { ConEvent } from '$lib/types';
import { Box, DialogContent, Typography } from '@mui/material';

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
                                                width: '35rem',
                                                display: 'flex',
                                                flexDirection: 'column',
                                                alignItems: 'flex-start',
                                            }}
                                        >
                                            <Typography>{event.title}</Typography>
                                            <Typography sx={{ fontWeight: 'bold', fontSize: '1.1rem' }}>
                                                {event.gameMaster}
                                            </Typography>
                                            <Typography sx={{ color: 'primary.main' }}>{event.system}</Typography>
                                        </Box>
                                        <Box
                                            sx={{
                                                width: '20rem',
                                                display: 'flex',
                                                flexDirection: 'column',
                                                alignItems: 'flex-start',
                                            }}
                                        >
                                            <Typography>Puje: Lørdag kveld</Typography>
                                            <Typography>Rom: Styrerom 1 Gerhard</Typography>
                                            <Typography>Rom: Symfoni</Typography>

                                            <Typography>Puje: Lørdag kveld</Typography>
                                            <Typography>Rom: Styrerom 1 Gerhard</Typography>
                                            <Typography>Rom: Symfoni</Typography>
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
