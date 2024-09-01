import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Box, Button, Typography } from '@mui/material';
import { GetTicketsByEmail, GetTicketsFromCheckIn } from './my-tickets/actions';

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
        <Box>
            <Button> Do some magic</Button>
            <Typography variant="h1">My Profile</Typography>
            <Typography variant="body1">
                Hello world! There'll be stuff here at some point, but bear with me for now. Anyways, how's your day? I
                hope you're doing well.{' '}
            </Typography>
        </Box>
    );
}
