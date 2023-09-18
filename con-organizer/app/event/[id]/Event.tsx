'use client';
import { useState } from 'react';
import { Box, Button, Card, CardActions } from '@mui/material';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';
import MainNavigator from '@/components/mainNavigator';
import { useSingleEvents } from '@/lib/hooks/UseSingleEvent';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const { event, loading } = useSingleEvents(id);
    const [openEdit, setOpenEdit] = useState(false);

    const handleCloseEdit = () => {
        setOpenEdit(false);
    };

    const handleOpenEdit = () => {
        setOpenEdit(true);
    };

    return (
        <Box sx={{ maxWidth: '1080px', margin: '0 auto' }}>
            {loading && <h1>Loading...</h1>}
            <EditDialog open={openEdit} handleClose={handleCloseEdit} conEvent={event} />

            <Card>
                <Button onClick={() => window.history.go(-1)}>Tilbake</Button>
            </Card>

            <EventUi conEvent={event} showSelect={true} />

            <Card>
                <CardActions>
                    <Button onClick={handleOpenEdit}>Endre</Button>
                </CardActions>
            </Card>
            <MainNavigator />
        </Box>
    );
};

export default Event;
