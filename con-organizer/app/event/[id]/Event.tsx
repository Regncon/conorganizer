'use client';
import { useEffect, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { Box, Button, Card, CardActions } from '@mui/material';
import { useAuth } from '@/components/AuthProvider';
import EditDialog from '@/components/EditDialog';
import EventBoundary from '@/components/ErrorBoundaries/EventBoundary';
import EventUi from '@/components/EventUi';
import { useSingleEvents } from '@/lib/hooks/UseSingleEvent';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
type Props = { id: string };
const Event = ({ id }: Props) => {
    const { event, loading } = useSingleEvents(id);
    const user = useAuth();
    const { userSettings } = useUserSettings(user?.uid);
    const [showEditButton, setShowEditButton] = useState<boolean>(false);
    const [openEdit, setOpenEdit] = useState(false);

    useEffect(() => {
        setShowEditButton(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

    const handleCloseEdit = () => {
        setOpenEdit(false);
    };
    const handleOpenEdit = () => {
        setOpenEdit(true);
    };
    // throw new Error(
    //     'lorem Ipsum error in conAuthor authorization dialog box - invalid Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ea quia in blanditiis mollitia exercitationem, asperiores nam quidem commodi nulla illum laborum, distinctio magnam debitis vitae rerum, maiores maxime sapiente! Quia! Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ea quia in blanditiis mollitia exercitationem, asperiores nam quidem commodi nulla illum laborum, distinctio magnam debitis vitae rerum, maiores maxime sapiente! Quia!'
    // );
    return (
        <Box sx={{ maxWidth: '1080px', margin: { xs: '0', md: '5rem auto' } }}>
            {loading && <h1>Loading...</h1>}
            <EditDialog open={openEdit} handleClose={handleCloseEdit} conEvent={event} />
            <Card>
                <Button onClick={() => window.history.go(-1)}>Tilbake</Button>
            </Card>
            <ErrorBoundary FallbackComponent={EventBoundary}>
                <EventUi conEvent={event} />
            </ErrorBoundary>
            <Card sx={showEditButton ? { display: 'block' } : { display: 'none' }}>
                <CardActions>
                    <Button onClick={handleOpenEdit}>Endre</Button>
                </CardActions>
            </Card>
        </Box>
    );
};

export default Event;
