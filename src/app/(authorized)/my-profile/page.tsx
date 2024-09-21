import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Paper, Typography } from '@mui/material';
import { GetTicketsByEmail } from './my-tickets/components/lib/actions/actions';
import MyEvents from '../dashboard/components/MyEvents';
import MyTickets from '../dashboard/components/MyTickets';

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
        <Paper sx={{ padding: '1rem' }}>
            <Typography variant="h1">Min profil</Typography>
            <MyEvents />
            <MyTickets />
        </Paper>
    );
}
