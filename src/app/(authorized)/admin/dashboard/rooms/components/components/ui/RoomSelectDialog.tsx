'use client';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import DialogTitle from '@mui/material/DialogTitle';
import Dialog from '@mui/material/Dialog';
import { ConEvent } from '$lib/types';
import { Box, DialogContent, Divider } from '@mui/material';
import RoomCard from './RoomCard';
import RoomPoolInfo from './RoomPoolInfo';
import { PoolName } from '$lib/enums';
import UnwantedPoolByGm from './UnwantedPoolByGm';
import { Fragment } from 'react';

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

    const pools = [
        { poolName: PoolName.fridayEvening },
        { poolName: PoolName.saturdayMorning },
        { poolName: PoolName.saturdayEvening },
        { poolName: PoolName.sundayMorning },
    ];

    return (
        <Dialog onClose={handleClose} open={open} fullWidth={true} maxWidth={'md'}>
            <DialogTitle>Velg arrangement</DialogTitle>
            <DialogContent>
                <List sx={{ pt: 0 }}>
                    {events
                        .toSorted((a, b) => a.title.localeCompare(b.title))
                        .map((conEvent) => {
                            return (
                                <Fragment key={conEvent.id}>
                                    <ListItem disableGutters>
                                        <ListItemButton onClick={() => handleListItemClick(conEvent.id ?? '')}>
                                            <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                                                <Box
                                                    sx={{
                                                        width: '24rem',
                                                        display: 'flex',
                                                        flexDirection: 'column',
                                                        alignItems: 'flex-start',
                                                    }}
                                                >
                                                    <RoomCard
                                                        title={conEvent.title}
                                                        gameMaster={conEvent.gameMaster}
                                                        system={conEvent.system}
                                                        imageUri={conEvent.smallImageURL ?? '/dice-small.webp'}
                                                    />
                                                </Box>
                                                <Box>
                                                    {pools.map((pool) => {
                                                        return (
                                                            <Fragment key={pool.poolName}>
                                                                <UnwantedPoolByGm
                                                                    poolName={pool.poolName}
                                                                    conEvent={conEvent}
                                                                />
                                                                <RoomPoolInfo
                                                                    poolName={pool.poolName}
                                                                    conEvent={conEvent}
                                                                />
                                                            </Fragment>
                                                        );
                                                    })}
                                                </Box>
                                            </Box>
                                        </ListItemButton>
                                    </ListItem>
                                    <Divider />
                                </Fragment>
                            );
                        })}
                </List>
            </DialogContent>
        </Dialog>
    );
};
export default RoomSelectDialog;
