import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Box, Button, Paper, Typography } from '@mui/material';
import { GetTicketsByEmail, GetTicketsFromCheckIn } from './my-tickets/actions';
import MyEvents from '../dashboard/MyEvents';
import MyTickets from '../dashboard/MyTickets';

export default async function MyProfile() {
    const { user } = await getAuthorizedAuth();

    const hasGoogleProvider = user?.providerData.some((provider) => provider.providerId === 'google.com');
    console.log(hasGoogleProvider, 'hasGoogleProvider');
    if (hasGoogleProvider) {
        const tickets = await GetTicketsByEmail(user?.email);
        console.log(tickets, 'tickets');
    }

    //console.log(JSON.stringify(await GetTicketsFromCheckIn()), 'GetTicketsFromCheckIn');

    return (
        <Paper>
            <Typography variant="h1">My Profile</Typography>
            <Typography variant="body1">
                Hello world! There'll be stuff here at some point, but bear with me for now. Anyways, how's your day? I
                hope you're doing well.
            </Typography>
            <Typography variant="body1"> Anyways, here are the events you sent in. </Typography>
            <MyEvents />
            <Typography variant="body1"> And here are your tickets. I think? </Typography>
            <MyTickets />
            <Typography variant="body1"> Is there something wrong? </Typography>
            <Button variant="contained" color="primary" href="/">
                I'm not supposed to be here!
            </Button>
        </Paper>
    );
}
