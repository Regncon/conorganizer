import CustomIcon from '$lib/icons/miscgame';
import { Paper, Typography } from '@mui/material';
import Ticket from './UI/Ticket';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';

type Props = {};

const Tickets = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    const verifiedEmail = user?.emailVerified ?? false;
    const verifiedCheckIn = false;

    if (verifiedEmail && verifiedCheckIn) {
        return (
            <Paper sx={{ marginBottom: '2rem', paddingLeft: '2rem', width: '320px' }}>
                <Typography variant="h1">My Tickets</Typography>
                <CustomIcon color="primary" size="large" />
                <CustomIcon color="secondary" size="small" />
                <Ticket />
                <Ticket />
            </Paper>
        );
    }
    return null;
};

export default Tickets;
