'use client';
import { useEffect, useState } from 'react';
import { Box, Button, Card, CardActions } from '@mui/material';
import { useAuth } from '@/components/AuthProvider';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';
import MainNavigator from '@/components/mainNavigator';
import { useUserSettings } from '@/lib/hooks/useUserSettings';
import { useSingleEvents } from '@/lib/hooks/UseSingleEvent';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const { event, loading } = useSingleEvents(id);
    const user = useAuth();
    const { conAuthorization } = useUserSettings(user?.uid);
    const [showEditButton, setShowEditButton] = useState<boolean>(false);
    const [openEdit, setOpenEdit] = useState(false);

    useEffect(() => {
        setShowEditButton(conAuthorization?.admin || false);
    }, [conAuthorization]);

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

            <EventUi conEvent={event} />

            <Card sx={conAuthorization?.admin ? { display: 'block' } : { display: 'none' }}>
                <CardActions>
                    <Button onClick={handleOpenEdit}>Endre</Button>
                </CardActions>
            </Card>
        </Box>
    );
};

export default Event;
