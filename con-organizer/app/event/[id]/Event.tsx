'use client';
import { useState } from 'react';
import { Button, Card, CardActions } from '@mui/material';
import { collection } from 'firebase/firestore';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';
import { useSingleEvents } from '@/lib/hooks/UseSingleEvent';
import db from '../../../lib/firebase';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const collectionRef = collection(db, 'schools');
    const [openEdit, setOpenEdit] = useState(false);

    const { event, loading } = useSingleEvents(id);
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
                    <EditDialog
                        open={openEdit}
                        handleClose={handleCloseEdit}
                        collectionRef={collectionRef}
                        conEvent={event}
                    />

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
