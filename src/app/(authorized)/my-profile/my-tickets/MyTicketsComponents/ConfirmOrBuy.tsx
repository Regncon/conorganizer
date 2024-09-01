import { Paper, Typography, Button, Box } from '@mui/material';
import ConfirmEmailButton from './UI/ConfirmEmail';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import GoogleSignInButton from '$app/(auth)/login/GoogleButton';
import LaunchIcon from '@mui/icons-material/Launch';

type Props = {};

const ConfirmOrBuy = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    return (
        <Paper sx={{ marginBottom: '2rem', paddingInline: '0.5rem', maxWidth: '320px' }}>
            <Typography variant="h1">Bekreft e-post/Mine billetter</Typography>
            <Typography variant="h2">Har billetter</Typography>
            <Box sx={{ display: 'grid', gap: '1rem' }}>
                <ConfirmEmailButton />
                {user?.emailVerified ? null : <GoogleSignInButton />}
            </Box>
            <Typography variant="h2">Har ikke billetter</Typography>
            <Button
                fullWidth
                variant="contained"
                href="https://event.checkin.no/73685/regncon-xxxii-2024"
                color="secondary"
            >
                Kjøp billett <LaunchIcon sx={{ marginInlineStart: '1rem' }} />
            </Button>
        </Paper>
    );
};

export default ConfirmOrBuy;
