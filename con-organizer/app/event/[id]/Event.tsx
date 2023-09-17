'use client';
import { useState } from 'react';
import { Button, Card, CardActions } from '@mui/material';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';
import { useSingleEvents } from '@/lib/hooks/UseSingleEvent';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const { event, loading } = useSingleEvents(id);
    const [openEdit, setOpenEdit] = useState(false);
    console.log(event);

    const handleCloseEdit = () => {
        setOpenEdit(false);
    };

    const handleOpenEdit = () => {
        setOpenEdit(true);
    };

    return (
        <>
            {loading ? (
                <h1>Loading...</h1>
            ) : (
                <>
                    <EditDialog open={openEdit} handleClose={handleCloseEdit} conEvent={event} />
                    <Card sx={{ maxWidth: '440px' }}>
                        <Button onClick={() => window.history.go(-1)}>Tilbake</Button>
                    </Card>
                    <EventUi conEvent={event} showSelect={true} />

                    <Card sx={{ maxWidth: '440px' }}>
                        <CardActions>
                            <Button onClick={handleOpenEdit}>Endre</Button>
                        </CardActions>
                    </Card>
                </>
            )}
        </>
    );
};

export default Event;
